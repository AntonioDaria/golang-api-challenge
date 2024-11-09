package services

import (
	"github.com/AntonioDaria/surfe/src/models"
	"github.com/AntonioDaria/surfe/src/repository"
)

type Service interface {
	GetUserByID(userID int) (*models.User, error)
}

type ServiceImpl struct {
	userRepo repository.Repository
}

func NewUserService(userRepo repository.Repository) *ServiceImpl {
	return &ServiceImpl{userRepo: userRepo}
}

// GetUserByID retrieves a user by ID through the repository
func (s *ServiceImpl) GetUserByID(userID int) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}
