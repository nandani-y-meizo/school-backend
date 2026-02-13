package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type CreateStudentRequest struct {
	BoardEntityID string `json:"board_entity_id" binding:"required"`
	ClassEntityID string `json:"class_entity_id" binding:"required"`
	RefNo         string `json:"ref_no" binding:"required"`
	Div           string `json:"div" binding:"required"`
	FirstName     string `json:"first_name" binding:"required"`
	MiddleName    string `json:"middle_name"`
	LastName      string `json:"last_name" binding:"required"`
}

type UpdateStudentRequest struct {
	BoardEntityID *string `json:"board_entity_id,omitempty"`
	ClassEntityID *string `json:"class_entity_id,omitempty"`
	RefNo         *string `json:"ref_no,omitempty"`
	Div           *string `json:"div,omitempty"`
	FirstName     *string `json:"first_name,omitempty"`
	MiddleName    *string `json:"middle_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	IsDeleted     *bool   `json:"is_deleted,omitempty"`
}

type UpdateStudentResponse struct {
	BoardEntityID string `json:"board_entity_id"`
	ClassEntityID string `json:"class_entity_id"`
	RefNo         string `json:"ref_no"`
	Div           string `json:"div"`
	FirstName     string `json:"first_name"`
	MiddleName    string `json:"middle_name"`
	LastName      string `json:"last_name"`
	IsDeleted     bool   `json:"is_deleted"`
}

//
// ================= CONSTRUCTORS =================
//

func NewCreateStudentRequest() *CreateStudentRequest {
	return &CreateStudentRequest{}
}

func NewUpdateStudentRequest() *UpdateStudentRequest {
	return &UpdateStudentRequest{}
}

//
// ================= VALIDATION =================
//

func (r *CreateStudentRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateStudentRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}
