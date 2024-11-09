package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetUserByID(t *testing.T) {
	// Arrange
	userRepo, err := NewUserRepo("./data/users.json")
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}

	// Act
	user, err := userRepo.GetUserByID(1)

	// Assert
	if err != nil {
		t.Fatalf("failed to get user by ID: %v", err)
	}

	if user.ID != 1 {
		t.Fatalf("expected user ID to be 1, got %d", user.ID)
	}

	assert.Equal(t, "Ferdinande", user.Name)
}
