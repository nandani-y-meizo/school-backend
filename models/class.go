package models

import (
	"shared/pkgs/uuids"
	"time"

	"github.com/nandani-y-meizo/school-backend/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID      string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	BoardEntityID string             `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassName     string             `json:"class_name,omitempty" bson:"class_name,omitempty"`
	IsDeleted     bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateClass struct {
	BoardEntityID *string `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassName     *string `json:"class_name,omitempty" bson:"class_name,omitempty"`
	IsDeleted     *bool   `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

//
//  Constructors
//

func NewClass() *Class {
	now := time.Now().UTC()
	id := primitive.NewObjectID()
	entityID, err := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	if err != nil {
		return nil
	}

	return &Class{
		ID:        id,
		EntityID: entityID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewUpdateClass() *UpdateClass {
	return &UpdateClass{}
}

//
// Bind create request → model
//

func (c *Class) Bind(req *requests.CreateClassRequest) {
	c.BoardEntityID = req.BoardEntityID
	c.ClassName = req.ClassName
}

//
// Bind update request → model
//

func (c *UpdateClass) Bind(req *requests.UpdateClassRequest) {
	if req.BoardEntityID != nil {
		c.BoardEntityID = req.BoardEntityID
	}
	if req.ClassName != nil {
		c.ClassName = req.ClassName
	}
	if req.IsDeleted != nil {
		c.IsDeleted = req.IsDeleted
	}
}
