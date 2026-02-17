package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"shared/middleware"

	"github.com/nandani-y-meizo/school-backend/requests"
	"github.com/nandani-y-meizo/school-backend/services"
)

func UpdateBoard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Company code and Board ID
	companyCode := c.Param("company_code")
	id := c.Param("id")
	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Bind and validate JSON request
	req := requests.NewUpdateBoardRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update board
	service := services.NewBoardService()
	updatedBoard, err := service.Update(ctx, companyCode, id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedBoard)
}

func UpdateClass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Company code and Class ID
	companyCode := c.Param("company_code")
	id := c.Param("id")
	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Bind and validate JSON request
	req := requests.NewUpdateClassRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update class
	service := services.NewClassService()
	updatedClass, err := service.Update(ctx, companyCode, id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedClass)
}
func UpdateBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Company code and Book ID
	companyCode := c.Param("company_code")
	id := c.Param("id")
	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Bind and validate JSON request
	req := requests.NewUpdateBookRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update book
	service := services.NewBookService()
	updatedBook, err := service.Update(ctx, companyCode, id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}

func UpdateExam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Company code and Exam ID
	companyCode := c.Param("company_code")
	id := c.Param("id")
	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Bind and validate JSON request
	req := requests.NewUpdateExamRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update exam
	service := services.NewExamService()
	updatedExam, err := service.Update(ctx, companyCode, id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedExam)
}
func UpdateStudent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Company code and Student ID
	companyCode := c.Param("company_code")
	id := c.Param("id")
	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Bind and validate JSON request
	req := requests.NewUpdateStudentRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update student
	service := services.NewStudentService()
	updatedStudent, err := service.Update(ctx, companyCode, id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedStudent)
}

func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Company code and Student ID
	companyCode := c.Param("company_code")
	id := c.Param("id")
	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Bind and validate JSON request
	req := requests.NewUpdateUserRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update student
	service := services.NewUserService()
	updatedStudent, err := service.Update(ctx, companyCode, id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedStudent)
}

func UpdatePaymentScanner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Company code and PaymentScanner ID
	companyCode := c.Param("company_code")
	id := c.Param("id")
	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Bind and validate JSON request
	req := requests.NewUpdatePaymentScannerRequest()
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update payment scanner
	service := services.NewPaymentScannerService()
	updatedPaymentScanner, err := service.Update(ctx, companyCode, id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPaymentScanner)
}
