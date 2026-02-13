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

const PaymentScannerCollection = "payment_scanners"

//
// ================= SERVICE INTERFACE =================
//

type PaymentScannerService interface {
	Create(ctx context.Context, companyCode string, req *requests.CreatePaymentScannerRequest) (*models.PaymentScanner, error)
	GetAll(ctx context.Context, companyCode string) ([]*models.PaymentScanner, error)
	GetByID(ctx context.Context, companyCode string, id string) (*models.PaymentScanner, error)
	GetByUUIDs(ctx context.Context, companyCode string, ids []string) ([]*models.PaymentScanner, error)
	Update(ctx context.Context, companyCode string, id string, req *requests.UpdatePaymentScannerRequest) (*models.PaymentScanner, error)
	Delete(ctx context.Context, companyCode string, id string) error
}

//
// ================= SERVICE STRUCT =================
//

type paymentScannerService struct{}

func NewPaymentScannerService() PaymentScannerService {
	return &paymentScannerService{}
}

//
// ================= CREATE =================
//

func (s *paymentScannerService) Create(
	ctx context.Context,
	companyCode string,
	req *requests.CreatePaymentScannerRequest,
) (*models.PaymentScanner, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentScannerCollection)

	paymentScanner := models.NewPaymentScanner()
	paymentScanner.Bind(req)

	_, err := collection.InsertOne(ctx, paymentScanner)
	if err != nil {
		return nil, err
	}

	return paymentScanner, nil
}

//
// ================= GET ALL =================
//

func (s *paymentScannerService) GetAll(
	ctx context.Context,
	companyCode string,
) ([]*models.PaymentScanner, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentScannerCollection)

	cursor, err := collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var paymentScanners []*models.PaymentScanner
	if err := cursor.All(ctx, &paymentScanners); err != nil {
		return nil, err
	}

	return paymentScanners, nil
}

//
// ================= GET BY ID =================
//

func (s *paymentScannerService) GetByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.PaymentScanner, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentScannerCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var paymentScanner models.PaymentScanner
	err := collection.FindOne(ctx, filter).Decode(&paymentScanner)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("payment scanner not found")
	}
	if err != nil {
		return nil, err
	}

	return &paymentScanner, nil
}

//
// ================= GET BY UUIDs =================
//

func (s *paymentScannerService) GetByUUIDs(
	ctx context.Context,
	companyCode string,
	ids []string,
) ([]*models.PaymentScanner, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentScannerCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"entity_id":  bson.M{"$in": ids},
		"is_deleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var paymentScanners []*models.PaymentScanner
	if err := cursor.All(ctx, &paymentScanners); err != nil {
		return nil, err
	}

	return paymentScanners, nil
}

//
// ================= UPDATE =================
//

func (s *paymentScannerService) Update(
	ctx context.Context,
	companyCode string,
	id string,
	req *requests.UpdatePaymentScannerRequest,
) (*models.PaymentScanner, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentScannerCollection)

	updateFields := bson.M{}

	if req.StudentEntityID != nil {
		updateFields["student_entity_id"] = *req.StudentEntityID
	}
	if req.ExamEntityID != nil {
		updateFields["exam_entity_id"] = *req.ExamEntityID
	}
	if req.PaymentID != nil {
		updateFields["payment_id"] = *req.PaymentID
	}
	if req.PaymentDate != nil {
		updateFields["payment_date"] = *req.PaymentDate
	}
	if req.PaymentMethod != nil {
		updateFields["payment_method"] = *req.PaymentMethod
	}
	if req.Amount != nil {
		updateFields["amount"] = *req.Amount
	}
	if req.Status != nil {
		updateFields["status"] = *req.Status
	}
	if req.TransactionID != nil {
		updateFields["transaction_id"] = *req.TransactionID
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

	var updated models.PaymentScanner
	err := collection.
		FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).
		Decode(&updated)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("payment scanner not found")
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

//
// ================= DELETE (SOFT DELETE) =================
//

func (s *paymentScannerService) Delete(
	ctx context.Context,
	companyCode string,
	id string,
) error {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentScannerCollection)

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
		return errors.New("payment scanner not found")
	}

	return nil
}
