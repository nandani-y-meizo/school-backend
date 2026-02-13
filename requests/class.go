package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type CreateClassRequest struct {
	BoardEntityID string `json:"board_entity_id" binding:"required"`
	ClassName     string `json:"class_name" binding:"required"`
}

type UpdateClassRequest struct {
	BoardEntityID *string `json:"board_entity_id,omitempty"`
	ClassName     *string `json:"class_name,omitempty"`
	IsDeleted     *bool   `json:"is_deleted,omitempty"`
}

type UpdateClassResponse struct {
	ClassName string `json:"class_name"`
	IsDeleted bool   `json:"is_deleted"`
}

// constructor

func NewCreateClassRequest() *CreateClassRequest {
	return &CreateClassRequest{}
}

func NewUpdateClassRequest() *UpdateClassRequest {
	return &UpdateClassRequest{}
}

//validation

func (r *CreateClassRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateClassRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}
