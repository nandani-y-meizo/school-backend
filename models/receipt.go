package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Receipt represents a student's payment receipt with complete details
type Receipt struct {
	StudentDetails  StudentPaymentDetails `json:"student_details"`
	PaymentHistory  []PaymentHistoryItem  `json:"payment_history"`
	TotalPaid       float64               `json:"total_paid"`
	TotalDue        float64               `json:"total_due"`
	PendingPayments []PendingPayment      `json:"pending_payments"`
	AvailableBooks  AvailableBooks        `json:"available_books"`
	AvailableExams  AvailableExams        `json:"available_exams"`
}

type AvailableBooks struct {
	Compulsory []Book `json:"compulsory"`
	Optional   []Book `json:"optional"`
}

type AvailableExams struct {
	Compulsory []Exam `json:"compulsory"`
	Optional   []Exam `json:"optional"`
}

// StudentPaymentDetails contains student information
type StudentPaymentDetails struct {
	ID            primitive.ObjectID `json:"id"`
	EntityID      string             `json:"entity_id"`
	RefNo         string             `json:"ref_no"`
	FirstName     string             `json:"first_name"`
	MiddleName    string             `json:"middle_name"`
	LastName      string             `json:"last_name"`
	Div           string             `json:"div"`
	BoardEntityID string             `json:"board_entity_id"`
	ClassEntityID string             `json:"class_entity_id"`
	BoardName     string             `json:"board_name,omitempty"`
	ClassName     string             `json:"class_name,omitempty"`
}

// PaymentHistoryItem contains individual payment details
type PaymentHistoryItem struct {
	ID            primitive.ObjectID `json:"id"`
	EntityID      string             `json:"entity_id"`
	ExamEntityID  string             `json:"exam_entity_id"`
	PaymentID     string             `json:"payment_id"`
	PaymentDate   time.Time          `json:"payment_date"`
	PaymentMethod string             `json:"payment_method"`
	Amount        float64            `json:"amount"`
	Status        string             `json:"status"`
	TransactionID string             `json:"transaction_id"`
	ExamName      string             `json:"exam_name"`
	ExamAmount    float64            `json:"exam_amount"`
}

// PendingPayment contains details of payments that need to be made
type PendingPayment struct {
	ExamEntityID string  `json:"exam_entity_id"`
	ExamName     string  `json:"exam_name"`
	ExamAmount   float64 `json:"exam_amount"`
	FeesPaid     bool    `json:"fees_paid"`
	DueAmount    float64 `json:"due_amount"`
}

// ReceiptRequest for looking up student by refNo
type ReceiptRequest struct {
	RefNo string `json:"ref_no" binding:"required"`
}
