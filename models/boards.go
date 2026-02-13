package models

import (
	"time"

	"shared/pkgs/uuids"

	"github.com/nandani-y-meizo/school-backend/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Board struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID  string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	BoardID   string             `json:"board_id,omitempty" bson:"board_id,omitempty"`
	BoardName string             `json:"board_name,omitempty" bson:"board_name,omitempty"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateBoard struct {
	BoardID   *string `json:"board_id,omitempty" bson:"board_id,omitempty"`
	BoardName *string `json:"board_name,omitempty" bson:"board_name,omitempty"`
	IsDeleted *bool   `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

//
//  Constructors
//

func NewBoard() *Board {
	now := time.Now().UTC()
	id := primitive.NewObjectID()
	entityID, err := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	if err != nil {
		return nil
	}

	return &Board{
		ID:        id,
		EntityID:  entityID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewUpdateBoard() *UpdateBoard {
	return &UpdateBoard{}
}

//
// Bind create request → model
//

func (b *Board) Bind(request *requests.CreateBoardRequest) {
	b.BoardID = request.BoardID
	b.BoardName = request.BoardName
}

//
// Bind update request → model
//

func (b *UpdateBoard) Bind(request *requests.UpdateBoardRequest) {
	if request.BoardID != nil {
		b.BoardID = request.BoardID
	}

	if request.BoardName != nil {
		b.BoardName = request.BoardName
	}
}
