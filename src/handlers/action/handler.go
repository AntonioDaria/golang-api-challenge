package action

import (
	action_s "github.com/AntonioDaria/surfe/src/services/action"
	"github.com/rs/zerolog"
)

type Handler struct {
	actionService action_s.Service
	logger        zerolog.Logger
}

func NewHandler(actionService action_s.Service, logger zerolog.Logger) *Handler {
	return &Handler{
		actionService: actionService,
		logger:        logger,
	}
}
