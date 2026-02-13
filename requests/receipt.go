package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type GetReceiptByRefNoRequest struct {
	RefNo string `json:"ref_no" binding:"required"`
}

//
// ================= CONSTRUCTORS =================
//

func NewGetReceiptByRefNoRequest() *GetReceiptByRefNoRequest {
	return &GetReceiptByRefNoRequest{}
}

//
// ================= VALIDATION =================
//

func (r *GetReceiptByRefNoRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}
