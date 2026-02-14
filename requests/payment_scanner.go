package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type CreatePaymentScannerRequest struct {
	MachineNo string `json:"machine_no" binding:"required"`
	Tid       string `json:"tid" binding:"required"`
	IsActive  bool   `json:"is_active"`
}

type UpdatePaymentScannerRequest struct {
	MachineNo *string `json:"machine_no,omitempty"`
	Tid       *string `json:"tid,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	IsDeleted *bool   `json:"is_deleted,omitempty"`
}

type UpdatePaymentScannerResponse struct {
	EntityID  string `json:"entity_id"`
	MachineNo string `json:"machine_no"`
	Tid       string `json:"tid"`
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"is_deleted"`
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
