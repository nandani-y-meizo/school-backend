package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type ConfirmPaymentRequest struct {
	StudentRefNo  string   `json:"student_ref_no" binding:"required"`
	PaymentMode   string   `json:"payment_mode" binding:"required"`
	SelectedExams []string `json:"selected_exams,omitempty"`
	SelectedBooks []string `json:"selected_books,omitempty"`
	TotalAmount   float64  `json:"total_amount" binding:"required"`
}

//
// ================= CONSTRUCTORS =================
//

func NewConfirmPaymentRequest() *ConfirmPaymentRequest {
	return &ConfirmPaymentRequest{}
}

//
// ================= VALIDATION =================
//

func (r *ConfirmPaymentRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}
