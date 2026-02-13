package requests

import (
	"errors"
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type GetUnpaidStudentsRequest struct {
	ClassEntityID *string `json:"class_entity_id,omitempty"`
	BoardEntityID *string `json:"board_entity_id,omitempty"`
	ItemType      *string `json:"item_type,omitempty"` // "exam", "book", or "all"
}

//
// ================= CONSTRUCTORS =================
//

func NewGetUnpaidStudentsRequest() *GetUnpaidStudentsRequest {
	return &GetUnpaidStudentsRequest{}
}

//
// ================= VALIDATION =================
//

func (r *GetUnpaidStudentsRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}

	// Validate item_type if provided
	if r.ItemType != nil {
		itemType := *r.ItemType
		if itemType != "exam" && itemType != "book" && itemType != "all" {
			return errors.New("item_type must be 'exam', 'book', or 'all'")
		}
	}

	return nil
}
