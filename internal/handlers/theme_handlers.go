package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"net/http"
)

func (h *Handler) fetchThemes(w http.ResponseWriter, r *http.Request) {
	subsectionIdStr := r.URL.Query().Get("subsection_id")

	subsectionId, err := uuid.Parse(subsectionIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't parse str to uuid: %v", err))
		return
	}

	themeList, err := h.services.Themes.GetThemeList(r.Context(), subsectionId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get themes: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, themeList)
}

func (h *Handler) createTheme(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	newTheme := new(dto.ThemeDto)
	err := json.NewDecoder(r.Body).Decode(&newTheme)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	theme, err := h.services.Themes.CreateTheme(r.Context(), newTheme.Name, newTheme.SubsectionID, accessToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create theme: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, theme)
}

func (h *Handler) deleteTheme(w http.ResponseWriter, r *http.Request) {
	themeId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get themeId: %v", err))
		return
	}

	deletedTheme, err := h.services.Themes.GetTheme(r.Context(), themeId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this theme: %v", err))
		return
	}

	err = h.services.Themes.DeleteTheme(r.Context(), themeId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't delete this theme: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, deletedTheme)
}

func (h *Handler) updateTheme(w http.ResponseWriter, r *http.Request) {
	themeId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get themeId: %v", err))
		return
	}

	oldTheme, err := h.services.Themes.GetTheme(r.Context(), themeId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this theme: %v", err))
		return
	}

	newThemeParams := new(dto.ThemeDto)
	err = json.NewDecoder(r.Body).Decode(&newThemeParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if newThemeParams.Name == oldTheme.Name {
		respondWithError(w, http.StatusNoContent, fmt.Sprintf("No content: %v", err))
		return
	}

	updatedTheme, err := h.services.Themes.UpdateTheme(r.Context(), themeId, newThemeParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't update theme name: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, updatedTheme)
}
