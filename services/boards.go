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

const BoardCollection = "boards"

// ================= SERVICE INTERFACE =================

type BoardService interface {
	Create(ctx context.Context, companyCode string, req *requests.CreateBoardRequest) (*models.Board, error)
	GetAll(ctx context.Context, companyCode string) ([]*models.Board, error)
	GetByID(ctx context.Context, companyCode string, id string) (*models.Board, error)
	GetByUUIDs(ctx context.Context, companyCode string, ids []string) ([]*models.Board, error)
	Update(ctx context.Context, companyCode string, id string, req *requests.UpdateBoardRequest) (*models.Board, error)
	Delete(ctx context.Context, companyCode string, id string) error
}

// ================= SERVICE STRUCT =================

type boardService struct{}

func NewBoardService() BoardService {
	return &boardService{}
}

// ================= CREATE =================

func (s *boardService) Create(
	ctx context.Context,
	companyCode string,
	req *requests.CreateBoardRequest,
) (*models.Board, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BoardCollection)

	board := models.NewBoard()
	board.Bind(req)

	_, err := collection.InsertOne(ctx, board)
	if err != nil {
		return nil, err
	}

	return board, nil
}

// ================= GET ALL =================

func (s *boardService) GetAll(
	ctx context.Context,
	companyCode string,
) ([]*models.Board, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BoardCollection)

	cursor, err := collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var boards []*models.Board
	if err := cursor.All(ctx, &boards); err != nil {
		return nil, err
	}

	return boards, nil
}

// ================= GET BY ID =================

func (s *boardService) GetByID(
	ctx context.Context,
	companyCode string,
	id string,
) (*models.Board, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BoardCollection)

	filter := bson.M{"entity_id": id, "is_deleted": false}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid, "is_deleted": false}
	}

	var board models.Board
	err := collection.FindOne(ctx, filter).Decode(&board)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("board not found")
	}
	if err != nil {
		return nil, err
	}

	return &board, nil
}

// ================= GET BY UUIDs =================

func (s *boardService) GetByUUIDs(
	ctx context.Context,
	companyCode string,
	ids []string,
) ([]*models.Board, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BoardCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"entity_id":  bson.M{"$in": ids},
		"is_deleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var boards []*models.Board
	if err := cursor.All(ctx, &boards); err != nil {
		return nil, err
	}

	return boards, nil
}

// ================= UPDATE =================

func (s *boardService) Update(
	ctx context.Context,
	companyCode string,
	id string,
	req *requests.UpdateBoardRequest,
) (*models.Board, error) {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BoardCollection)

	updateFields := bson.M{}

	if req.BoardID != nil {
		updateFields["board_id"] = *req.BoardID
	}
	if req.BoardName != nil {
		updateFields["board_name"] = *req.BoardName
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

	var updated models.Board
	err := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": updateFields}, opts).Decode(&updated)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("board not found")
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

// ================= DELETE (SOFT DELETE) =================

func (s *boardService) Delete(
	ctx context.Context,
	companyCode string,
	id string,
) error {

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).
		Collection(BoardCollection)

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
		return errors.New("board not found")
	}

	return nil
}
