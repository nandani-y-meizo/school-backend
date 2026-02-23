package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nandani-y-meizo/school-backend/services"
)

// RoleVerificationRequest represents the request for role verification
type RoleVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyUserRole verifies user role and returns authentication method
func VerifyUserRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Bind and validate JSON payload
	var req RoleVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf(" [RoleVerification] Request binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf(" [RoleVerification] Received request: %+v\n", req)

	// Call role verification service
	service := services.NewRoleVerificationService()
	result, err := service.VerifyUserRole(ctx, req.Email)
	if err != nil {
		fmt.Printf(" [RoleVerification] Service error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Internal server error",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf(" [RoleVerification] Service result: %+v\n", result)

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "success",
		"data":    result,
	})
}
