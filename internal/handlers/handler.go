package handlers

import (
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/service"
)

type Handler struct {
	userService *service.UserService
}

func New(queries *database.Queries) *Handler {
	return &Handler{
		userService: service.NewUserService(queries),
	}
}
