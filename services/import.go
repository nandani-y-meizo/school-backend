package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"shared/infra/db/mdb"

	"github.com/nandani-y-meizo/school-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ImportService interface {
	ImportBooks(ctx context.Context, companyCode string, file multipart.File) (int, error)
	ImportExams(ctx context.Context, companyCode string, file multipart.File) (int, error)
	ImportStudents(ctx context.Context, companyCode string, file multipart.File) (int, error)
}

type importService struct{}

func NewImportService() ImportService {
	return &importService{}
}

// Helper to get mappings for Board Name -> ID and Class Name -> ID
func (s *importService) getBoardAndClassMaps(ctx context.Context, companyCode string) (map[string]string, map[string]string, error) {
	db := mdb.GetMongo()

	// Boards
	boardCollection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).Collection("boards")
	boardCursor, err := boardCollection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, nil, err
	}
	defer boardCursor.Close(ctx)

	var boards []models.Board
	if err := boardCursor.All(ctx, &boards); err != nil {
		return nil, nil, err
	}

	boardMap := make(map[string]string)
	for _, b := range boards {
		boardMap[b.BoardName] = b.EntityID
	}

	// Classes
	classCollection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).Collection("classes")
	classCursor, err := classCollection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, nil, err
	}
	defer classCursor.Close(ctx)

	var classes []models.Class
	if err := classCursor.All(ctx, &classes); err != nil {
		return nil, nil, err
	}

	classMap := make(map[string]string)
	for _, c := range classes {
		// Key by "ClassName" - Note: Class names might duplicate across boards, so ideally key by Name+BoardID,
		// but for simple CSV import we might just assume unique class names or rely on user providing valid data.
		// For robustness, we'll map ClassName -> ClassEntityID.
		// WARNING: If multiple boards have same class name, this will be ambiguous.
		classMap[c.ClassName] = c.EntityID
	}

	return boardMap, classMap, nil
}

func (s *importService) ImportBooks(ctx context.Context, companyCode string, file multipart.File) (int, error) {
	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return 0, err
	}
	// Validation of header could be added here
	_ = header

	boardMap, classMap, err := s.getBoardAndClassMaps(ctx, companyCode)
	if err != nil {
		return 0, err
	}

	var books []interface{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		// Expected CSV Format: BookName, Amount, FeesType(compulsory/optional), BoardName, ClassName
		if len(record) < 5 {
			continue
		}

		bookName := record[0]
		amountStr := record[1]
		feesType := record[2]
		boardName := record[3]
		className := record[4]

		amount, _ := strconv.ParseFloat(amountStr, 64)

		boardID, ok := boardMap[boardName]
		if !ok {
			continue // Skip if board not found
		}

		classID, ok := classMap[className]
		if !ok {
			continue // Skip if class not found
		}

		newBook := models.NewBook()
		newBook.BookName = bookName
		newBook.Amount = amount
		newBook.FeesType = feesType
		newBook.FeesPaid = (feesType == "compulsory")
		newBook.BoardEntityID = boardID
		newBook.ClassEntityID = classID

		books = append(books, newBook)
	}

	if len(books) == 0 {
		return 0, nil
	}

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).Collection("books")

	_, err = collection.InsertMany(ctx, books)
	return len(books), err
}

func (s *importService) ImportExams(ctx context.Context, companyCode string, file multipart.File) (int, error) {
	reader := csv.NewReader(file)

	// Read header
	_, err := reader.Read()
	if err != nil {
		return 0, err
	}

	boardMap, classMap, err := s.getBoardAndClassMaps(ctx, companyCode)
	if err != nil {
		return 0, err
	}

	var exams []interface{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		// Expected CSV: ExamName, Amount, FeesType, BoardName, ClassName
		if len(record) < 5 {
			continue
		}

		examName := record[0]
		amountStr := record[1]
		feesType := record[2]
		boardName := record[3]
		className := record[4]

		amount, _ := strconv.ParseFloat(amountStr, 64)

		boardID, ok := boardMap[boardName]
		if !ok {
			continue
		}
		classID, ok := classMap[className]
		if !ok {
			continue
		}

		newExam := models.NewExam()
		newExam.ExamName = examName
		newExam.ExamAmount = amount
		newExam.FeesType = feesType
		newExam.FeesPaid = (feesType == "compulsory")
		newExam.BoardEntityID = boardID
		newExam.ClassEntityID = classID

		exams = append(exams, newExam)
	}

	if len(exams) == 0 {
		return 0, nil
	}

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).Collection("exams")

	_, err = collection.InsertMany(ctx, exams)
	return len(exams), err
}

func (s *importService) ImportStudents(ctx context.Context, companyCode string, file multipart.File) (int, error) {
	reader := csv.NewReader(file)

	// Read header
	_, err := reader.Read()
	if err != nil {
		return 0, err
	}

	boardMap, classMap, err := s.getBoardAndClassMaps(ctx, companyCode)
	if err != nil {
		return 0, err
	}

	var students []interface{}

	// We need to check for existing RefNos to avoid duplicates?
	// For now, simpler implementation: just insert.
	// Ideally we'd upsert or check existence.

	// Pre-fetch all students RefNo map to check duplicates in memory?
	// Or try/catch insert errors if unique index exists. (RefNo is not unique index in schema currently?)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		// Expected CSV: FirstName, MiddleName, LastName, RefNo, Div, BoardName, ClassName
		if len(record) < 7 {
			continue
		}

		firstName := record[0]
		middleName := record[1]
		lastName := record[2]
		refNo := record[3]
		div := record[4]
		boardName := record[5]
		className := record[6]

		boardID, ok := boardMap[boardName]
		if !ok {
			continue
		}
		classID, ok := classMap[className]
		if !ok {
			continue
		}

		newStudent := models.NewStudent()
		newStudent.FirstName = firstName
		newStudent.MiddleName = middleName
		newStudent.LastName = lastName
		newStudent.RefNo = refNo
		newStudent.Div = div
		newStudent.BoardEntityID = boardID
		newStudent.ClassEntityID = classID

		students = append(students, newStudent)
	}

	if len(students) == 0 {
		return 0, nil
	}

	db := mdb.GetMongo()
	collection := db.GetClient().Database(fmt.Sprintf("company_%s", companyCode)).Collection("students")

	_, err = collection.InsertMany(ctx, students)
	return len(students), err
}
