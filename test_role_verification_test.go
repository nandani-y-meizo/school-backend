package main

import (
	"fmt"
	"testing"

	"github.com/nandani-y-meizo/school-backend/services"
)

func TestRoleVerification(t *testing.T) {
	service := services.NewRoleVerificationService()

	// Test with known user email
	email := "yeligetinandini@gmail.com"

	result, err := service.VerifyUserRole(nil, email)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
		return
	}

	fmt.Printf("Role verification result for %s: %+v\n", email, result)
}
