package user

import (
	user_s "github.com/AntonioDaria/surfe/src/services/user"
	"github.com/rs/zerolog"
)

type Handler struct {
	userService user_s.Service
	logger      zerolog.Logger
}

func NewHandler(userService user_s.Service, logger zerolog.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}
