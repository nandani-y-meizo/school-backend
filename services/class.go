package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nandani-y-meizo/school-backend/models"
	"github.com/nandani-y-meizo/school-backend/requests"
	"shared/infra/db/mdb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ClassCollection = "classes"

//
// ================= SERVICE INTERFACE =================
//

type ClassService interface {
	Create(ctx context.Context, companyCode string, req *requests.CreateClassRequest) (*models.Class, error)
	GetAll(ctx context.Context, companyCode string) ([]*models.Class, error)
	GetByID(ctx context.Context, companyCode string, id string) (*models.Class, error)
	GetByUUIDs(ctx context.Context, companyCode string, ids []string) ([]*models.Class, error)
	Update(ctx context.Context, companyCode string, id string, req *requests.UpdateClassRequest) (*models.Class, error)
	Delete(ctx context.Context, companyCode string, id string) error
}

//
// ================= SERVICE STRUCT =================
//

type classService struct{}

func NewClassService() ClassService {
	return &classService{}
}

//
// ================= CREATE =================
//

func (s *classService) Create(
	ctx context.Context,
	companyCode string,
	req *requests.CreateClassRequest,
) (*models.Class, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ClassCollection)

	class := models.NewClass()
	class.Bind(req)

	_, err := collection.InsertOne(ctx, class)
	if err != nil {
		return nil, err
	}

	return class, nil
}

//
// ================= GET ALL =================
//

func (s *classService) GetAll(
	ctx context.Context,
	companyCode string,
) ([]*models.Class, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ClassCollection)

	cursor, err := collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var classes []*models.Class
	if err := cursor.All(ctx, &classes); err != nil {
		return nil, err
	}

	return classes, nil
}

//
// ================= GET BY ID =================
//

func (s *classService) GetByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.Class, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ClassCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var class models.Class
	err := collection.FindOne(ctx, filter).Decode(&class)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("class not found")
	}
	if err != nil {
		return nil, err
	}

	return &class, nil
}

//
// ================= GET BY UUIDs =================
//

func (s *classService) GetByUUIDs(
	ctx context.Context,
	companyCode string,
	ids []string,
) ([]*models.Class, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ClassCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"entity_id":  bson.M{"$in": ids},
		"is_deleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var classes []*models.Class
	if err := cursor.All(ctx, &classes); err != nil {
		return nil, err
	}

	return classes, nil
}

//
// ================= UPDATE =================
//

func (s *classService) Update(
	ctx context.Context,
	companyCode string,
	id string,
	req *requests.UpdateClassRequest,
) (*models.Class, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ClassCollection)

	updateFields := bson.M{}

	if req.BoardEntityID != nil {
		updateFields["board_entity_id"] = *req.BoardEntityID
	}
	if req.ClassName != nil {
		updateFields["class_name"] = *req.ClassName
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

	var updated models.Class
	err := collection.
		FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).
		Decode(&updated)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("class not found")
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

//
// ================= DELETE (SOFT DELETE) =================
//

func (s *classService) Delete(
	ctx context.Context,
	companyCode string,
	id string,
) error {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ClassCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	result, err := collection.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("class not found")
	}

	return nil
}
