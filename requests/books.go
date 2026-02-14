package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type CreateBookRequest struct {
	BookID        string  `json:"book_id" binding:"required"`
	BoardEntityID string  `json:"board_entity_id" binding:"required"`
	ClassEntityID string  `json:"class_entity_id" binding:"required"`
	BookName      string  `json:"book_name" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	FeesPaid      bool    `json:"fees_paid"`
	FeesType      string  `json:"fees_type,omitempty"`
}

type UpdateBookRequest struct {
	BookID        *string  `json:"book_id,omitempty"`
	BoardEntityID *string  `json:"board_entity_id,omitempty"`
	ClassEntityID *string  `json:"class_entity_id,omitempty"`
	BookName      *string  `json:"book_name,omitempty"`
	Amount        *float64 `json:"amount,omitempty"`
	FeesPaid      *bool    `json:"fees_paid,omitempty"`
	FeesType      *string  `json:"fees_type,omitempty"`
	IsDeleted     *bool    `json:"is_deleted,omitempty"`
}

type UpdateBookResponse struct {
	BookName  string  `json:"book_name"`
	Amount    float64 `json:"amount"`
	FeesPaid  bool    `json:"fees_paid"`
	IsDeleted bool    `json:"is_deleted"`
}

//
// ================= CONSTRUCTORS =================
//

func NewCreateBookRequest() *CreateBookRequest {
	return &CreateBookRequest{}
}

func NewUpdateBookRequest() *UpdateBookRequest {
	return &UpdateBookRequest{}
}

//
// ================= VALIDATION =================
//

func (r *CreateBookRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateBookRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}
