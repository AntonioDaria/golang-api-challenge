package services

import (
	"testing"

	"github.com/AntonioDaria/surfe/src/repository/action/mock"
	"github.com/golang/mock/gomock"
)

func Test_CountActionsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new mock repository
	actionRepo := mock.NewMockRepository(ctrl)
	actionService := NewActionService(actionRepo)

	// Define expected behavior for CountActionsByUserID
	actionRepo.EXPECT().CountActionsByUserID(1).Return(2)

	// Act
	count := actionService.CountActionsByUserID(1)

	// Assert
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}
