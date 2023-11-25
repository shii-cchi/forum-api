package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (h *Handler) RegisterHTTPEndpoints(r chi.Router) {
	r.Mount("/auth", h.authHandlers())
}

func (h *Handler) authHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Post("/register", h.registerHandler)
		r.Post("/login", h.loginHandler)
		r.Get("/refresh", h.refreshHandler)
		r.Get("/logout", h.logoutHandler)
	})

	return rg
}
