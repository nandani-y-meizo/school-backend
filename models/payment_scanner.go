package models

import (
	"time"

	"shared/pkgs/uuids"

	"github.com/nandani-y-meizo/school-backend/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentScanner struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID        string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	StudentEntityID string             `json:"student_entity_id,omitempty" bson:"student_entity_id,omitempty"`
	ExamEntityID    string             `json:"exam_entity_id,omitempty" bson:"exam_entity_id,omitempty"`
	PaymentID       string             `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	PaymentDate     time.Time          `json:"payment_date,omitempty" bson:"payment_date,omitempty"`
	PaymentMethod   string             `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	Amount          float64            `json:"amount,omitempty" bson:"amount,omitempty"`
	Status          string             `json:"status,omitempty" bson:"status,omitempty"` // paid, pending, failed
	TransactionID   string             `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	IsDeleted       bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdatePaymentScanner struct {
	StudentEntityID *string    `json:"student_entity_id,omitempty" bson:"student_entity_id,omitempty"`
	ExamEntityID    *string    `json:"exam_entity_id,omitempty" bson:"exam_entity_id,omitempty"`
	PaymentID       *string    `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	PaymentDate     *time.Time `json:"payment_date,omitempty" bson:"payment_date,omitempty"`
	PaymentMethod   *string    `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	Amount          *float64   `json:"amount,omitempty" bson:"amount,omitempty"`
	Status          *string    `json:"status,omitempty" bson:"status,omitempty"`
	TransactionID   *string    `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	IsDeleted       *bool      `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

//
// ================= CONSTRUCTORS =================
//

func NewPaymentScanner() *PaymentScanner {
	now := time.Now().UTC()
	id := primitive.NewObjectID()

	entityID, err := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	if err != nil {
		return nil
	}

	return &PaymentScanner{
		ID:        id,
		EntityID:  entityID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewUpdatePaymentScanner() *UpdatePaymentScanner {
	return &UpdatePaymentScanner{}
}

//
// ================= BIND CREATE =================
//

func (p *PaymentScanner) Bind(req *requests.CreatePaymentScannerRequest) {
	p.StudentEntityID = req.StudentEntityID
	p.ExamEntityID = req.ExamEntityID
	p.PaymentID = req.PaymentID
	p.PaymentDate = req.PaymentDate
	p.PaymentMethod = req.PaymentMethod
	p.Amount = req.Amount
	p.Status = req.Status
	p.TransactionID = req.TransactionID
}

//
// ================= BIND UPDATE =================
//

func (p *UpdatePaymentScanner) Bind(req *requests.UpdatePaymentScannerRequest) {
	if req.StudentEntityID != nil {
		p.StudentEntityID = req.StudentEntityID
	}
	if req.ExamEntityID != nil {
		p.ExamEntityID = req.ExamEntityID
	}
	if req.PaymentID != nil {
		p.PaymentID = req.PaymentID
	}
	if req.PaymentDate != nil {
		p.PaymentDate = req.PaymentDate
	}
	if req.PaymentMethod != nil {
		p.PaymentMethod = req.PaymentMethod
	}
	if req.Amount != nil {
		p.Amount = req.Amount
	}
	if req.Status != nil {
		p.Status = req.Status
	}
	if req.TransactionID != nil {
		p.TransactionID = req.TransactionID
	}
	if req.IsDeleted != nil {
		p.IsDeleted = req.IsDeleted
	}
}
