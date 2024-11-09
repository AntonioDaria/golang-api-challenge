package services

import (
	"github.com/AntonioDaria/surfe/src/repository/action"
)

//go:generate mockgen -source=$GOFILE -destination=mock/action_service_mock.go -package=mock
type Service interface {
	GetActionCountByUserID(userID int) (int, error)
}

type ServiceImpl struct {
	actionRepo action.Repository
}

func NewActionService(actionRepo action.Repository) *ServiceImpl {
	return &ServiceImpl{actionRepo: actionRepo}
}

// CountActionsByUserID counts the number of actions performed by a user
func (s *ServiceImpl) GetActionCountByUserID(userID int) (int, error) {
	if !s.actionRepo.UserExists(userID) {
		return 0, action.ErrUserNotFound
	}

	// Get the action count if the user exists
	return s.actionRepo.CountActionsByUserID(userID), nil
}
