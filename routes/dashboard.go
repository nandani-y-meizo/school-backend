package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nandani-y-meizo/school-backend/services"
)

func GetDashboardStats(c *gin.Context) {
	companyCode := c.Param("company_code")

	service := services.NewDashboardService()
	stats, err := service.GetDashboardStats(c.Request.Context(), companyCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
