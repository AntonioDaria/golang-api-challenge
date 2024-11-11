package action

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"encoding/json"

	"github.com/AntonioDaria/surfe/src/models"
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

func TestGetNextActionProbabilitiesHandler(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Set up mock service
	mockService := action_mock.NewMockService(ctrl)
	handler := NewHandler(mockService, logger)

	mockReturn := map[models.ActionType]float64{
		models.ActionTypeAddContact:  0.5,
		models.ActionTypeEditContact: 0.5,
	}

	// Define the expected behavior and result
	mockService.EXPECT().GetNextActionProbabilities(models.ActionTypeAddContact).Return(mockReturn)
	app := fiber.New()
	app.Get("/actions/:actionType/probabilities", handler.GetNextActionProbabilitiesHandler)

	// Perform the request
	req := httptest.NewRequest(http.MethodGet, "/actions/ADD_CONTACT/probabilities", nil)
	resp, _ := app.Test(req, -1)

	// Assert the status and response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetNextActionProbabilitiesHandler_Precision(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	// Minimal dataset for testing precision
	actionRepo := &action.RepositoryImpl{
		Actions: []models.Action{
			{ID: 1, UserID: 1, Type: models.ActionTypeAddContact, CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
			{ID: 2, UserID: 1, Type: models.ActionTypeViewContacts, CreatedAt: time.Date(2021, 1, 1, 0, 0, 1, 0, time.UTC)},
			{ID: 3, UserID: 1, Type: models.ActionTypeEditContact, CreatedAt: time.Date(2021, 1, 1, 0, 0, 2, 0, time.UTC)},
			{ID: 4, UserID: 1, Type: models.ActionTypeReferUser, CreatedAt: time.Date(2021, 1, 1, 0, 0, 3, 0, time.UTC)},
		},
	}

	actionService := action_s.NewActionService(actionRepo)
	handler := NewHandler(actionService, logger)

	app := fiber.New()
	app.Get("/actions/:actionType/probabilities", handler.GetNextActionProbabilitiesHandler)

	// Send request to the endpoint
	req := httptest.NewRequest(http.MethodGet, "/actions/ADD_CONTACT/probabilities", nil)
	resp, _ := app.Test(req, -1)

	// Assert response status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response body
	var probabilitiesResponse NextActionProbabilitiesResponse
	err := json.NewDecoder(resp.Body).Decode(&probabilitiesResponse)
	assert.NoError(t, err)

	// Check that each probability is rounded to 2 decimal places
	for actionType, probability := range probabilitiesResponse.Probabilities {
		assert.InDelta(t, probability, math.Round(probability*100)/100, 0.001,
			"Probability for %v should be rounded to 2 decimal places but got %v", actionType, probability)
	}

	// Check the probability values for VIEW_CONTACTS, EDIT_CONTACT, and REFER_USER
	assert.Equal(t, 0.33, probabilitiesResponse.Probabilities[models.ActionTypeViewContacts])
	assert.Equal(t, 0.33, probabilitiesResponse.Probabilities[models.ActionTypeEditContact])
	assert.Equal(t, 0.33, probabilitiesResponse.Probabilities[models.ActionTypeReferUser])
}

func TestGetNextActionProbabilitiesHandler_Integration(t *testing.T) {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	// Minimal dataset for testing
	actionRepo := &action.RepositoryImpl{
		Actions: []models.Action{
			{ID: 1, UserID: 1, Type: models.ActionTypeAddContact, CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
			{ID: 2, UserID: 1, Type: models.ActionTypeViewContacts, CreatedAt: time.Date(2021, 1, 1, 0, 0, 1, 0, time.UTC)},
			{ID: 3, UserID: 1, Type: models.ActionTypeEditContact, CreatedAt: time.Date(2021, 1, 1, 0, 0, 2, 0, time.UTC)},
		},
	}

	actionService := action_s.NewActionService(actionRepo)
	handler := NewHandler(actionService, logger)

	app := fiber.New()
	app.Get("/actions/:actionType/probabilities", handler.GetNextActionProbabilitiesHandler)

	// Request setup
	req := httptest.NewRequest(http.MethodGet, "/actions/ADD_CONTACT/probabilities", nil)
	resp, _ := app.Test(req, -1)

	// Assert response status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response body
	var probabilitiesResponse NextActionProbabilitiesResponse
	err := json.NewDecoder(resp.Body).Decode(&probabilitiesResponse)
	assert.NoError(t, err)

	// Check the probability values for VIEW_CONTACTS and EDIT_CONTACT
	assert.Equal(t, 0.5, probabilitiesResponse.Probabilities[models.ActionTypeViewContacts])
	assert.Equal(t, 0.5, probabilitiesResponse.Probabilities[models.ActionTypeEditContact])
}
