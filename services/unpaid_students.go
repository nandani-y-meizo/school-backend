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

	// Get payment scanner collection
	paymentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("payment_scanners")

	// Get exam collection
	examCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("exams")

	// Get book collection
	bookCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("books")

	var unpaidStudents []models.UnpaidStudent

	// Process each student
	for _, student := range students {
		// Get paid items for this student
		paymentCursor, err := paymentCollection.Find(ctx, bson.M{
			"student_entity_id": student.EntityID,
			"is_deleted":        false,
			"status":            "paid",
		})
		if err != nil {
			continue
		}

		var paidPayments []models.PaymentScanner
		if err := paymentCursor.All(ctx, &paidPayments); err != nil {
			paymentCursor.Close(ctx)
			continue
		}
		paymentCursor.Close(ctx)

		// Create map of paid items
		paidItems := make(map[string]bool)
		for _, payment := range paidPayments {
			paidItems[payment.ExamEntityID] = true
		}

		// Get all exams for this student's class
		examFilter := bson.M{
			"class_entity_id": student.ClassEntityID,
			"is_deleted":      false,
		}
		examCursor, err := examCollection.Find(ctx, examFilter)
		if err != nil {
			continue
		}
		defer examCursor.Close(ctx)

		var exams []models.Exam
		if err := examCursor.All(ctx, &exams); err != nil {
			continue
		}

		// Get all books for this student's class
		bookFilter := bson.M{
			"class_entity_id": student.ClassEntityID,
			"is_deleted":      false,
		}
		bookCursor, err := bookCollection.Find(ctx, bookFilter)
		if err != nil {
			continue
		}
		defer bookCursor.Close(ctx)

		var books []models.Book
		if err := bookCursor.All(ctx, &books); err != nil {
			continue
		}

		// Build pending items
		var pendingItems []models.PendingItem
		var totalDue float64

		// Add unpaid exams
		for _, exam := range exams {
			if !paidItems[exam.EntityID] {
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
					IsCompulsory: exam.FeesPaid,
				})
				totalDue += exam.ExamAmount
			}
		}

		// Add unpaid books
		for _, book := range books {
			bookKey := "book_" + book.EntityID
			if !paidItems[bookKey] {
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
					IsCompulsory: book.FeesPaid,
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
