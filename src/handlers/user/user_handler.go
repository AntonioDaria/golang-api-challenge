package user

import (
	"strconv"

	"github.com/AntonioDaria/surfe/src/handlers/utils"
	not_found_err "github.com/AntonioDaria/surfe/src/repository"
	"github.com/gofiber/fiber/v2"
)

type UserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

// GetUserByIDHandler handles requests to retrieve a user by ID
func (h *Handler) GetUserByIDHandler(c *fiber.Ctx) error {
	// Parse user ID from the request parameters
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse user ID")
		return utils.JsonError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Retrieve the user using the service layer
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if err == not_found_err.ErrUserNotFound {
			h.logger.Error().Err(err).Msg("User not found")
			return utils.JsonError(c, fiber.StatusNotFound, "User not found")
		}
		h.logger.Error().Err(err).Msg("Failed to retrieve user")
		return utils.JsonError(c, fiber.StatusInternalServerError, "Failed to retrieve user")
	}

	// Return the user as JSON if found
	return c.JSON(UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
	})
}
