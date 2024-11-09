package services

import (
	"github.com/AntonioDaria/surfe/src/models"

	"github.com/AntonioDaria/surfe/src/repository/mock"

	"testing"

	"github.com/AntonioDaria/surfe/src/repository/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new mock repository
	userRepo := mock.NewMockRepository(ctrl)
	userService := NewUserService(userRepo)

	// Define expected behavior for GetUserByID
	expectedUser := &models.User{ID: 1, Name: "Ferdinande"}
	userRepo.EXPECT().GetUserByID(1).Return(expectedUser, nil)

	// Act
	user, err := userService.GetUserByID(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

// test user not found
func TestGetUserByID_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new mock repository
	userRepo := mock.NewMockRepository(ctrl)
	userService := NewUserService(userRepo)

	// Define expected behavior for GetUserByID
	userRepo.EXPECT().GetUserByID(1).Return(nil, user.ErrUserNotFound)

	// Act
	found_user, err := userService.GetUserByID(1)

	// Assert
	assert.ErrorIs(t, err, user.ErrUserNotFound)
	assert.Nil(t, found_user)
}
