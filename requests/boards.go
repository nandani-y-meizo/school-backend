package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type CreateBoardRequest struct {
	BoardID   string `json:"board_id" binding:"required"`
	BoardName string `json:"board_name" binding:"required"`
}

type UpdateBoardRequest struct {
	BoardID   *string `json:"board_id,omitempty"`
	BoardName *string `json:"board_name,omitempty"`
	IsDeleted *bool   `json:"is_deleted,omitempty"`
}

//
//  Constructors
//

func NewCreateBoardRequest() *CreateBoardRequest {
	return &CreateBoardRequest{}
}

func NewUpdateBoardRequest() *UpdateBoardRequest {
	return &UpdateBoardRequest{}
}

//
//  Validation
//

func (r *CreateBoardRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	// Add any custom validation logic here if needed
	return nil
}

func (r *UpdateBoardRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	// Add any custom validation logic here if needed
	return nil
}

