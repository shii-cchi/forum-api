package handlers

import (
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/service"
)

type Handler struct {
	services *service.Services
	cfg      *config.Config
}

func NewHandler(services *service.Services, c *config.Config) *Handler {
	return &Handler{
		services: services,
		cfg:      c,
	}
}
