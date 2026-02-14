package models

import (
	"time"

	"shared/pkgs/uuids"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentDevice struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID  string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	MachineNo string             `json:"machine_no,omitempty" bson:"machine_no,omitempty"`
	Tid       string             `json:"tid,omitempty" bson:"tid,omitempty"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdatePaymentDevice struct {
	MachineNo *string `json:"machine_no,omitempty" bson:"machine_no,omitempty"`
	Tid       *string `json:"tid,omitempty" bson:"tid,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsDeleted *bool   `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

func NewPaymentDevice() *PaymentDevice {
	now := time.Now().UTC()
	id := primitive.NewObjectID()

	entityID, err := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	if err != nil {
		return nil
	}

	return &PaymentDevice{
		ID:        id,
		EntityID:  entityID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
