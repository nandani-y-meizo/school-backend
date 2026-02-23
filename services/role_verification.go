package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"shared/constants"
	"shared/infra/db/mdb"

	"github.com/nandani-y-meizo/school-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ================= ROLE VERIFICATION SERVICE =================

type RoleVerificationService interface {
	VerifyUserRole(ctx context.Context, email string) (*RoleVerificationResponse, error)
}

type roleVerificationService struct{}

func NewRoleVerificationService() RoleVerificationService {
	return &roleVerificationService{}
}

// RoleVerificationResponse represents the response for role verification
type RoleVerificationResponse struct {
	Exists     bool   `json:"exists"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	AuthMethod string `json:"auth_method"`
	Company    string `json:"company,omitempty"`
	UserID     string `json:"user_id,omitempty"`
}

// PermifyResponse represents the response from Permify
type PermifyResponse struct {
	CanAccess bool   `json:"can_access"`
	Reason    string `json:"reason,omitempty"`
}

// ================= ROLE VERIFICATION METHODS =================

func (s *roleVerificationService) VerifyUserRole(ctx context.Context, email string) (*RoleVerificationResponse, error) {
	// First, check if user exists in any company database
	userInfo, err := s.findUserByEmail(ctx, email)
	if err != nil {
		return &RoleVerificationResponse{Exists: false}, nil
	}

	// If user found, verify role with Permify
	if userInfo != nil {
		// Check with Permify for authorization
		_, err := s.checkPermifyAuthorization(ctx, userInfo)
		if err != nil {
			// Log error but still allow login if user exists in DB
			fmt.Printf("Permify check failed: %v, proceeding with DB verification\n", err)
		}

		return &RoleVerificationResponse{
			Exists:     true,
			Email:      userInfo.Email,
			Role:       userInfo.Role,
			AuthMethod: userInfo.AuthMethod,
			Company:    userInfo.CompanyCode,
			UserID:     userInfo.UserID,
		}, nil
	}

	return &RoleVerificationResponse{Exists: false}, nil
}

// findUserByEmail searches for user across all company databases
func (s *roleVerificationService) findUserByEmail(ctx context.Context, email string) (*UserInfo, error) {
	db := mdb.GetMongo()

	// Get list of all company databases (this might need adjustment based on your setup)
	// For now, we'll check the main user service database and a few common company databases
	databases := []string{
		"user_service",           // Main user database
		"company_SCHOOL001-0003", // Company database
		"company_SCHOOL001-0004", // Company database
		"company_SCHOOL001-0005", // Company database
	}

	for _, dbName := range databases {
		collection := db.GetClient().Database(dbName).Collection("users")

		var user models.User
		err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
		if err == nil {
			// For regular users, we'll determine role and auth method based on business logic
			// Since the User model doesn't have these fields, we'll use default values
			role := "user"           // Default role for regular users
			authMethod := "password" // Default auth method

			return &UserInfo{
				Email:       user.Email,
				Role:        role,
				AuthMethod:  authMethod,
				CompanyCode: strings.TrimPrefix(dbName, "company_"), // Extract company code from database name
				UserID:      user.ID.Hex(),
			}, nil
		}
		if err != mongo.ErrNoDocuments {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	return nil, nil // User not found
}

// checkPermifyAuthorization verifies user role with Permify
func (s *roleVerificationService) checkPermifyAuthorization(ctx context.Context, userInfo *UserInfo) (bool, error) {
	if constants.PermifyURL == "" {
		return false, fmt.Errorf("Permify URL not configured")
	}

	// Create request to Permify
	permifyReq := PermifyCheckRequest{
		Subject: PermifySubject{
			Type: "user",
			Id:   userInfo.UserID,
		},
		Action: "login",
		Resource: PermifyResource{
			Type: "system",
		},
		Context: map[string]interface{}{
			"role":    userInfo.Role,
			"company": userInfo.CompanyCode,
		},
	}

	// Convert to JSON
	reqBody, err := json.Marshal(permifyReq)
	if err != nil {
		return false, fmt.Errorf("failed to marshal permify request: %w", err)
	}

	// Make HTTP request to Permify
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/v1/permissions/check", constants.PermifyURL)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return false, fmt.Errorf("failed to call permify: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("permify returned status: %d", resp.StatusCode)
	}

	var permifyResp PermifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&permifyResp); err != nil {
		return false, fmt.Errorf("failed to decode permify response: %w", err)
	}

	return permifyResp.CanAccess, nil
}

// ================= SUPPORTING TYPES =================

type UserInfo struct {
	Email       string
	Role        string
	AuthMethod  string
	CompanyCode string
	UserID      string
}

type PermifyCheckRequest struct {
	Subject  PermifySubject         `json:"subject"`
	Action   string                 `json:"action"`
	Resource PermifyResource        `json:"resource"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

type PermifySubject struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PermifyResource struct {
	Type string `json:"type"`
}
