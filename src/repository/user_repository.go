package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AntonioDaria/surfe/src/models"
)

type Repository interface {
	GetUserByID(userID int) (*models.User, error)
}

type RepositoryImpl struct {
	users []models.User
}

// NewUserRepo loads user data from a JSON file and initializes UserRepo
func NewUserRepo(filePath string) (*RepositoryImpl, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user file: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(file, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &RepositoryImpl{users: users}, nil
}

// GetUserByID retrieves a user by their ID
func (r *RepositoryImpl) GetUserByID(userID int) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with ID %d not found", userID)
}
