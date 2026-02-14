package models

import (
	"time"

	"shared/pkgs/uuids"

	"github.com/nandani-y-meizo/school-backend/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID      string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	BookID        string             `json:"book_id,omitempty" bson:"book_id,omitempty"`
	BoardEntityID string             `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassEntityID string             `json:"class_entity_id,omitempty" bson:"class_entity_id,omitempty"`
	BookName      string             `json:"book_name,omitempty" bson:"book_name,omitempty"`
	Amount        float64            `json:"amount,omitempty" bson:"amount,omitempty"`
	FeesPaid      bool               `json:"fees_paid" bson:"fees_paid"`                     // true = compulsory, false = optional
	FeesType      string             `json:"fees_type,omitempty" bson:"fees_type,omitempty"` // "compulsory" or "optional"
	IsDeleted     bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateBook struct {
	BookID        *string  `json:"book_id,omitempty" bson:"book_id,omitempty"`
	BoardEntityID *string  `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassEntityID *string  `json:"class_entity_id,omitempty" bson:"class_entity_id,omitempty"`
	BookName      *string  `json:"book_name,omitempty" bson:"book_name,omitempty"`
	Amount        *float64 `json:"amount,omitempty" bson:"amount,omitempty"`
	FeesPaid      *bool    `json:"fees_paid,omitempty" bson:"fees_paid,omitempty"`
	FeesType      *string  `json:"fees_type,omitempty" bson:"fees_type,omitempty"`
	IsDeleted     *bool    `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

//
// ================= CONSTRUCTORS =================
//

func NewBook() *Book {
	now := time.Now().UTC()
	id := primitive.NewObjectID()

	entityID, err := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	if err != nil {
		return nil
	}

	return &Book{
		ID:        id,
		EntityID:  entityID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewUpdateBook() *UpdateBook {
	return &UpdateBook{}
}

//
// ================= BIND CREATE =================
//

func (b *Book) Bind(req *requests.CreateBookRequest) {
	b.BookID = req.BookID
	b.BoardEntityID = req.BoardEntityID
	b.ClassEntityID = req.ClassEntityID
	b.BookName = req.BookName
	b.Amount = req.Amount
	b.FeesPaid = req.FeesPaid
	b.FeesType = req.FeesType
}

//
// ================= BIND UPDATE =================
//

func (b *UpdateBook) Bind(req *requests.UpdateBookRequest) {

	if req.BookID != nil {
		b.BookID = req.BookID
	}
	if req.BoardEntityID != nil {
		b.BoardEntityID = req.BoardEntityID
	}
	if req.ClassEntityID != nil {
		b.ClassEntityID = req.ClassEntityID
	}
	if req.BookName != nil {
		b.BookName = req.BookName
	}
	if req.Amount != nil {
		b.Amount = req.Amount
	}
	if req.FeesPaid != nil {
		b.FeesPaid = req.FeesPaid
	}
	if req.FeesType != nil {
		b.FeesType = req.FeesType
	}
	if req.IsDeleted != nil {
		b.IsDeleted = req.IsDeleted
	}
}
