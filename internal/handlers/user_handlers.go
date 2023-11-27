package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"net/http"
	"time"
)

func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	newUser := new(dto.UserDto)
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, refreshToken, err := h.userService.CreateUser(r.Context(), newUser)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create user: %v", err))
		return
	}

	cookie := http.Cookie{
		Name:    "Token",
		Value:   refreshToken,
		Expires: time.Now().Add(720 * time.Hour),
		Path:    "/auth/register",
	}

	http.SetCookie(w, &cookie)

	respondWithJSON(w, http.StatusCreated, user)
}

func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) refreshHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {

}
