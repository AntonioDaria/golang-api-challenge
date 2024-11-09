package action

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"encoding/json"

	"github.com/AntonioDaria/surfe/src/repository/action"
	action_s "github.com/AntonioDaria/surfe/src/services/action"
	action_mock "github.com/AntonioDaria/surfe/src/services/action/mock"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetActionCountByUserIDHandler_Success(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Set up mock service
	mockService := action_mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	// Define the expected behavior and result
	mockService.EXPECT().GetActionCountByUserID(1).Return(100, nil)

	// Set up the Fiber app
	app := fiber.New()
	app.Get("/users/:id/actions/count", handler.GetActionCountByUserIDHandler)

	// Perform the request
	req := httptest.NewRequest(http.MethodGet, "/users/1/actions/count", nil)
	resp, _ := app.Test(req, -1)

	// Assert the status and response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetActionCountByUserIDHandler_User_Not_Found(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Set up mock service
	mockService := action_mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	// Define the expected behavior and result
	mockService.EXPECT().GetActionCountByUserID(1).Return(0, action.ErrUserNotFound)

	// Set up the Fiber app
	app := fiber.New()
	app.Get("/users/:id/actions/count", handler.GetActionCountByUserIDHandler)

	// Perform the request
	req := httptest.NewRequest(http.MethodGet, "/users/1/actions/count", nil)
	resp, _ := app.Test(req, -1)

	// Assert the status and response
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetActionCountByUserIDHandler_Invalid_User_ID(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Set up mock service
	mockService := action_mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	// Set up the Fiber app
	app := fiber.New()
	app.Get("/users/:id/actions/count", handler.GetActionCountByUserIDHandler)

	// Perform the request
	req := httptest.NewRequest(http.MethodGet, "/users/invalid/actions/count", nil)
	resp, _ := app.Test(req, -1)

	// Assert the status and response
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetActionCountByUserIDHandler_Internal_Server_Error(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Set up mock service
	mockService := action_mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	// Define the expected behavior and result
	mockService.EXPECT().GetActionCountByUserID(1).Return(0, errors.New("internal server error"))

	// Set up the Fiber app
	app := fiber.New()
	app.Get("/users/:id/actions/count", handler.GetActionCountByUserIDHandler)

	// Perform the request
	req := httptest.NewRequest(http.MethodGet, "/users/1/actions/count", nil)
	resp, _ := app.Test(req, -1)

	// Assert the status and response
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetActionCountByUserIDIntegration(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	actionRepo, err := action.NewActionRepo("../../repository/data/actions.json")
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	actionService := action_s.NewActionService(actionRepo)
	handler := NewHandler(actionService, logger)

	app := fiber.New()
	app.Get("/users/:id/actions/count", handler.GetActionCountByUserIDHandler)

	// Happy path
	req := httptest.NewRequest(http.MethodGet, "/users/1/actions/count", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// check action count response for user 1
	var countResponse ActionCountResponse
	err = json.NewDecoder(resp.Body).Decode(&countResponse)
	assert.NoError(t, err)

	assert.Equal(t, 49, countResponse.Count)

	// User not found
	req = httptest.NewRequest(http.MethodGet, "/users/1000/actions/count", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Invalid user ID 404
	req = httptest.NewRequest(http.MethodGet, "/users/invalid/actions/count", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
