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

const BookCollection = "books"

//
// ================= SERVICE INTERFACE =================
//

type BookService interface {
	Create(ctx context.Context, companyCode string, req *requests.CreateBookRequest) (*models.Book, error)
	GetAll(ctx context.Context, companyCode string) ([]*models.Book, error)
	GetByID(ctx context.Context, companyCode string, id string) (*models.Book, error)
	GetByUUIDs(ctx context.Context, companyCode string, ids []string) ([]*models.Book, error)
	Update(ctx context.Context, companyCode string, id string, req *requests.UpdateBookRequest) (*models.Book, error)
	Delete(ctx context.Context, companyCode string, id string) error
}

//
// ================= SERVICE STRUCT =================
//

type bookService struct{}

func NewBookService() BookService {
	return &bookService{}
}

//
// ================= CREATE =================
//

func (s *bookService) Create(
	ctx context.Context,
	companyCode string,
	req *requests.CreateBookRequest,
) (*models.Book, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BookCollection)

	book := models.NewBook()
	book.Bind(req)

	_, err := collection.InsertOne(ctx, book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

//
// ================= GET ALL =================
//

func (s *bookService) GetAll(
	ctx context.Context,
	companyCode string,
) ([]*models.Book, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BookCollection)

	cursor, err := collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []*models.Book
	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return books, nil
}

//
// ================= GET BY ID =================
//

func (s *bookService) GetByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.Book, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BookCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var book models.Book
	err := collection.FindOne(ctx, filter).Decode(&book)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("book not found")
	}
	if err != nil {
		return nil, err
	}

	return &book, nil
}

//
// ================= GET BY UUIDs =================
//

func (s *bookService) GetByUUIDs(
	ctx context.Context,
	companyCode string,
	ids []string,
) ([]*models.Book, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BookCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"entity_id":  bson.M{"$in": ids},
		"is_deleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []*models.Book
	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return books, nil
}

//
// ================= UPDATE =================
//

func (s *bookService) Update(
	ctx context.Context,
	companyCode string,
	id string,
	req *requests.UpdateBookRequest,
) (*models.Book, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BookCollection)

	updateFields := bson.M{}

	if req.BoardEntityID != nil {
		updateFields["board_entity_id"] = *req.BoardEntityID
	}
	if req.ClassEntityID != nil {
		updateFields["class_entity_id"] = *req.ClassEntityID
	}
	if req.BookName != nil {
		updateFields["book_name"] = *req.BookName
	}
	if req.Amount != nil {
		updateFields["amount"] = *req.Amount
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

	var updated models.Book
	err := collection.
		FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).
		Decode(&updated)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("book not found")
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

//
// ================= DELETE (SOFT DELETE) =================
//

func (s *bookService) Delete(
	ctx context.Context,
	companyCode string,
	id string,
) error {

	db := mdb.GetMongo()
	collection := db.GetClient().
		Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BookCollection)

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
		return errors.New("book not found")
	}

	return nil
}
