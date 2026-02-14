package services

import (
	"context"
	"fmt"

	"shared/infra/db/mdb"

	"github.com/nandani-y-meizo/school-backend/models"
	"github.com/nandani-y-meizo/school-backend/requests"

	"go.mongodb.org/mongo-driver/bson"
)

//
// ================= SERVICE INTERFACE =================
//

type UnpaidStudentsService interface {
	GetUnpaidStudents(ctx context.Context, companyCode string, req *requests.GetUnpaidStudentsRequest) (*models.UnpaidStudentsResponse, error)
}

//
// ================= SERVICE STRUCT =================
//

type unpaidStudentsService struct{}

func NewUnpaidStudentsService() UnpaidStudentsService {
	return &unpaidStudentsService{}
}

//
// ================= GET UNPAID STUDENTS =================
//

func (s *unpaidStudentsService) GetUnpaidStudents(
	ctx context.Context,
	companyCode string,
	req *requests.GetUnpaidStudentsRequest,
) (*models.UnpaidStudentsResponse, error) {

	db := mdb.GetMongo()

	// Get student collection
	studentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("students")

	// Build filter for students
	studentFilter := bson.M{"is_deleted": false}
	if req.ClassEntityID != nil {
		studentFilter["class_entity_id"] = *req.ClassEntityID
	}
	if req.BoardEntityID != nil {
		studentFilter["board_entity_id"] = *req.BoardEntityID
	}

	// Find all students
	cursor, err := studentCollection.Find(ctx, studentFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var students []models.Student
	if err := cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	// Get collections
	examCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("exams")

	bookCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("books")

	paymentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("payment_scanners")

	// Get all exams
	cursorExams, err := examCollection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursorExams.Close(ctx)
	var allExams []models.Exam
	if err := cursorExams.All(ctx, &allExams); err != nil {
		return nil, err
	}

	// Map exams by class
	classExams := make(map[string][]models.Exam)
	for _, exam := range allExams {
		classExams[exam.ClassEntityID] = append(classExams[exam.ClassEntityID], exam)
	}

	// Get all books
	cursorBooks, err := bookCollection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursorBooks.Close(ctx)
	var allBooks []models.Book
	if err := cursorBooks.All(ctx, &allBooks); err != nil {
		return nil, err
	}

	// Map books by class
	classBooks := make(map[string][]models.Book)
	for _, book := range allBooks {
		classBooks[book.ClassEntityID] = append(classBooks[book.ClassEntityID], book)
	}

	// Get all payments to identify paid items
	cursorPayments, err := paymentCollection.Find(ctx, bson.M{"is_deleted": false, "status": "paid"})
	if err != nil {
		return nil, err
	}
	defer cursorPayments.Close(ctx)
	var allPayments []models.PaymentScanner
	if err := cursorPayments.All(ctx, &allPayments); err != nil {
		return nil, err
	}

	// Map paid items by student
	studentPaidItems := make(map[string]map[string]bool)
	for _, payment := range allPayments {
		if studentPaidItems[payment.StudentEntityID] == nil {
			studentPaidItems[payment.StudentEntityID] = make(map[string]bool)
		}
		studentPaidItems[payment.StudentEntityID][payment.ExamEntityID] = true
	}

	var unpaidStudents []models.UnpaidStudent

	// Process each student
	for _, student := range students {
		paidItems := studentPaidItems[student.EntityID]
		exams := classExams[student.ClassEntityID]
		books := classBooks[student.ClassEntityID]

		// Build pending items
		var pendingItems []models.PendingItem
		var totalDue float64

		// Add unpaid exams
		for _, exam := range exams {
			if paidItems == nil || !paidItems[exam.EntityID] {
				// Determine if compulsory
				isCompulsory := false
				if exam.FeesType != "" {
					isCompulsory = exam.FeesType == "compulsory"
				} else {
					isCompulsory = exam.FeesPaid
				}

				itemType := "exam"
				if req.ItemType != nil && *req.ItemType != "all" && *req.ItemType != itemType {
					continue
				}

				pendingItems = append(pendingItems, models.PendingItem{
					ItemType:     "exam",
					ItemEntityID: exam.EntityID,
					ItemName:     exam.ExamName,
					ItemAmount:   exam.ExamAmount,
					DueAmount:    exam.ExamAmount,
					IsCompulsory: isCompulsory,
				})
				totalDue += exam.ExamAmount
			}
		}

		// Add unpaid books
		for _, book := range books {
			if paidItems == nil || !paidItems[book.EntityID] {
				// Determine if compulsory
				isCompulsory := false
				if book.FeesType != "" {
					isCompulsory = book.FeesType == "compulsory"
				} else {
					isCompulsory = book.FeesPaid
				}

				itemType := "book"
				if req.ItemType != nil && *req.ItemType != "all" && *req.ItemType != itemType {
					continue
				}

				pendingItems = append(pendingItems, models.PendingItem{
					ItemType:     "book",
					ItemEntityID: book.EntityID,
					ItemName:     book.BookName,
					ItemAmount:   book.Amount,
					DueAmount:    book.Amount,
					IsCompulsory: isCompulsory,
				})
				totalDue += book.Amount
			}
		}

		// Only add student if they have pending items
		if len(pendingItems) > 0 {
			unpaidStudents = append(unpaidStudents, models.UnpaidStudent{
				ID:            student.ID,
				EntityID:      student.EntityID,
				RefNo:         student.RefNo,
				FirstName:     student.FirstName,
				MiddleName:    student.MiddleName,
				LastName:      student.LastName,
				Div:           student.Div,
				BoardEntityID: student.BoardEntityID,
				ClassEntityID: student.ClassEntityID,
				PendingItems:  pendingItems,
				TotalDue:      totalDue,
			})
		}
	}

	return &models.UnpaidStudentsResponse{
		Students: unpaidStudents,
		Total:    len(unpaidStudents),
	}, nil
}
