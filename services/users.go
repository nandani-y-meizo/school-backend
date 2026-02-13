package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shared/infra/db/mdb"

	"github.com/nandani-y-meizo/school-backend/models"
	"github.com/nandani-y-meizo/school-backend/requests"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const UsersCollection = "users"

// ================= SERVICE INTERFACE =================

type UserService interface {
	Create(ctx context.Context, companyCode string, req *requests.CreateUserRequest) (*models.User, error)
	GetAll(ctx context.Context, companyCode string) ([]*models.User, error)
	GetByID(ctx context.Context, companyCode string, id string) (*models.User, error)
	GetByUUIDs(ctx context.Context, companyCode string, ids []string) ([]*models.User, error)
	Update(ctx context.Context, companyCode string, id string, req *requests.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, companyCode string, id string) error
}

// ================= SERVICE STRUCT =================

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

// ================= CREATE =================

func (s *userService) Create(
	ctx context.Context,
	companyCode string,
	req *requests.CreateUserRequest,
) (*models.User, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(UsersCollection)

	user := models.NewUser()
	user.Bind(req)

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ================= GET ALL =================

func (s *userService) GetAll(
	ctx context.Context,
	companyCode string,
) ([]*models.User, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(UsersCollection)

	cursor, err := collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// ================= GET BY ID =================

func (s *userService) GetByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.User, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(UsersCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ================= GET BY UUIDs =================

func (s *userService) GetByUUIDs(
	ctx context.Context,
	companyCode string,
	ids []string,
) ([]*models.User, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(UsersCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"entity_id":  bson.M{"$in": ids},
		"is_deleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// ================= UPDATE =================

func (s *userService) Update(
	ctx context.Context,
	companyCode string,
	id string,
	req *requests.UpdateUserRequest,
) (*models.User, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(UsersCollection)

	updateFields := bson.M{}

	if req.Name != nil {
		updateFields["name"] = *req.Name
	}
	if req.Email != nil {
		updateFields["email"] = *req.Email
	}
	if req.IsDeleted != nil {
		updateFields["is_deleted"] = *req.IsDeleted
	}

	if len(updateFields) == 0 {
		return nil, errors.New("no fields to update")
	}

	updateFields["updated_at"] = time.Now()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var updated models.User
	err := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).Decode(&updated)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

// ================= DELETE (SOFT DELETE) =================

func (s *userService) Delete(
	ctx context.Context,
	companyCode string,
	id string,
) error {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(UsersCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"is_deleted": true,
		"updated_at": time.Now(),
	}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
