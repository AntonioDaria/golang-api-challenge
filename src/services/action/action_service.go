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
	GetReferralIndex() map[int]int
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

	// Iterate through sorted actions to find occurrences of actionType
	for i := 0; i < len(sortedActions); i++ {
		currentAction := sortedActions[i]

		// If current action matches actionType, count subsequent actions by the same user
		if currentAction.Type == actionType {
			for j := i + 1; j < len(sortedActions); j++ {
				nextAction := sortedActions[j]

				// Stop counting if we encounter another instance of actionType or a different user
				if nextAction.UserID != currentAction.UserID || nextAction.Type == actionType {
					break
				}

				// Count this next action
				nextActionCounts[nextAction.Type]++
				totalCount++
			}
		}
	}

	// Calculate probabilities by dividing each next action count by the total count
	probabilities := make(map[act_type.ActionType]float64)
	for action, count := range nextActionCounts {
		probability := float64(count) / float64(totalCount)
		probabilities[action] = math.Round(probability*100) / 100 // rounds to 2 decimal places
	}

	return probabilities
}

func (s *ServiceImpl) GetReferralIndex() map[int]int {
	// Build an adjacency list from refer actions
	adjacencyList := make(map[int][]int)
	for _, action := range s.actionRepo.GetAllActions() {
		if action.Type == act_type.ActionTypeReferUser {
			referrer := action.UserID
			referred := action.TargetUser
			adjacencyList[referrer] = append(adjacencyList[referrer], referred)
		}
	}

	// Final referral index map to store results for each user
	referralIndex := make(map[int]int)
	visited := make(map[int]bool)
	inCycle := make(map[int]bool)

	// DFS function to calculate referral index with cycle detection
	var dfs func(userID int, path map[int]bool) int
	dfs = func(userID int, path map[int]bool) int {
		if path[userID] {
			// If a node is revisited in the same path, we have a cycle
			fmt.Printf("Cycle detected for User %d\n", userID)
			for node := range path {
				inCycle[node] = true // Mark all nodes in the path as part of a cycle
			}
			return 0
		}

		if visited[userID] {
			// If already visited and calculated, return the cached result
			return referralIndex[userID]
		}

		// Mark this user as visited in the current path
		path[userID] = true
		visited[userID] = true
		count := 0

		// Traverse each referral
		for _, referred := range adjacencyList[userID] {
			count += 1 + dfs(referred, path)
		}

		// Unmark this user from the current path
		path[userID] = false
		referralIndex[userID] = count
		return count
	}

	// Calculate the referral index for each user
	for userID := range adjacencyList {
		if !visited[userID] {
			dfs(userID, make(map[int]bool))
		}
	}

	// Adjust referral indices for nodes in cycles
	for node := range inCycle {
		referralIndex[node] = len(inCycle) - 1
	}

	return referralIndex
}
