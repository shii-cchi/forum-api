package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"net/http"
)

func (h *Handler) fetchSections(w http.ResponseWriter, r *http.Request) {
	sectionList, err := h.services.Sections.GetSectionList(r.Context())

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get sections: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, sectionList)
}

func (h *Handler) createSection(w http.ResponseWriter, r *http.Request) {
	newSection := new(dto.SectionDto)
	err := json.NewDecoder(r.Body).Decode(&newSection)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	section, err := h.services.Sections.CreateSection(r.Context(), newSection.Name)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create section: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, section)
}

func (h *Handler) deleteSection(w http.ResponseWriter, r *http.Request) {
	sectionId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get sectionId: %v", err))
		return
	}

	deletedSection, err := h.services.Sections.GetSection(r.Context(), sectionId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this section: %v", err))
		return
	}

	err = h.services.Sections.DeleteSection(r.Context(), sectionId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't delete this section: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, deletedSection)
}

func (h *Handler) updateSection(w http.ResponseWriter, r *http.Request) {
	sectionId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get sectionId: %v", err))
		return
	}

	oldSection, err := h.services.Sections.GetSection(r.Context(), sectionId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this section: %v", err))
		return
	}

	newSectionParams := new(dto.SectionDto)
	err = json.NewDecoder(r.Body).Decode(&newSectionParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if newSectionParams.Name == oldSection.Name {
		respondWithError(w, http.StatusNoContent, fmt.Sprintf("No content: %v", err))
		return
	}

	updatedSection, err := h.services.Sections.UpdateSectionName(r.Context(), sectionId, newSectionParams.Name)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't update section name: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, updatedSection)
}
