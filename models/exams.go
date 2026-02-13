package models

import (
	"time"

	"shared/pkgs/uuids"

	"github.com/nandani-y-meizo/school-backend/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exam struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID      string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	BoardEntityID string             `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassEntityID string             `json:"class_entity_id,omitempty" bson:"class_entity_id,omitempty"`
	ExamName      string             `json:"exam_name,omitempty" bson:"exam_name,omitempty"`
	ExamAmount    float64            `json:"exam_amount,omitempty" bson:"exam_amount,omitempty"`
	FeesPaid      bool               `json:"fees_paid" bson:"fees_paid"` // true = compulsory, false = optional
	IsDeleted     bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateExam struct {
	BoardEntityID *string  `json:"board_entity_id,omitempty" bson:"board_entity_id,omitempty"`
	ClassEntityID *string  `json:"class_entity_id,omitempty" bson:"class_entity_id,omitempty"`
	ExamName      *string  `json:"exam_name,omitempty" bson:"exam_name,omitempty"`
	ExamAmount    *float64 `json:"exam_amount,omitempty" bson:"exam_amount,omitempty"`
	FeesPaid      *bool    `json:"fees_paid,omitempty" bson:"fees_paid,omitempty"`
	IsDeleted     *bool    `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

//
// ================= CONSTRUCTORS =================
//

func NewExam() *Exam {
	now := time.Now().UTC()
	id := primitive.NewObjectID()

	entityID, err := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	if err != nil {
		return nil
	}

	return &Exam{
		ID:        id,
		EntityID:  entityID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewUpdateExam() *UpdateExam {
	return &UpdateExam{}
}

//
// ================= BIND CREATE =================
//

func (b *Exam) Bind(req *requests.CreateExamRequest) {
	b.BoardEntityID = req.BoardEntityID
	b.ClassEntityID = req.ClassEntityID
	b.ExamName = req.ExamName
	b.ExamAmount = req.ExamAmount
	b.FeesPaid = req.FeesPaid
}

//
// ================= BIND UPDATE =================
//

func (b *UpdateExam) Bind(req *requests.UpdateExamRequest) {

	if req.BoardEntityID != nil {
		b.BoardEntityID = req.BoardEntityID
	}
	if req.ClassEntityID != nil {
		b.ClassEntityID = req.ClassEntityID
	}
	if req.ExamName != nil {
		b.ExamName = req.ExamName
	}
	if req.ExamAmount != nil {
		b.ExamAmount = req.ExamAmount
	}
	if req.FeesPaid != nil {
		b.FeesPaid = req.FeesPaid
	}
	if req.IsDeleted != nil {
		b.IsDeleted = req.IsDeleted
	}
}
