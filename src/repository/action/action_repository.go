package action

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AntonioDaria/surfe/src/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock/action_repository_mock.go -package=mock
type Repository interface {
	CountActionsByUserID(userID int) int
}

type RepositoryImpl struct {
	actions []models.Action
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

	return &RepositoryImpl{actions: actions}, nil
}

// CountActionsByUserID counts the number of actions performed by a user
func (r *RepositoryImpl) CountActionsByUserID(userID int) int {
	count := 0
	for _, action := range r.actions {
		if action.UserID == userID {
			count++
		}
	}
	return count
}
