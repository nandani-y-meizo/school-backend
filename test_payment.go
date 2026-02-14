package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"shared/infra/db/mdb"

	"github.com/nandani-y-meizo/school-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

func testPayment() {
	// Connect to MongoDB
	db := mdb.GetMongo()

	// Test payment creation
	paymentCollection := db.GetClient().Database("company_SCHOOL001").Collection("payment_scanners")

	// Create a test payment record
	payment := models.NewPaymentScanner()
	payment.StudentEntityID = "test_student"
	payment.ExamEntityID = "test_exam"
	payment.PaymentID = "TEST_PAY_001"
	payment.PaymentDate = time.Now()
	payment.PaymentMethod = "cash"
	payment.Amount = 100.0
	payment.Status = "paid"
	payment.TransactionID = "TEST_TXN_001"

	result, err := paymentCollection.InsertOne(context.Background(), payment)
	if err != nil {
		log.Printf("Error creating test payment: %v", err)
		return
	}

	fmt.Printf("Test payment created with ID: %s\n", result.InsertedID)

	// Query payments to verify
	cursor, err := paymentCollection.Find(context.Background(), bson.M{"is_deleted": false})
	if err != nil {
		log.Printf("Error querying payments: %v", err)
		return
	}
	defer cursor.Close(context.Background())

	var payments []models.PaymentScanner
	cursor.All(context.Background(), &payments)

	fmt.Printf("Found %d payments in database:\n", len(payments))
	for i, p := range payments {
		fmt.Printf("%d. Student: %s, Exam: %s, Amount: %.2f, Status: %s\n",
			i+1, p.StudentEntityID, p.ExamEntityID, p.Amount, p.Status)
	}
}
