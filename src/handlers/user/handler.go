package user

import (
	user "github.com/AntonioDaria/surfe/src/services"
	"github.com/rs/zerolog"
)

type Handler struct {
	userService user.Service
	logger      zerolog.Logger
}

func NewHandler(userService user.Service, logger zerolog.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}
