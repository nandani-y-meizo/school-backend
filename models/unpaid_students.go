package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UnpaidStudent represents a student with pending payments
type UnpaidStudent struct {
	ID            primitive.ObjectID `json:"id"`
	EntityID      string             `json:"entity_id"`
	RefNo         string             `json:"ref_no"`
	FirstName     string             `json:"first_name"`
	MiddleName    string             `json:"middle_name"`
	LastName      string             `json:"last_name"`
	Div           string             `json:"div"`
	BoardEntityID string             `json:"board_entity_id"`
	ClassEntityID string             `json:"class_entity_id"`
	PendingItems  []PendingItem      `json:"pending_items"`
	TotalDue      float64            `json:"total_due"`
}

// PendingItem represents an unpaid item (exam or book)
type PendingItem struct {
	ItemType     string  `json:"item_type"` // "exam" or "book"
	ItemEntityID string  `json:"item_entity_id"`
	ItemName     string  `json:"item_name"`
	ItemAmount   float64 `json:"item_amount"`
	DueAmount    float64 `json:"due_amount"`
	IsCompulsory bool    `json:"is_compulsory"`
}

// UnpaidStudentsResponse for API response
type UnpaidStudentsResponse struct {
	Students []UnpaidStudent `json:"students"`
	Total    int             `json:"total"`
}

// UnpaidStudentsRequest for filtering
type UnpaidStudentsRequest struct {
	ClassEntityID *string `json:"class_entity_id,omitempty"`
	BoardEntityID *string `json:"board_entity_id,omitempty"`
	ItemType      *string `json:"item_type,omitempty"` // "exam", "book", or "all"
}
