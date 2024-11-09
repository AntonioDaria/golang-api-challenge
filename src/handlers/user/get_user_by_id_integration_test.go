package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AntonioDaria/surfe/src/repository"
	"github.com/AntonioDaria/surfe/src/services"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByIDIntegration(t *testing.T) {
	// Initialize the repository with test data
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	userRepo, err := repository.NewUserRepo("../../repository/data/users.json")
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
