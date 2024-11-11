package action

import (
	"errors"
	"strconv"

	"github.com/AntonioDaria/surfe/src/handlers/utils"
	"github.com/AntonioDaria/surfe/src/models"
	"github.com/AntonioDaria/surfe/src/repository/action"
	"github.com/gofiber/fiber/v2"
)

type ActionCountResponse struct {
	Count int `json:"count"`
}

// GetActionCountByUserIDHandler retrieves the total count of actions for a given user ID
func (h *Handler) GetActionCountByUserIDHandler(c *fiber.Ctx) error {
	// Parse user ID from the request parameters
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse user ID")
		return utils.JsonError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Retrieve the action count using the service layer
	count, err := h.actionService.GetActionCountByUserID(userID)
	if err != nil {
		if errors.Is(err, action.ErrUserNotFound) {
			h.logger.Error().Err(err).Msg("User not found")
			return utils.JsonError(c, fiber.StatusNotFound, "User not found")
		}

		h.logger.Error().Err(err).Msg("Failed to retrieve action count")
		return utils.JsonError(c, fiber.StatusInternalServerError, "Failed to retrieve action count")
	}

	// Return the count as JSON if found
	return c.JSON(ActionCountResponse{Count: count})
}

type NextActionProbabilitiesResponse struct {
	Probabilities map[models.ActionType]float64 `json:"probabilities"`
}

func (h *Handler) GetNextActionProbabilitiesHandler(c *fiber.Ctx) error {
	actionType := models.ActionType(c.Params("actionType"))

	probabilities := h.actionService.GetNextActionProbabilities(actionType)

	return c.JSON(NextActionProbabilitiesResponse{Probabilities: probabilities})
}
