package handlers

import (
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/service"
	"github.com/shii-cchi/forum-api/pkg/hash"
)

type Handler struct {
	userService *service.UserService
	cfg         *config.Config
}

func New(queries *database.Queries, hasher *hash.SHA1Hasher, cfg *config.Config) *Handler {
	return &Handler{
		userService: service.NewUserService(queries, hasher, cfg),
		cfg:         cfg,
	}
}
