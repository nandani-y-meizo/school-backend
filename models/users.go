package models

import (
	"time"

	"shared/pkgs/uuids"

	"github.com/nandani-y-meizo/school-backend/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EntityID  string             `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateUser struct {
	Name      *string `json:"name,omitempty" bson:"name,omitempty"`
	Email     *string `json:"email,omitempty" bson:"email,omitempty"`
	Password  *string `json:"password,omitempty" bson:"password,omitempty"`
	IsDeleted *bool   `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

//
//  Constructors
//

func NewUser() *User {
	now := time.Now().UTC()
	id := primitive.NewObjectID()
	entityID, err := uuids.NewUUID5(id.Hex(), uuids.OidNamespace)
	if err != nil {
		return nil
	}

	return &User{
		ID:        id,
		EntityID:  entityID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewUpdateUser() *UpdateUser {
	return &UpdateUser{}
}

//
// Bind create request → model
//

func (b *User) Bind(request *requests.CreateUserRequest) {
	b.Name = request.Name
	b.Email = request.Email
	b.Password = request.Password
}

//
// Bind update request → model
//

func (b *UpdateUser) Bind(request *requests.UpdateUserRequest) {
	if request.Name != nil {
		b.Name = request.Name
	}

	if request.Password != nil {
		b.Password = request.Password
	}
}
