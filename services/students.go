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

const StudentCollection = "students"

//
// ================= SERVICE INTERFACE =================
//

type StudentService interface {
	Create(ctx context.Context, companyCode string, req *requests.CreateStudentRequest) (*models.Student, error)
	GetAll(ctx context.Context, companyCode string) ([]*models.Student, error)
	GetByID(ctx context.Context, companyCode string, id string) (*models.Student, error)
	GetByUUIDs(ctx context.Context, companyCode string, ids []string) ([]*models.Student, error)
	Update(ctx context.Context, companyCode string, id string, req *requests.UpdateStudentRequest) (*models.Student, error)
	Delete(ctx context.Context, companyCode string, id string) error
}

//
// ================= SERVICE STRUCT =================
//

type studentService struct{}

func NewStudentService() StudentService {
	return &studentService{}
}

//
// ================= CREATE =================
//

func (s *studentService) Create(
	ctx context.Context,
	companyCode string,
	req *requests.CreateStudentRequest,
) (*models.Student, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(StudentCollection)

	student := models.NewStudent()
	student.Bind(req)

	_, err := collection.InsertOne(ctx, student)
	if err != nil {
		return nil, err
	}

	return student, nil
}

//
// ================= GET ALL =================
//

func (s *studentService) GetAll(
	ctx context.Context,
	companyCode string,
) ([]*models.Student, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(StudentCollection)

	cursor, err := collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var students []*models.Student
	if err := cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	return students, nil
}

//
// ================= GET BY ID =================
//

func (s *studentService) GetByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.Student, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(StudentCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var student models.Student
	err := collection.FindOne(ctx, filter).Decode(&student)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("student not found")
	}
	if err != nil {
		return nil, err
	}

	return &student, nil
}

//
// ================= GET BY UUIDs =================
//

func (s *studentService) GetByUUIDs(
	ctx context.Context,
	companyCode string,
	ids []string,
) ([]*models.Student, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(StudentCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"entity_id":  bson.M{"$in": ids},
		"is_deleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var students []*models.Student
	if err := cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	return students, nil
}

//
// ================= UPDATE =================
//

func (s *studentService) Update(
	ctx context.Context,
	companyCode string,
	id string,
	req *requests.UpdateStudentRequest,
) (*models.Student, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(StudentCollection)

	updateFields := bson.M{}

	if req.BoardEntityID != nil {
		updateFields["board_entity_id"] = *req.BoardEntityID
	}
	if req.ClassEntityID != nil {
		updateFields["class_entity_id"] = *req.ClassEntityID
	}
	if req.RefNo != nil {
		updateFields["ref_no"] = *req.RefNo
	}
	if req.Div != nil {
		updateFields["div"] = *req.Div
	}
	if req.FirstName != nil {
		updateFields["first_name"] = *req.FirstName
	}
	if req.MiddleName != nil {
		updateFields["middle_name"] = *req.MiddleName
	}
	if req.LastName != nil {
		updateFields["last_name"] = *req.LastName
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

	var updated models.Student
	err := collection.
		FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).
		Decode(&updated)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("student not found")
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

//
// ================= DELETE (SOFT DELETE) =================
//

func (s *studentService) Delete(
	ctx context.Context,
	companyCode string,
	id string,
) error {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(StudentCollection)

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
		return errors.New("student not found")
	}

	return nil
}
