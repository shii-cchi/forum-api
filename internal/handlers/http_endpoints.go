package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (h *Handler) RegisterHTTPEndpoints(r chi.Router) {
	r.Mount("/auth", h.authHandlers())
	r.Mount("/section", h.sectionHandlers())
	r.Mount("/subsection", h.subsectionHandlers())
	r.Mount("/theme", h.themeHandlers())
	r.Mount("/thread", h.threadHandlers())
	r.Mount("/message", h.messageHandlers())
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

func (h *Handler) sectionHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", h.fetchSections)
		r.Post("/", h.createSection)
		r.Delete("/{id}", h.deleteSection)
		r.Patch("/{id}", h.updateSection)
	})

	return rg
}

func (h *Handler) subsectionHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", h.fetchSubsections)
		r.Post("/", h.createSubsection)
		r.Delete("/{id}", h.deleteSubsection)
		r.Patch("/{id}", h.updateSubsection)
	})

	return rg
}

func (h *Handler) themeHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", h.fetchThemes)
		r.Post("/", h.createTheme)
		r.Delete("/{id}", h.deleteTheme)
		r.Patch("/{id}", h.updateTheme)
	})

	return rg
}

func (h *Handler) threadHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", h.fetchThreads)
		r.Post("/", h.createThread)
		r.Delete("/{id}", h.deleteThread)
		r.Patch("/{id}", h.updateThread)
	})

	return rg
}

func (h *Handler) messageHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", h.fetchMessages)
		r.Post("/", h.createMessage)
		r.Delete("/{id}", h.deleteMessage)
		r.Patch("/{id}", h.updateMessage)
	})

	return rg
}
