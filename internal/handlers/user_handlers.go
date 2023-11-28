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
	checkedUser := new(dto.UserDto)
	err := json.NewDecoder(r.Body).Decode(&checkedUser)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, refreshToken, err := h.userService.Login(r.Context(), checkedUser)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't login: %v", err))
		return
	}

	cookie := http.Cookie{
		Name:    "Token",
		Value:   refreshToken,
		Expires: time.Now().Add(720 * time.Hour),
		Path:    "/auth/login",
	}

	http.SetCookie(w, &cookie)

	respondWithJSON(w, http.StatusOK, user)
}

func (h *Handler) refreshHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	if accessToken == "" {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("Forbidden"))
		return
	}

	err := h.userService.Logout(r.Context(), accessToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't logout: %v", err))
		return
	}

	cookie := http.Cookie{
		Name:    "Token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
		Path:    "/auth/register",
	}

	http.SetCookie(w, &cookie)
}
