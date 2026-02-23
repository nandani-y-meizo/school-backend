package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name      *string `json:"name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
	IsDeleted *bool   `json:"is_deleted,omitempty"`
}

//
//  Constructors
//

func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{}
}

func NewUpdateUserRequest() *UpdateUserRequest {
	return &UpdateUserRequest{}
}

//
//  Validation
//

func (r *CreateUserRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	// Add any custom validation logic here if needed
	return nil
}

func (r *UpdateUserRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	// Add any custom validation logic here if needed
	return nil
}
