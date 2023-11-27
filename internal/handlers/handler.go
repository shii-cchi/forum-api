package handlers

import (
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/service"
	"github.com/shii-cchi/forum-api/pkg/hash"
)

type Handler struct {
	userService *service.UserService
}

func New(queries *database.Queries, hasher *hash.SHA1Hasher) *Handler {
	return &Handler{
		userService: service.NewUserService(queries, hasher),
	}
}
