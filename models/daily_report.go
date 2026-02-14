package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyReport struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID       string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	ReportDate     time.Time          `json:"report_date,omitempty" bson:"report_date,omitempty"`
	TotalPayments  int                `json:"total_payments,omitempty" bson:"total_payments,omitempty"`
	TotalAmount    float64            `json:"total_amount,omitempty" bson:"total_amount,omitempty"`
	PaymentDetails []PaymentDetail    `json:"payment_details,omitempty" bson:"payment_details,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type PaymentDetail struct {
	ID              string    `json:"id,omitempty" bson:"id,omitempty"`
	PaymentID       string    `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	StudentEntityID string    `json:"student_entity_id,omitempty" bson:"student_entity_id,omitempty"`
	StudentRefNo    string    `json:"student_ref_no,omitempty" bson:"student_ref_no,omitempty"`
	StudentName     string    `json:"student_name,omitempty" bson:"student_name,omitempty"`
	BoardEntityID   string    `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassEntityID   string    `json:"class_entity_id,omitempty" bson:"class_entity_id,omitempty"`
	ExamEntityID    string    `json:"exam_entity_id,omitempty" bson:"exam_entity_id,omitempty"`
	BookEntityID    string    `json:"book_entity_id,omitempty" bson:"book_entity_id,omitempty"`
	ItemType        string    `json:"item_type,omitempty" bson:"item_type,omitempty"` // "exam" or "book"
	ItemName        string    `json:"item_name,omitempty" bson:"item_name,omitempty"`
	FeesType        string    `json:"fees_type,omitempty" bson:"fees_type,omitempty"` // "compulsory" or "optional"
	Amount          float64   `json:"amount,omitempty" bson:"amount,omitempty"`
	PaymentMethod   string    `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	Status          string    `json:"status,omitempty" bson:"status,omitempty"`
	TransactionID   string    `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	PaymentTime     time.Time `json:"payment_time,omitempty" bson:"payment_time,omitempty"`
}

type ReportSummary struct {
	TotalPayments  int      `json:"total_payments,omitempty"`
	TotalAmount    float64  `json:"total_amount,omitempty"`
	TotalCash      float64  `json:"total_cash,omitempty"`
	TotalUPI       float64  `json:"total_upi,omitempty"`
	PaymentMethods []string `json:"payment_methods,omitempty"`
	PaymentStatus  []string `json:"payment_status,omitempty"`
}

type DailyReportResponse struct {
	Reports []DailyReport `json:"reports,omitempty"`
	Total   int           `json:"total,omitempty"`
	Summary ReportSummary `json:"summary,omitempty"`
}
