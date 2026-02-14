package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nandani-y-meizo/school-backend/services"
)

func ImportBooks(c *gin.Context) {
	companyCode := c.Param("company_code")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	service := services.NewImportService()
	count, err := service.ImportBooks(c.Request.Context(), companyCode, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Books imported successfully", "count": count})
}

func ImportExams(c *gin.Context) {
	companyCode := c.Param("company_code")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	service := services.NewImportService()
	count, err := service.ImportExams(c.Request.Context(), companyCode, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exams imported successfully", "count": count})
}

func ImportStudents(c *gin.Context) {
	companyCode := c.Param("company_code")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	service := services.NewImportService()
	count, err := service.ImportStudents(c.Request.Context(), companyCode, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Students imported successfully", "count": count})
}
