package services

import (
	"context"
	"fmt"
	"time"

	"shared/infra/db/mdb"
	"shared/pkgs/uuids"

	"github.com/nandani-y-meizo/school-backend/models"
	"github.com/nandani-y-meizo/school-backend/requests"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentConfirmationService interface {
	ConfirmPayment(ctx context.Context, companyCode string, req *requests.ConfirmPaymentRequest) error
}

type paymentConfirmationService struct{}

func NewPaymentConfirmationService() PaymentConfirmationService {
	return &paymentConfirmationService{}
}

func (s *paymentConfirmationService) ConfirmPayment(
	ctx context.Context,
	companyCode string,
	req *requests.ConfirmPaymentRequest,
) error {
	fmt.Printf("ConfirmPayment called with: %+v\n", req)

	db := mdb.GetMongo()

	// Get student collection
	studentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("students")

	// Find student by refNo
	var student models.Student
	err := studentCollection.FindOne(ctx, bson.M{"ref_no": req.StudentRefNo, "is_deleted": false}).Decode(&student)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("student not found with ref no: %s", req.StudentRefNo)
	}
	if err != nil {
		return err
	}

	// Get payment scanner collection
	paymentCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("payment_scanners")

	// Get exam collection for exam details
	examCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("exams")

	// Get book collection for book details
	bookCollection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection("books")

	// Process selected exams
	for _, examEntityID := range req.SelectedExams {
		// Check if payment already exists for this exam
		var existingPayment models.PaymentScanner
		err := paymentCollection.FindOne(ctx, bson.M{
			"student_entity_id": student.EntityID,
			"exam_entity_id":    examEntityID,
			"is_deleted":        false,
		}).Decode(&existingPayment)

		if err == mongo.ErrNoDocuments {
			// Get exam details
			var exam models.Exam
			err := examCollection.FindOne(ctx, bson.M{
				"entity_id":  examEntityID,
				"is_deleted": false,
			}).Decode(&exam)
			if err != nil {
				continue // Skip if exam not found
			}

			// Create new payment record
			paymentScanner := models.NewPaymentScanner()
			paymentScanner.StudentEntityID = student.EntityID
			paymentScanner.ExamEntityID = examEntityID
			paymentScanner.PaymentID = generatePaymentID()
			paymentScanner.PaymentDate = time.Now()
			paymentScanner.PaymentMethod = req.PaymentMode
			paymentScanner.Amount = exam.ExamAmount
			paymentScanner.Status = "paid"
			paymentScanner.TransactionID = generateTransactionID()

			_, err = paymentCollection.InsertOne(ctx, paymentScanner)
			if err != nil {
				return fmt.Errorf("failed to create payment for exam %s: %v", examEntityID, err)
			}

			// Update exam fees_paid status to true
			_, err = examCollection.UpdateOne(ctx, bson.M{
				"entity_id":  examEntityID,
				"is_deleted": false,
			}, bson.M{
				"$set": bson.M{
					"fees_paid":  true,
					"updated_at": time.Now(),
				},
			})
			if err != nil {
				fmt.Printf("Error updating exam fees_paid for %s: %v\n", examEntityID, err)
				return fmt.Errorf("failed to update exam fees_paid status for %s: %v", examEntityID, err)
			}
			fmt.Printf("Updated exam %s fees_paid to true\n", examEntityID)
		}
	}

	// Process selected books
	for _, bookEntityID := range req.SelectedBooks {
		// Check if payment already exists for this book
		var existingPayment models.PaymentScanner
		err := paymentCollection.FindOne(ctx, bson.M{
			"student_entity_id": student.EntityID,
			"exam_entity_id":    bookEntityID, // Using exam_entity_id field for books as well
			"is_deleted":        false,
		}).Decode(&existingPayment)

		if err == mongo.ErrNoDocuments {
			// Get book details
			var book models.Book
			err := bookCollection.FindOne(ctx, bson.M{
				"entity_id":  bookEntityID,
				"is_deleted": false,
			}).Decode(&book)
			if err != nil {
				continue // Skip if book not found
			}

			// Create new payment record
			paymentScanner := models.NewPaymentScanner()
			paymentScanner.StudentEntityID = student.EntityID
			paymentScanner.ExamEntityID = bookEntityID // Using exam_entity_id field for books
			paymentScanner.PaymentID = generatePaymentID()
			paymentScanner.PaymentDate = time.Now()
			paymentScanner.PaymentMethod = req.PaymentMode
			paymentScanner.Amount = book.Amount
			paymentScanner.Status = "paid"
			paymentScanner.TransactionID = generateTransactionID()

			_, err = paymentCollection.InsertOne(ctx, paymentScanner)
			if err != nil {
				return fmt.Errorf("failed to create payment for book %s: %v", bookEntityID, err)
			}

			// Update book fees_paid status to true
			_, err = bookCollection.UpdateOne(ctx, bson.M{
				"entity_id":  bookEntityID,
				"is_deleted": false,
			}, bson.M{
				"$set": bson.M{
					"fees_paid":  true,
					"updated_at": time.Now(),
				},
			})
			if err != nil {
				return fmt.Errorf("failed to update book fees_paid status for %s: %v", bookEntityID, err)
			}
		}
	}

	return nil
}

func generatePaymentID() string {
	return fmt.Sprintf("PAY_%d", time.Now().Unix())
}

func generateTransactionID() string {
	id := primitive.NewObjectID()
	entityID, _ := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	return fmt.Sprintf("TXN_%s", entityID[:8])
}
