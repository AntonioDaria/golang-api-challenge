package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/AntonioDaria/surfe/src/models"
	"github.com/AntonioDaria/surfe/src/services"
	"github.com/AntonioDaria/surfe/src/services/mock"
	"github.com/rs/zerolog"

	"github.com/AntonioDaria/surfe/src/repository/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByIDHandler_Success(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Set up mock service
	mockService := mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	// Set up test data
	mockUser := &models.User{
		ID:        1,
		Name:      "John Doe",
		CreatedAt: time.Now(),
	}
	mockService.EXPECT().GetUserByID(1).Return(mockUser, nil)

	// Create a new Fiber app and test request
	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetUserByIDHandler_NotFound(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	// Simulate user not found error
	mockService.EXPECT().GetUserByID(2).Return(nil, user.ErrUserNotFound)

	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetUserByIDHandler_BadRequest(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetUserByIDHandler_InternalServerError(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	// Simulate internal server error
	mockService.EXPECT().GetUserByID(3).Return(nil, errors.New("internal server error"))

	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/users/3", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetUserByIDIntegration(t *testing.T) {
	// Initialize the repository with test data
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	userRepo, err := user.NewUserRepo("../../repository/data/users.json")
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Set up service and handler
	userService := services.NewUserService(userRepo)
	userHandler := NewHandler(userService, logger)

	// Create a new Fiber app and register the route
	app := fiber.New()
	app.Get("/users/:id", userHandler.GetUserByIDHandler)

	// Happy Path: Test retrieving a user that exists
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check response content:
	var user UserResponse
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Ferdinande", user.Name)
	assert.Equal(t, "2020-07-14T05:48:54.798Z", user.CreatedAt)

	// Not Found: Test retrieving a user that does not exist
	req = httptest.NewRequest(http.MethodGet, "/users/9999", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Bad Request: Test with invalid user ID
	req = httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
