package requests

import (
	"shared/pkgs/validations"

	"github.com/gin-gonic/gin"
)

type DailyReportRequest struct {
	StartDate     *string `json:"start_date,omitempty"`
	EndDate       *string `json:"end_date,omitempty"`
	ItemType      *string `json:"item_type,omitempty"` // "exam", "book", "all"
	Status        *string `json:"status,omitempty"`    // "paid", "pending", "failed", "all"
	ClassEntityID *string `json:"class_entity_id,omitempty"`
	BoardEntityID *string `json:"board_entity_id,omitempty"`
	ExamEntityID  *string `json:"exam_entity_id,omitempty"`
	BookEntityID  *string `json:"book_entity_id,omitempty"`
}

//
// ================= CONSTRUCTORS =================
//

func NewDailyReportRequest() *DailyReportRequest {
	return &DailyReportRequest{}
}

//
// ================= VALIDATION =================
//

func (r *DailyReportRequest) Validate(c *gin.Context) error {
	if err := validations.ValidateJSON(c, r); err != nil {
		return err
	}
	return nil
}
