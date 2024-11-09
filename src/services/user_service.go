package services

import (
	"github.com/AntonioDaria/surfe/src/models"
	"github.com/AntonioDaria/surfe/src/repository/user"
)

//go:generate mockgen -source=$GOFILE -destination=mock/services_mock.go -package=mock

type Service interface {
	GetUserByID(userID int) (*models.User, error)
}

type ServiceImpl struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) *ServiceImpl {
	return &ServiceImpl{userRepo: userRepo}
}

// GetUserByID retrieves a user by ID through the repository
func (s *ServiceImpl) GetUserByID(userID int) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}
