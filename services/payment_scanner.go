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

const PaymentDeviceCollection = "payment_scanner_devices"

//
// ================= SERVICE INTERFACE =================
//

type PaymentScannerService interface {
	Create(ctx context.Context, companyCode string, req *requests.CreatePaymentScannerRequest) (*models.PaymentDevice, error)
	GetAll(ctx context.Context, companyCode string) ([]*models.PaymentDevice, error)
	GetByID(ctx context.Context, companyCode string, id string) (*models.PaymentDevice, error)
	GetByUUIDs(ctx context.Context, companyCode string, ids []string) ([]*models.PaymentDevice, error)
	Update(ctx context.Context, companyCode string, id string, req *requests.UpdatePaymentScannerRequest) (*models.PaymentDevice, error)
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
) (*models.PaymentDevice, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentDeviceCollection)

	device := models.NewPaymentDevice()
	device.MachineNo = req.MachineNo
	device.Tid = req.Tid
	device.IsActive = req.IsActive

	_, err := collection.InsertOne(ctx, device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

//
// ================= GET ALL =================
//

func (s *paymentScannerService) GetAll(
	ctx context.Context,
	companyCode string,
) ([]*models.PaymentDevice, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentDeviceCollection)

	cursor, err := collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var devices []*models.PaymentDevice
	if err := cursor.All(ctx, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

//
// ================= GET BY ID =================
//

func (s *paymentScannerService) GetByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.PaymentDevice, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentDeviceCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var device models.PaymentDevice
	err := collection.FindOne(ctx, filter).Decode(&device)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("payment device not found")
	}
	if err != nil {
		return nil, err
	}

	return &device, nil
}

//
// ================= GET BY UUIDs =================
//

func (s *paymentScannerService) GetByUUIDs(
	ctx context.Context,
	companyCode string,
	ids []string,
) ([]*models.PaymentDevice, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentDeviceCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"entity_id":  bson.M{"$in": ids},
		"is_deleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var devices []*models.PaymentDevice
	if err := cursor.All(ctx, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

//
// ================= UPDATE =================
//

func (s *paymentScannerService) Update(
	ctx context.Context,
	companyCode string,
	id string,
	req *requests.UpdatePaymentScannerRequest,
) (*models.PaymentDevice, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentDeviceCollection)

	updateFields := bson.M{}

	if req.MachineNo != nil {
		updateFields["machine_no"] = *req.MachineNo
	}
	if req.Tid != nil {
		updateFields["tid"] = *req.Tid
	}
	if req.IsActive != nil {
		updateFields["is_active"] = *req.IsActive
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

	var updated models.PaymentDevice
	err := collection.
		FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).
		Decode(&updated)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("payment device not found")
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

//
// ================= DELETE =================
//

func (s *paymentScannerService) Delete(
	ctx context.Context,
	companyCode string,
	id string,
) error {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(PaymentDeviceCollection)

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
		return errors.New("payment device not found")
	}

	return nil
}
