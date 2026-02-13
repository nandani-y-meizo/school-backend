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

const ExamCollection = "exams"

//
// ================= SERVICE INTERFACE =================
//

type ExamService interface {
	Create(ctx context.Context, companyCode string, req *requests.CreateExamRequest) (*models.Exam, error)
	GetAll(ctx context.Context, companyCode string) ([]*models.Exam, error)
	GetByID(ctx context.Context, companyCode string, id string) (*models.Exam, error)
	GetByUUIDs(ctx context.Context, companyCode string, ids []string) ([]*models.Exam, error)
	Update(ctx context.Context, companyCode string, id string, req *requests.UpdateExamRequest) (*models.Exam, error)
	Delete(ctx context.Context, companyCode string, id string) error
}

//
// ================= SERVICE STRUCT =================
//

type examService struct{}

func NewExamService() ExamService {
	return &examService{}
}

//
// ================= CREATE =================
//

func (s *examService) Create(
	ctx context.Context,
	companyCode string,
	req *requests.CreateExamRequest,
) (*models.Exam, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ExamCollection)

	exam := models.NewExam()
	exam.Bind(req)

	_, err := collection.InsertOne(ctx, exam)
	if err != nil {
		return nil, err
	}

	return exam, nil
}

//
// ================= GET ALL =================
//

func (s *examService) GetAll(
	ctx context.Context,
	companyCode string,
) ([]*models.Exam, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ExamCollection)

	cursor, err := collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var exams []*models.Exam
	if err := cursor.All(ctx, &exams); err != nil {
		return nil, err
	}

	return exams, nil
}

//
// ================= GET BY ID =================
//

func (s *examService) GetByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.Exam, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ExamCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var exam models.Exam
	err := collection.FindOne(ctx, filter).Decode(&exam)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("exam not found")
	}
	if err != nil {
		return nil, err
	}

	return &exam, nil
}

//
// ================= GET BY UUIDs =================
//

func (s *examService) GetByUUIDs(
	ctx context.Context,
	companyCode string,
	ids []string,
) ([]*models.Exam, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ExamCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"entity_id":  bson.M{"$in": ids},
		"is_deleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var exams []*models.Exam
	if err := cursor.All(ctx, &exams); err != nil {
		return nil, err
	}

	return exams, nil
}

//
// ================= UPDATE =================
//

func (s *examService) Update(
	ctx context.Context,
	companyCode string,
	id string,
	req *requests.UpdateExamRequest,
) (*models.Exam, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ExamCollection)

	updateFields := bson.M{}

	if req.BoardEntityID != nil {
		updateFields["board_entity_id"] = *req.BoardEntityID
	}
	if req.ClassEntityID != nil {
		updateFields["class_entity_id"] = *req.ClassEntityID
	}
	if req.ExamName != nil {
		updateFields["exam_name"] = *req.ExamName
	}
	if req.ExamAmount != nil {
		updateFields["exam_amount"] = *req.ExamAmount
	}
	if req.FeesPaid != nil {
		updateFields["fees_paid"] = *req.FeesPaid
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

	var updated models.Exam
	err := collection.
		FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).
		Decode(&updated)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("exam not found")
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

//
// ================= DELETE (SOFT DELETE) =================
//

func (s *examService) Delete(
	ctx context.Context,
	companyCode string,
	id string,
) error {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(ExamCollection)

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
		return errors.New("exam not found")
	}

	return nil
}
