package routes

import "github.com/gin-gonic/gin"

// Routes sets up all API routes
func Routes(api *gin.RouterGroup) {
	// Board routes
	boards := api.Group("/companies/:company_code/boards")
	{
		boards.POST("", CreateBoard)
		boards.GET("", GetBoards)
		boards.GET("/:id", GetBoardByID)
		boards.PUT("/:id", UpdateBoard)
		boards.DELETE("/:id", DeleteBoard)

		// Additional board routes
		boards.POST("/batch", GetBoardsByUUIDs)
	}

	//class routes
	class := api.Group("/companies/:company_code/classes")
	{
		class.POST("", CreateClass)
		class.GET("", GetClasses)
		class.GET("/:id", GetClassByID)
		class.PUT("/:id", UpdateClass)
		class.DELETE("/:id", DeleteClass)
	}

	//books routes
	books := api.Group("/companies/:company_code/books")
	{
		books.POST("", CreateBook)
		books.GET("", GetBooks)
		books.GET("/:id", GetBookByID)
		books.PUT("/:id", UpdateBook)
		books.DELETE("/:id", DeleteBook)
	}

	//exams routes
	exams := api.Group("/companies/:company_code/exams")
	{
		exams.POST("", CreateExam)
		exams.GET("", GetExams)
		exams.GET("/:id", GetExamByID)
		exams.PUT("/:id", UpdateExam)
		exams.DELETE("/:id", DeleteExam)

		// Additional exam routes
		exams.POST("/batch", GetExamsByUUIDs)
	}

	students := api.Group("/companies/:company_code/students")
	{
		students.POST("", CreateStudent)
		students.GET("", GetStudents)
		students.GET("/:id", GetStudentByID)
		students.PUT("/:id", UpdateStudent)
		students.DELETE("/:id", DeleteStudent)

		// Additional student routes
		students.POST("/batch", GetStudentsByUUIDs)
	}

	users := api.Group("/companies/:company_code/users")
	{
		users.POST("", CreateUser)
		users.GET("", GetUsers)
		users.GET("/:id", GetUserByID)
		users.PUT("/:id", UpdateUser)
		users.DELETE("/:id", DeleteUser)

		// Additional user routes
		users.POST("/batch", GetUsersByUUIDs)
	}

	paymentScanners := api.Group("/companies/:company_code/payment-scanners")
	{
		paymentScanners.POST("", CreatePaymentScanner)
		paymentScanners.GET("", GetPaymentScanners)
		paymentScanners.GET("/:id", GetPaymentScannerByID)
		paymentScanners.PUT("/:id", UpdatePaymentScanner)
		paymentScanners.DELETE("/:id", DeletePaymentScanner)

		// Additional payment scanner routes
		paymentScanners.POST("/batch", GetPaymentScannersByUUIDs)
	}

	receipts := api.Group("/companies/:company_code/receipts")
	{
		receipts.POST("/lookup", GetReceiptByRefNo)
	}

	unpaidStudents := api.Group("/companies/:company_code/unpaid-students")
	{
		unpaidStudents.POST("", GetUnpaidStudents)
	}
}
