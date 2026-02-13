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
	var paymentHistory []models.PaymentHistoryItem
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

	// Get all exams for this student's class to determine pending payments
	examCursor, err := examCollection.Find(ctx, bson.M{
		"class_entity_id": student.ClassEntityID,
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

	// Build pending payments
	var pendingPayments []models.PendingPayment
	var totalDue float64
	paidExamMap := make(map[string]bool)

	// Mark exams that have been paid
	for _, payment := range paymentScanners {
		if payment.Status == "paid" {
			paidExamMap[payment.ExamEntityID] = true
		}
	}

	// Find unpaid exams
	for _, exam := range allExams {
		if !paidExamMap[exam.EntityID] {
			pendingPayments = append(pendingPayments, models.PendingPayment{
				ExamEntityID: exam.EntityID,
				ExamName:     exam.ExamName,
				ExamAmount:   exam.ExamAmount,
				FeesPaid:     exam.FeesPaid,
				DueAmount:    exam.ExamAmount,
			})
			totalDue += exam.ExamAmount
		}
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
		},
		PaymentHistory:  paymentHistory,
		TotalPaid:       totalPaid,
		TotalDue:        totalDue,
		PendingPayments: pendingPayments,
	}

	return receipt, nil
}
