package services

import (
	"context"
	"fmt"
	"time"

	"shared/infra/db/mdb"

	"github.com/nandani-y-meizo/school-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

type DashboardService interface {
	GetDashboardStats(ctx context.Context, companyCode string) (*models.DashboardStats, error)
}

type dashboardService struct{}

func NewDashboardService() DashboardService {
	return &dashboardService{}
}

func (s *dashboardService) GetDashboardStats(ctx context.Context, companyCode string) (*models.DashboardStats, error) {
	db := mdb.GetMongo()
	dbName := fmt.Sprintf("company_%s", companyCode)

	// Get payment scanner collection for collection stats
	paymentCollection := db.GetClient().Database(dbName).Collection("payment_scanners")

	// Calculate collection stats
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":       nil,
				"totalPaid": bson.M{"$sum": "$amount"},
				"cashAmount": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$payment_method", "cash"}},
							"$amount",
							0,
						},
					},
				},
				"upiAmount": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$payment_method", "upi"}},
							"$amount",
							0,
						},
					},
				},
			},
		},
	}

	cursor, err := paymentCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var collectionResults []struct {
		TotalPaid  float64 `bson:"totalPaid"`
		CashAmount float64 `bson:"cashAmount"`
		UPIAmount  float64 `bson:"upiAmount"`
	}

	if err := cursor.All(ctx, &collectionResults); err != nil {
		return nil, err
	}

	collectionStats := models.CollectionStats{
		TotalPaidAmount: 0,
		CashAmount:      0,
		UPIAmount:       0,
	}

	if len(collectionResults) > 0 {
		collectionStats.TotalPaidAmount = collectionResults[0].TotalPaid
		collectionStats.CashAmount = collectionResults[0].CashAmount
		collectionStats.UPIAmount = collectionResults[0].UPIAmount
	}

	// Get student stats
	studentCollection := db.GetClient().Database(dbName).Collection("students")
	examCollection := db.GetClient().Database(dbName).Collection("exams")
	bookCollection := db.GetClient().Database(dbName).Collection("books")
	paymentScannerCollection := db.GetClient().Database(dbName).Collection("payment_scanners")

	// 1. Fetch all students
	cursorStudents, err := studentCollection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursorStudents.Close(ctx)
	var students []models.Student
	if err := cursorStudents.All(ctx, &students); err != nil {
		return nil, err
	}

	// 2. Fetch all exams and books to build required items map per class
	cursorExams, err := examCollection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursorExams.Close(ctx)
	var exams []models.Exam
	if err := cursorExams.All(ctx, &exams); err != nil {
		return nil, err
	}

	cursorBooks, err := bookCollection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursorBooks.Close(ctx)
	var books []models.Book
	if err := cursorBooks.All(ctx, &books); err != nil {
		return nil, err
	}

	classRequiredItems := make(map[string][]string)
	for _, exam := range exams {
		classRequiredItems[exam.ClassEntityID] = append(classRequiredItems[exam.ClassEntityID], exam.EntityID)
	}
	for _, book := range books {
		classRequiredItems[book.ClassEntityID] = append(classRequiredItems[book.ClassEntityID], book.EntityID)
	}

	// 3. Fetch all payments to build paid items map per student
	cursorPayments, err := paymentScannerCollection.Find(ctx, bson.M{"is_deleted": false, "status": "paid"})
	if err != nil {
		return nil, err
	}
	defer cursorPayments.Close(ctx)
	var payments []models.PaymentScanner
	if err := cursorPayments.All(ctx, &payments); err != nil {
		return nil, err
	}

	studentPaidItems := make(map[string]map[string]bool)
	for _, payment := range payments {
		if studentPaidItems[payment.StudentEntityID] == nil {
			studentPaidItems[payment.StudentEntityID] = make(map[string]bool)
		}
		studentPaidItems[payment.StudentEntityID][payment.ExamEntityID] = true
	}

	// 4. Calculate paid/unpaid students
	paidStudentsCount := 0
	unpaidStudentsCount := 0

	for _, student := range students {
		requiredItems := classRequiredItems[student.ClassEntityID]
		paidItems := studentPaidItems[student.EntityID]

		allPaid := true
		for _, itemID := range requiredItems {
			if !paidItems[itemID] {
				allPaid = false
				break
			}
		}

		if allPaid {
			paidStudentsCount++
		} else {
			unpaidStudentsCount++
		}
	}

	feesStatus := models.FeesStatusStats{
		PaidStudents:   paidStudentsCount,
		UnpaidStudents: unpaidStudentsCount,
		TotalStudents:  len(students),
	}

	// Create holidays list (you can make this dynamic from DB if needed)
	holidays := []models.Holiday{
		{
			Name: "Republic Day",
			Date: time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC),
			Flag: "🇮🇳",
		},
		{
			Name: "Holi",
			Date: time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC),
			Flag: "🎨",
		},
		{
			Name: "Independence Day",
			Date: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC),
			Flag: "🇮🇳",
		},
		{
			Name: "Gandhi Jayanti",
			Date: time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC),
			Flag: "🕊️",
		},
		{
			Name: "Diwali",
			Date: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
			Flag: "🪔",
		},
	}

	return &models.DashboardStats{
		Collection: collectionStats,
		FeesStatus: feesStatus,
		Holidays:   holidays,
	}, nil
}
