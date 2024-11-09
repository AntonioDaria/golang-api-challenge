package services

import (
	"testing"

	"github.com/AntonioDaria/surfe/src/repository/action"
	"github.com/AntonioDaria/surfe/src/repository/action/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_CountActionsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new mock repository
	actionRepo := mock.NewMockRepository(ctrl)
	actionService := NewActionService(actionRepo)

	// Define expected behavior for UserExists
	actionRepo.EXPECT().UserExists(1).Return(true)

	// Define expected behavior for CountActionsByUserID
	actionRepo.EXPECT().CountActionsByUserID(1).Return(2)

	// Act
	count, err := actionService.GetActionCountByUserID(1)
	assert.NoError(t, err)

	// Assert
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func Test_CountActionsByUserID_User_Not_Found(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new mock repository
	actionRepo := mock.NewMockRepository(ctrl)
	actionService := NewActionService(actionRepo)

	// Define expected behavior for UserExists
	actionRepo.EXPECT().UserExists(1).Return(false)

	// Act
	count, err := actionService.GetActionCountByUserID(1)
	assert.Error(t, err)
	assert.Equal(t, action.ErrUserNotFound, err)
	assert.Equal(t, 0, count)
}
