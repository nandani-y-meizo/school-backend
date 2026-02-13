package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nandani-y-meizo/school-backend/requests"
	"github.com/nandani-y-meizo/school-backend/services"
	// "shared/middleware"
)

func CreateBoard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// // Check access claims
	// _, err := middleware.GetAccessClaims(c)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// Get company code from URL param
	companyCode := c.Param("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code is required"})
		return
	}

	// Bind and validate JSON payload
	req := requests.NewCreateBoardRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create board
	service := services.NewBoardService()
	board, err := service.Create(ctx, companyCode, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, board)
}

func CreateClass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	companyCode := c.Param("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code is required"})
		return
	}

	req := requests.NewCreateClassRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service := services.NewClassService()
	class, err := service.Create(ctx, companyCode, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, class)
}

func CreateBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// // Check access claims
	// _, err := middleware.GetAccessClaims(c)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// Get company code from URL param
	companyCode := c.Param("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code is required"})
		return
	}

	// Bind and validate JSON payload
	req := requests.NewCreateBookRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create book
	service := services.NewBookService()
	book, err := service.Create(ctx, companyCode, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func CreateExam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get company code from URL param
	companyCode := c.Param("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code is required"})
		return
	}

	// Bind and validate JSON payload
	req := requests.NewCreateExamRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create exam
	service := services.NewExamService()
	exam, err := service.Create(ctx, companyCode, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exam)
}

func CreateStudent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// // Check access claims
	// _, err := middleware.GetAccessClaims(c)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// Get company code from URL param
	companyCode := c.Param("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code is required"})
		return
	}

	// Bind and validate JSON payload
	req := requests.NewCreateStudentRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create student
	service := services.NewStudentService()
	student, err := service.Create(ctx, companyCode, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// // Check access claims
	// _, err := middleware.GetAccessClaims(c)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// Get company code from URL param
	companyCode := c.Param("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code is required"})
		return
	}

	// Bind and validate JSON payload
	req := requests.NewCreateUserRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create user
	service := services.NewUserService()
	user, err := service.Create(ctx, companyCode, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func CreatePaymentScanner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// // Check access claims
	// _, err := middleware.GetAccessClaims(c)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// Get company code from URL param
	companyCode := c.Param("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code is required"})
		return
	}

	// Bind and validate JSON payload
	req := requests.NewCreatePaymentScannerRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create payment scanner
	service := services.NewPaymentScannerService()
	paymentScanner, err := service.Create(ctx, companyCode, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentScanner)
}

func GetReceiptByRefNo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// // Check access claims
	// _, err := middleware.GetAccessClaims(c)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// Get company code from URL param
	companyCode := c.Param("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code is required"})
		return
	}

	// Bind and validate JSON payload
	req := requests.NewGetReceiptByRefNoRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to get receipt
	service := services.NewReceiptService()
	receipt, err := service.GetReceiptByRefNo(ctx, companyCode, req.RefNo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receipt)
}
