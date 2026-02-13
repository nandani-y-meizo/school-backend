package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type CreateExamRequest struct {
	BoardEntityID string  `json:"board_entity_id" binding:"required"`
	ClassEntityID string  `json:"class_entity_id" binding:"required"`
	ExamName      string  `json:"exam_name" binding:"required"`
	ExamAmount    float64 `json:"exam_amount" binding:"required,gt=0"`
	FeesPaid      bool    `json:"fees_paid"`
}

type UpdateExamRequest struct {
	BoardEntityID *string  `json:"board_entity_id,omitempty"`
	ClassEntityID *string  `json:"class_entity_id,omitempty"`
	ExamName      *string  `json:"exam_name,omitempty"`
	ExamAmount    *float64 `json:"exam_amount,omitempty"`
	FeesPaid      *bool    `json:"fees_paid,omitempty"`
	IsDeleted     *bool    `json:"is_deleted,omitempty"`
}

type UpdateExamResponse struct {
	ExamName  string  `json:"exam_name"`
	ExamAmount    float64 `json:"exam_amount"`
	FeesPaid  bool    `json:"fees_paid"`
	IsDeleted bool    `json:"is_deleted"`
}

//
// ================= CONSTRUCTORS =================
//

func NewCreateExamRequest() *CreateExamRequest {
	return &CreateExamRequest{}
}

func NewUpdateExamRequest() *UpdateExamRequest {
	return &UpdateExamRequest{}
}

//
// ================= VALIDATION =================
//

func (r *CreateExamRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateExamRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}
