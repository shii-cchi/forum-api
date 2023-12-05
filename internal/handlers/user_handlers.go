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
		Name:     "Token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
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
		Name:     "Token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	http.SetCookie(w, &cookie)

	respondWithJSON(w, http.StatusOK, user)
}

func (h *Handler) refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")

	if err != nil {
		if err == http.ErrNoCookie {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cookie not found", err))
			return
		}

		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return

	}

	refreshToken := cookie.Value

	ok, err := h.userService.IsValidToken(refreshToken, h.cfg.RefreshSigningKey)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if !ok {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unauthorized"))
		return
	}

	user, refreshToken, err := h.userService.Refresh(r.Context(), refreshToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	cookie = &http.Cookie{
		Name:     "Token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	http.SetCookie(w, cookie)

	respondWithJSON(w, http.StatusOK, user)
}

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	ok, err := h.userService.IsValidToken(accessToken, h.cfg.AccessSigningKey)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if !ok {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unauthorized"))
		return
	}

	err = h.userService.Logout(r.Context(), accessToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't logout: %v", err))
		return
	}

	cookie := http.Cookie{
		Name:     "Token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
	}

	http.SetCookie(w, &cookie)
}
