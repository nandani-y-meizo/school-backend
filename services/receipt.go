package services

import (
	"context"
	"errors"
	"fmt"

	"shared/infra/db/mdb"

	"github.com/nandani-y-meizo/school-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//
// ================= SERVICE INTERFACE =================
//

type ReceiptService interface {
	GetReceiptByRefNo(ctx context.Context, companyCode string, refNo string) (*models.Receipt, error)
}

//
// ================= SERVICE STRUCT =================
//

type receiptService struct{}

func NewReceiptService() ReceiptService {
	return &receiptService{}
}

//
// ================= GET RECEIPT BY REF NO =================
//

func (s *receiptService) GetReceiptByRefNo(
	ctx context.Context,
	companyCode string,
	refNo string,
) (*models.Receipt, error) {

	db := mdb.GetMongo()

	// Get student collection
	studentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("students")

	// Find student by refNo
	var student models.Student
	err := studentCollection.FindOne(ctx, bson.M{"ref_no": refNo, "is_deleted": false}).Decode(&student)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("student not found with this ref no")
	}
	if err != nil {
		return nil, err
	}

	// Get payment scanner collection
	paymentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("payment_scanners")

	// Find all payments for this student
	cursor, err := paymentCollection.Find(ctx, bson.M{
		"student_entity_id": student.EntityID,
		"is_deleted":        false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var paymentScanners []models.PaymentScanner
	if err := cursor.All(ctx, &paymentScanners); err != nil {
		return nil, err
	}

	// Get exam collection for exam details
	examCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("exams")

	// Build payment history with exam details
	paymentHistory := make([]models.PaymentHistoryItem, 0)
	var totalPaid float64
	examMap := make(map[string]models.Exam)

	// First, get all exam details
	for _, payment := range paymentScanners {
		if _, exists := examMap[payment.ExamEntityID]; !exists {
			var exam models.Exam
			err := examCollection.FindOne(ctx, bson.M{
				"entity_id":  payment.ExamEntityID,
				"is_deleted": false,
			}).Decode(&exam)
			if err == nil {
				examMap[payment.ExamEntityID] = exam
			}
		}
	}

	// Build payment history
	for _, payment := range paymentScanners {
		if exam, exists := examMap[payment.ExamEntityID]; exists {
			paymentHistory = append(paymentHistory, models.PaymentHistoryItem{
				ID:            payment.ID,
				EntityID:      payment.EntityID,
				ExamEntityID:  payment.ExamEntityID,
				PaymentID:     payment.PaymentID,
				PaymentDate:   payment.PaymentDate,
				PaymentMethod: payment.PaymentMethod,
				Amount:        payment.Amount,
				Status:        payment.Status,
				TransactionID: payment.TransactionID,
				ExamName:      exam.ExamName,
				ExamAmount:    exam.ExamAmount,
			})
			if payment.Status == "paid" {
				totalPaid += payment.Amount
			}
		}
	}

	// Get all exams for this student's class and board to determine pending payments
	examCursor, err := examCollection.Find(ctx, bson.M{
		"class_entity_id": student.ClassEntityID,
		"board_entity_id": student.BoardEntityID,
		"is_deleted":      false,
	})
	if err != nil {
		return nil, err
	}
	defer examCursor.Close(ctx)

	var allExams []models.Exam
	if err := examCursor.All(ctx, &allExams); err != nil {
		return nil, err
	}

	// Get book collection
	bookCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("books")

	// Get all books for this student's class and board
	bookCursor, err := bookCollection.Find(ctx, bson.M{
		"class_entity_id": student.ClassEntityID,
		"board_entity_id": student.BoardEntityID,
		"is_deleted":      false,
	})
	if err != nil {
		return nil, err
	}
	defer bookCursor.Close(ctx)

	var allBooks []models.Book
	if err := bookCursor.All(ctx, &allBooks); err != nil {
		return nil, err
	}

	// Build pending payments and available exams
	var pendingPayments []models.PendingPayment
	var availableExams models.AvailableExams

	// Initialize slices to ensure empty array in JSON instead of null
	availableExams.Compulsory = make([]models.Exam, 0)
	availableExams.Optional = make([]models.Exam, 0)

	var totalDue float64
	paidExamMap := make(map[string]bool)
	paidBookMap := make(map[string]bool)

	// Get all payment scanners for this student to check actual payment status
	paymentCursor, err := paymentCollection.Find(ctx, bson.M{
		"student_entity_id": student.EntityID,
		"is_deleted":        false,
	})
	if err != nil {
		return nil, err
	}
	defer paymentCursor.Close(ctx)

	var paymentScannersForStatus []models.PaymentScanner
	if err := paymentCursor.All(ctx, &paymentScannersForStatus); err != nil {
		return nil, err
	}

	for _, payment := range paymentScannersForStatus {
		if payment.Status == "paid" {
			paidExamMap[payment.ExamEntityID] = true
			paidBookMap[payment.ExamEntityID] = true
			fmt.Printf("Found paid payment for entity_id: %s\n", payment.ExamEntityID)
		}
	}

	// Process exams
	for _, exam := range allExams {
		// Handle both old fees_paid (boolean) and new fees_type (string) fields
		isCompulsory := false
		if exam.FeesType != "" {
			// New format: use fees_type string
			isCompulsory = exam.FeesType == "compulsory"
		} else {
			// Old format: use fees_paid boolean
			isCompulsory = exam.FeesPaid
		}

		// Check if this exam is actually paid
		isPaid := paidExamMap[exam.EntityID]

		// Update the exam's fees_paid status to reflect actual payment status
		exam.FeesPaid = isPaid

		fmt.Printf("Exam %s (%s): isPaid=%t, original fees_paid=%t, new fees_paid=%t\n",
			exam.ExamName, exam.EntityID, isPaid, exam.FeesPaid, exam.FeesPaid)

		// Categorize for available exams
		if isCompulsory {
			availableExams.Compulsory = append(availableExams.Compulsory, exam)
		} else {
			availableExams.Optional = append(availableExams.Optional, exam)
		}

		// Calculate pending payments (unpaid exams)
		if !isPaid {
			pendingPayments = append(pendingPayments, models.PendingPayment{
				ExamEntityID: exam.EntityID,
				ExamName:     exam.ExamName,
				ExamAmount:   exam.ExamAmount,
				FeesPaid:     isPaid,
				DueAmount:    exam.ExamAmount,
			})
			totalDue += exam.ExamAmount
		}
	}

	// Process books
	var availableBooks models.AvailableBooks
	availableBooks.Compulsory = make([]models.Book, 0)
	availableBooks.Optional = make([]models.Book, 0)

	for _, book := range allBooks {
		// Handle both old fees_paid (boolean) and new fees_type (string) fields
		isCompulsory := false
		if book.FeesType != "" {
			// New format: use fees_type string
			isCompulsory = book.FeesType == "compulsory"
		} else {
			// Old format: use fees_paid boolean
			isCompulsory = book.FeesPaid
		}

		// Check if this book is actually paid
		isPaid := paidBookMap[book.EntityID]

		// Update the book's fees_paid status to reflect actual payment status
		book.FeesPaid = isPaid

		if isCompulsory {
			availableBooks.Compulsory = append(availableBooks.Compulsory, book)
		} else {
			availableBooks.Optional = append(availableBooks.Optional, book)
		}
	}

	// Get board collection for board name
	boardCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("boards")

	// Get class collection for class name
	classCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("classes")

	// Fetch board name
	var board models.Board
	boardName := ""
	err = boardCollection.FindOne(ctx, bson.M{
		"entity_id":  student.BoardEntityID,
		"is_deleted": false,
	}).Decode(&board)
	if err == nil {
		boardName = board.BoardName
	}

	// Fetch class name
	var class models.Class
	className := ""
	err = classCollection.FindOne(ctx, bson.M{
		"entity_id":  student.ClassEntityID,
		"is_deleted": false,
	}).Decode(&class)
	if err == nil {
		className = class.ClassName
	}

	// Build receipt
	receipt := &models.Receipt{
		StudentDetails: models.StudentPaymentDetails{
			ID:            student.ID,
			EntityID:      student.EntityID,
			RefNo:         student.RefNo,
			FirstName:     student.FirstName,
			MiddleName:    student.MiddleName,
			LastName:      student.LastName,
			Div:           student.Div,
			BoardEntityID: student.BoardEntityID,
			ClassEntityID: student.ClassEntityID,
			BoardName:     boardName,
			ClassName:     className,
		},
		PaymentHistory:  paymentHistory,
		TotalPaid:       totalPaid,
		TotalDue:        totalDue,
		PendingPayments: pendingPayments,
		AvailableBooks:  availableBooks,
		AvailableExams:  availableExams,
	}

	return receipt, nil
}
