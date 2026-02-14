package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"shared/infra/db/mdb"
	"shared/pkgs/uuids"

	"github.com/nandani-y-meizo/school-backend/models"
	"github.com/nandani-y-meizo/school-backend/requests"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyReportService interface {
	GetDailyReports(ctx context.Context, companyCode string, req *requests.DailyReportRequest) (*models.DailyReportResponse, error)
	GetDailyReportByID(ctx context.Context, companyCode string, id string) (*models.DailyReport, error)
	DeleteDailyReport(ctx context.Context, companyCode string, id string) error
}

type dailyReportService struct{}

func NewDailyReportService() DailyReportService {
	return &dailyReportService{}
}

func (s *dailyReportService) GetDailyReports(
	ctx context.Context,
	companyCode string,
	req *requests.DailyReportRequest,
) (*models.DailyReportResponse, error) {
	db := mdb.GetMongo()

	// Get payment scanner collection
	paymentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("payment_scanners")

	// Get student collection for student details
	studentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("students")

	// Get exam collection for exam details
	examCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("exams")

	// Get book collection for book details
	bookCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("books")

	// Build filter
	filter := bson.M{"is_deleted": false}

	// Date range filter
	if req.StartDate != nil && req.EndDate != nil {
		startDate, _ := time.Parse("2006-01-02", *req.StartDate)
		endDate, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = endDate.Add(24 * time.Hour) // Include end date
		filter["payment_date"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
	}

	// Status filter
	if req.Status != nil && *req.Status != "all" {
		filter["status"] = *req.Status
	}

	// Find all payments
	cursor, err := paymentCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var paymentScanners []models.PaymentScanner
	if err := cursor.All(ctx, &paymentScanners); err != nil {
		return nil, err
	}

	// Debug: Log what we found
	fmt.Printf("Found %d payment scanners\n", len(paymentScanners))

	// Get all students for details
	studentIDs := make([]string, 0, len(paymentScanners))
	studentMap := make(map[string]models.Student)
	for _, payment := range paymentScanners {
		if _, exists := studentMap[payment.StudentEntityID]; !exists {
			studentIDs = append(studentIDs, payment.StudentEntityID)
		}
	}

	if len(studentIDs) > 0 {
		studentCursor, err := studentCollection.Find(ctx, bson.M{
			"entity_id":  bson.M{"$in": studentIDs},
			"is_deleted": false,
		})
		if err == nil {
			defer studentCursor.Close(ctx)
			var students []models.Student
			studentCursor.All(ctx, &students)
			for _, student := range students {
				studentMap[student.EntityID] = student
			}
		}
	}

	// Get all exams for details
	examIDs := make([]string, 0, len(paymentScanners))
	examMap := make(map[string]models.Exam)
	for _, payment := range paymentScanners {
		if _, exists := examMap[payment.ExamEntityID]; !exists {
			examIDs = append(examIDs, payment.ExamEntityID)
		}
	}

	if len(examIDs) > 0 {
		examCursor, err := examCollection.Find(ctx, bson.M{
			"entity_id":  bson.M{"$in": examIDs},
			"is_deleted": false,
		})
		if err == nil {
			defer examCursor.Close(ctx)
			var exams []models.Exam
			examCursor.All(ctx, &exams)
			for _, exam := range exams {
				examMap[exam.EntityID] = exam
			}
		}
	}

	// Get all books for details
	bookIDs := make([]string, 0, len(paymentScanners))
	bookMap := make(map[string]models.Book)
	for _, payment := range paymentScanners {
		if _, exists := bookMap[payment.ExamEntityID]; !exists {
			bookIDs = append(bookIDs, payment.ExamEntityID)
		}
	}

	if len(bookIDs) > 0 {
		bookCursor, err := bookCollection.Find(ctx, bson.M{
			"entity_id":  bson.M{"$in": bookIDs},
			"is_deleted": false,
		})
		if err == nil {
			defer bookCursor.Close(ctx)
			var books []models.Book
			bookCursor.All(ctx, &books)
			for _, book := range books {
				bookMap[book.EntityID] = book
			}
		}
	}

	// Build payment details
	var paymentDetails []models.PaymentDetail
	paymentMethods := make(map[string]bool)
	paymentStatus := make(map[string]bool)
	totalAmount := 0.0
	totalCash := 0.0
	totalUPI := 0.0

	for _, payment := range paymentScanners {
		student, studentExists := studentMap[payment.StudentEntityID]
		if !studentExists {
			continue
		}

		var itemName string
		var itemType string
		var itemEntityID string
		var feesType string

		// Check if it's an exam or book
		if exam, exists := examMap[payment.ExamEntityID]; exists {
			itemName = exam.ExamName
			itemType = "exam"
			itemEntityID = exam.EntityID
			feesType = exam.FeesType
			if feesType == "" {
				if exam.FeesPaid {
					feesType = "compulsory"
				} else {
					feesType = "optional"
				}
			}
		} else if book, exists := bookMap[payment.ExamEntityID]; exists {
			itemName = book.BookName
			itemType = "book"
			itemEntityID = book.EntityID
			feesType = book.FeesType
			if feesType == "" {
				if book.FeesPaid {
					feesType = "compulsory"
				} else {
					feesType = "optional"
				}
			}
		} else {
			continue
		}

		// Apply item type filter
		if req.ItemType != nil && *req.ItemType != "all" && *req.ItemType != itemType {
			continue
		}

		// Apply Class filter
		if req.ClassEntityID != nil && *req.ClassEntityID != "" && *req.ClassEntityID != student.ClassEntityID {
			continue
		}

		// Apply Board filter
		if req.BoardEntityID != nil && *req.BoardEntityID != "" && *req.BoardEntityID != student.BoardEntityID {
			continue
		}

		// Apply specific Exam/Book filter
		if itemType == "exam" && req.ExamEntityID != nil && *req.ExamEntityID != "" && *req.ExamEntityID != itemEntityID {
			continue
		}
		if itemType == "book" && req.BookEntityID != nil && *req.BookEntityID != "" && *req.BookEntityID != itemEntityID {
			continue
		}

		studentName := fmt.Sprintf("%s %s", student.FirstName, student.LastName)
		if student.MiddleName != "" {
			studentName = fmt.Sprintf("%s %s %s", student.FirstName, student.MiddleName, student.LastName)
		}

		detail := models.PaymentDetail{
			ID:              payment.ID.Hex(),
			PaymentID:       payment.PaymentID,
			StudentEntityID: payment.StudentEntityID,
			StudentRefNo:    student.RefNo,
			StudentName:     studentName,
			BoardEntityID:   student.BoardEntityID,
			ClassEntityID:   student.ClassEntityID,
			ExamEntityID:    itemEntityID,
			ItemType:        itemType,
			ItemName:        itemName,
			FeesType:        feesType,
			Amount:          payment.Amount,
			PaymentMethod:   payment.PaymentMethod,
			Status:          payment.Status,
			TransactionID:   payment.TransactionID,
			PaymentTime:     payment.PaymentDate,
		}

		if itemType == "book" {
			detail.BookEntityID = itemEntityID
		} else {
			detail.ExamEntityID = itemEntityID
		}

		paymentDetails = append(paymentDetails, detail)
		paymentMethods[payment.PaymentMethod] = true
		paymentStatus[payment.Status] = true
		totalAmount += payment.Amount
		if strings.ToLower(payment.PaymentMethod) == "cash" {
			totalCash += payment.Amount
		} else if strings.ToLower(payment.PaymentMethod) == "upi" {
			totalUPI += payment.Amount
		}
	}

	// Build payment methods list
	methods := make([]string, 0, len(paymentMethods))
	for method := range paymentMethods {
		methods = append(methods, method)
	}

	// Build payment status list
	statuses := make([]string, 0, len(paymentStatus))
	for status := range paymentStatus {
		statuses = append(statuses, status)
	}

	// Group by date for reports
	reportsMap := make(map[string]*models.DailyReport)
	for _, detail := range paymentDetails {
		dateKey := detail.PaymentTime.Format("2006-01-02")

		if report, exists := reportsMap[dateKey]; !exists {
			report = &models.DailyReport{
				EntityID:       generateEntityID(),
				ReportDate:     detail.PaymentTime,
				TotalPayments:  0,
				TotalAmount:    0,
				PaymentDetails: []models.PaymentDetail{},
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			reportsMap[dateKey] = report
		}

		report := reportsMap[dateKey]
		report.TotalPayments++
		report.TotalAmount += detail.Amount
		report.PaymentDetails = append(report.PaymentDetails, detail)
	}

	// Convert map to slice
	var reports []models.DailyReport
	for _, report := range reportsMap {
		reports = append(reports, *report)
	}

	summary := models.ReportSummary{
		TotalPayments:  len(paymentDetails),
		TotalAmount:    totalAmount,
		TotalCash:      totalCash,
		TotalUPI:       totalUPI,
		PaymentMethods: methods,
		PaymentStatus:  statuses,
	}

	return &models.DailyReportResponse{
		Reports: reports,
		Total:   len(reports),
		Summary: summary,
	}, nil
}

func (s *dailyReportService) GetDailyReportByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.DailyReport, error) {
	// This would be implemented if we store daily reports in a separate collection
	// For now, we can return an error or implement as needed
	return nil, fmt.Errorf("daily report by ID not implemented")
}

func (s *dailyReportService) DeleteDailyReport(
	ctx context.Context,
	companyCode string,
	id string,
) error {
	// This would be implemented if we store daily reports in a separate collection
	// For now, we can return an error or implement as needed
	return fmt.Errorf("delete daily report not implemented")
}

func generateEntityID() string {
	id := primitive.NewObjectID()
	entityID, _ := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	return entityID
}
