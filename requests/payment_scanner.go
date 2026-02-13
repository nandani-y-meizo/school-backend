package requests

import (
	"shared/pkgs/validations"
	"time"

	"github.com/gin-gonic/gin"
)

type CreatePaymentScannerRequest struct {
	StudentEntityID string    `json:"student_entity_id" binding:"required"`
	ExamEntityID    string    `json:"exam_entity_id" binding:"required"`
	PaymentID       string    `json:"payment_id" binding:"required"`
	PaymentDate     time.Time `json:"payment_date" binding:"required"`
	PaymentMethod   string    `json:"payment_method" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	Status          string    `json:"status" binding:"required"`
	TransactionID   string    `json:"transaction_id" binding:"required"`
}

type UpdatePaymentScannerRequest struct {
	StudentEntityID *string    `json:"student_entity_id,omitempty"`
	ExamEntityID    *string    `json:"exam_entity_id,omitempty"`
	PaymentID       *string    `json:"payment_id,omitempty"`
	PaymentDate     *time.Time `json:"payment_date,omitempty"`
	PaymentMethod   *string    `json:"payment_method,omitempty"`
	Amount          *float64   `json:"amount,omitempty"`
	Status          *string    `json:"status,omitempty"`
	TransactionID   *string    `json:"transaction_id,omitempty"`
	IsDeleted       *bool      `json:"is_deleted,omitempty"`
}

type UpdatePaymentScannerResponse struct {
	StudentEntityID string    `json:"student_entity_id"`
	ExamEntityID    string    `json:"exam_entity_id"`
	PaymentID       string    `json:"payment_id"`
	PaymentDate     time.Time `json:"payment_date"`
	PaymentMethod   string    `json:"payment_method"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	TransactionID   string    `json:"transaction_id"`
	IsDeleted       bool      `json:"is_deleted"`
}

//
// ================= CONSTRUCTORS =================
//

func NewCreatePaymentScannerRequest() *CreatePaymentScannerRequest {
	return &CreatePaymentScannerRequest{}
}

func NewUpdatePaymentScannerRequest() *UpdatePaymentScannerRequest {
	return &UpdatePaymentScannerRequest{}
}

//
// ================= VALIDATION =================
//

func (r *CreatePaymentScannerRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}

func (r *UpdatePaymentScannerRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}
