package services

import (
	"fmt"
	"math"

	act_type "github.com/AntonioDaria/surfe/src/models"
	"github.com/AntonioDaria/surfe/src/repository/action"
)

//go:generate mockgen -source=$GOFILE -destination=mock/action_service_mock.go -package=mock
type Service interface {
	GetActionCountByUserID(userID int) (int, error)
	GetNextActionProbabilities(actionType act_type.ActionType) map[act_type.ActionType]float64
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

func (s *ServiceImpl) GetNextActionProbabilities(actionType act_type.ActionType) map[act_type.ActionType]float64 {
	sortedActions := s.actionRepo.GetSortedActions()
	nextActionCounts := make(map[act_type.ActionType]int)
	totalCount := 0

	// Iterate through sorted actions to find each occurrence of actionType
	for i := 0; i < len(sortedActions)-1; i++ {
		currentAction := sortedActions[i]
		nextAction := sortedActions[i+1]

		// Check if the current action matches the specified action type and if the next action is from the same user
		if currentAction.Type == actionType && currentAction.UserID == nextAction.UserID {
			// Count this next action
			nextActionCounts[nextAction.Type]++
			totalCount++
		}
	}

	// Debugging output to verify next action counts and total count
	fmt.Printf("Final Next Action Counts: %v, Total Count: %d\n", nextActionCounts, totalCount)

	// Calculate probabilities by dividing each next action count by the total count of all next actions
	probabilities := make(map[act_type.ActionType]float64)
	for action, count := range nextActionCounts {
		probability := float64(count) / float64(totalCount)
		probabilities[action] = math.Round(probability*100) / 100 // rounds to 2 decimal places
	}

	// Debugging output to check probabilities
	fmt.Printf("Probabilities: %v\n", probabilities)

	return probabilities
}
