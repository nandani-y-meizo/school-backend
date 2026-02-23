package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"shared/infra/db/mdb"
	"shared/pkgs/jwtmanager"

	"github.com/gin-gonic/gin"
	"github.com/nandani-y-meizo/school-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

// RegularUserLoginRequest represents login request for regular users
type RegularUserLoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	CompanyCode string `json:"company_code" binding:"required"`
}

// RegularUserLoginResponse represents response for regular user login
type RegularUserLoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        *models.User `json:"user"`
	Company     string       `json:"company"`
}

// LoginRegularUser handles login for regular users from company databases
func LoginRegularUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Bind and validate JSON payload
	var req RegularUserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Connect to company database
	db := mdb.GetMongo()
	dbName := fmt.Sprintf("company_%s", req.CompanyCode)
	collection := db.GetClient().Database(dbName).Collection("users")

	// Find user by email
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email, "is_deleted": false}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "User not found",
			"error":   "user not found",
		})
		return
	}

	// TODO: In production, you should hash the password and compare with bcrypt
	// For now, we'll do simple comparison (NOT SECURE - FIX THIS!)
	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "Invalid credentials",
			"error":   "invalid password",
		})
		return
	}

	// Generate proper JWT token with user claims
	customClaims := map[string]interface{}{
		"user_id":      user.ID.Hex(),
		"email":        user.Email,
		"name":         user.Name,
		"username":     user.Name, // Add username field for middleware
		"role":         "user",
		"company_code": req.CompanyCode,
		"group_code":   "", // Add group_code field for middleware
	}

	tokenPair, err := jwtmanager.GenerateTokenPair(customClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
		return
	}
	accessToken := tokenPair.AccessToken

	// Return success response
	companyObj := gin.H{
		"company_code": req.CompanyCode,
		"code":         req.CompanyCode,
		"entity_id":    req.CompanyCode,
		"name":         req.CompanyCode,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "success",
		"data": gin.H{
			"access_token": accessToken,
			"user":         user,
			"company":      companyObj,
		},
	})
}
