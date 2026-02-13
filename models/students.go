package models

import (
	"time"

	"shared/pkgs/uuids"

	"github.com/nandani-y-meizo/school-backend/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID      string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	BoardEntityID string             `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassEntityID string             `json:"class_entity_id,omitempty" bson:"class_entity_id,omitempty"`
	RefNo         string             `json:"ref_no,omitempty" bson:"ref_no,omitempty"`
	Div           string             `json:"div,omitempty" bson:"div,omitempty"`
	FirstName     string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	MiddleName    string             `json:"middle_name,omitempty" bson:"middle_name,omitempty"`
	LastName      string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	IsDeleted     bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateStudent struct {
	BoardEntityID *string  `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassEntityID *string  `json:"class_entity_id,omitempty" bson:"class_entity_id,omitempty"`
    RefNo         *string  `json:"ref_no,omitempty" bson:"ref_no,omitempty"`
	Div           *string  `json:"div,omitempty" bson:"div,omitempty"`
	FirstName     *string  `json:"first_name,omitempty" bson:"first_name,omitempty"`
	MiddleName    *string  `json:"middle_name,omitempty" bson:"middle_name,omitempty"`
	LastName      *string  `json:"last_name,omitempty" bson:"last_name,omitempty"`
	IsDeleted     *bool    `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

//
// ================= CONSTRUCTORS =================
//

func NewStudent() *Student {
	now := time.Now().UTC()
	id := primitive.NewObjectID()

	entityID, err := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	if err != nil {
		return nil
	}

	return &Student{
		ID:        id,
		EntityID:  entityID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewUpdateStudent() *UpdateStudent {
	return &UpdateStudent{}
}

//
// ================= BIND CREATE =================
//

func (b *Student) Bind(req *requests.CreateStudentRequest) {
	b.BoardEntityID = req.BoardEntityID
	b.ClassEntityID = req.ClassEntityID
	b.RefNo = req.RefNo
	b.Div = req.Div
	b.FirstName = req.FirstName
	b.MiddleName = req.MiddleName
	b.LastName = req.LastName
}

//
// ================= BIND UPDATE =================
//
func (b *UpdateStudent) Bind(req *requests.UpdateStudentRequest) {

	if req.BoardEntityID != nil {
		b.BoardEntityID = req.BoardEntityID
	}
	if req.ClassEntityID != nil {
		b.ClassEntityID = req.ClassEntityID
	}
	if req.RefNo != nil {
		b.RefNo = req.RefNo
	}
	if req.Div != nil {
		b.Div = req.Div
	}
	if req.FirstName != nil {
		b.FirstName = req.FirstName
	}
	if req.MiddleName != nil {
		b.MiddleName = req.MiddleName
	}
	if req.LastName != nil {
		b.LastName = req.LastName
	}
	if req.IsDeleted != nil {
		b.IsDeleted = req.IsDeleted
	}
}

