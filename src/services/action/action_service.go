package services

import "github.com/AntonioDaria/surfe/src/repository/action"

//go:generate mockgen -source=$GOFILE -destination=mock/action_service_mock.go -package=mock
type Service interface {
	CountActionsByUserID(userID int) int
}

type ServiceImpl struct {
	actionRepo action.Repository
}

func NewActionService(actionRepo action.Repository) *ServiceImpl {
	return &ServiceImpl{actionRepo: actionRepo}
}

// CountActionsByUserID counts the number of actions performed by a user
func (s *ServiceImpl) CountActionsByUserID(userID int) int {
	return s.actionRepo.CountActionsByUserID(userID)
}
