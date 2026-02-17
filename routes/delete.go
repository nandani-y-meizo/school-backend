package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"shared/middleware"

	"github.com/nandani-y-meizo/school-backend/services"
)

func DeleteBoard(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get company code and board ID
	companyCode := c.Param("company_code")
	id := c.Param("id")

	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Call service to delete
	service := services.NewBoardService()
	err = service.Delete(ctx, companyCode, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Board deleted successfully"})
}
func DeleteClass(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get company code and class ID
	companyCode := c.Param("company_code")
	id := c.Param("id")

	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Call service to delete
	service := services.NewClassService()
	err = service.Delete(ctx, companyCode, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}

func DeleteBook(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get company code and book ID
	companyCode := c.Param("company_code")
	id := c.Param("id")

	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Call service to delete
	service := services.NewBookService()
	err = service.Delete(ctx, companyCode, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func DeleteExam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get company code and exam ID
	companyCode := c.Param("company_code")
	id := c.Param("id")

	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Call service to delete
	service := services.NewExamService()
	err = service.Delete(ctx, companyCode, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exam deleted successfully"})
}
func DeleteStudent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get company code and student ID
	companyCode := c.Param("company_code")
	id := c.Param("id")

	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Call service to delete
	service := services.NewStudentService()
	err = service.Delete(ctx, companyCode, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get company code and student ID
	companyCode := c.Param("company_code")
	id := c.Param("id")

	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Call service to delete
	service := services.NewUserService()
	err = service.Delete(ctx, companyCode, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func DeletePaymentScanner(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access check
	_, err := middleware.GetAccessClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get company code and payment scanner ID
	companyCode := c.Param("company_code")
	id := c.Param("id")

	if companyCode == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code and id are required"})
		return
	}

	// Call service to delete
	service := services.NewPaymentScannerService()
	err = service.Delete(ctx, companyCode, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment scanner deleted successfully"})
}
