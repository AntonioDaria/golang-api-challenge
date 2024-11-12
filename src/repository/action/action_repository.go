package action

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/AntonioDaria/surfe/src/models"
)

var ErrUserNotFound = fmt.Errorf("user not found")

//go:generate mockgen -source=$GOFILE -destination=mock/action_repository_mock.go -package=mock
type Repository interface {
	CountActionsByUserID(userID int) int
	UserExists(userID int) bool
	GetSortedActions() []models.Action
	GetAllActions() []models.Action
}

type RepositoryImpl struct {
	Actions []models.Action
}

// NewActionRepo loads action data from a JSON file and initializes ActionRepo
func NewActionRepo(filePath string) (*RepositoryImpl, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read action file: %w", err)
	}

	var actions []models.Action
	if err := json.Unmarshal(file, &actions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal action data: %w", err)
	}

	return &RepositoryImpl{Actions: actions}, nil
}

// CountActionsByUserID counts the number of actions performed by a user
func (r *RepositoryImpl) CountActionsByUserID(userID int) int {
	count := 0
	for _, action := range r.Actions {
		if action.UserID == userID {
			count++
		}
	}
	return count
}

// UserExists checks if a user has performed any actions
func (r *RepositoryImpl) UserExists(userID int) bool {
	for _, action := range r.Actions {
		if action.UserID == userID {
			return true
		}
	}
	return false
}

// GetSortedActions returns all actions sorted by user and timestamp.
// This allows to analyze the sequence of actions by user.
func (r *RepositoryImpl) GetSortedActions() []models.Action {
	sortedActions := make([]models.Action, len(r.Actions))
	copy(sortedActions, r.Actions)

	// Sort actions by UserID and then by Timestamp within each UserID
	sort.Slice(sortedActions, func(i, j int) bool {
		if sortedActions[i].UserID == sortedActions[j].UserID {
			return sortedActions[i].CreatedAt.Before(sortedActions[j].CreatedAt)
		}
		return sortedActions[i].UserID < sortedActions[j].UserID
	})

	return sortedActions
}

// GetAllActions returns all actions
func (r *RepositoryImpl) GetAllActions() []models.Action {
	return r.Actions
}
