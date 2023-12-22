package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"net/http"
)

func (h *Handler) fetchSubsections(w http.ResponseWriter, r *http.Request) {
	sectionIdStr := r.URL.Query().Get("section_id")

	sectionId, err := uuid.Parse(sectionIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't parse str to uuid: %v", err))
		return
	}

	subsectionList, err := h.services.Subsections.GetSubsectionList(r.Context(), sectionId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get subsections: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, subsectionList)
}

func (h *Handler) createSubsection(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	newSubsection := new(dto.SubsectionDto)
	err := json.NewDecoder(r.Body).Decode(&newSubsection)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	subsection, err := h.services.Subsections.CreateSubsection(r.Context(), newSubsection.Name, newSubsection.SectionID, accessToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create subsection: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, subsection)
}

func (h *Handler) deleteSubsection(w http.ResponseWriter, r *http.Request) {
	subsectionId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get subsectionId: %v", err))
		return
	}

	deletedSubsection, err := h.services.Subsections.GetSubsection(r.Context(), subsectionId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this subsection: %v", err))
		return
	}

	err = h.services.Subsections.DeleteSubsection(r.Context(), subsectionId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't delete this subsection: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, deletedSubsection)
}

func (h *Handler) updateSubsection(w http.ResponseWriter, r *http.Request) {
	subsectionId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get subsectionId: %v", err))
		return
	}

	oldSubsection, err := h.services.Subsections.GetSubsection(r.Context(), subsectionId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this subsection: %v", err))
		return
	}

	newSubsectionParams := new(dto.SubsectionDto)
	err = json.NewDecoder(r.Body).Decode(&newSubsectionParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if newSubsectionParams.Name == oldSubsection.Name {
		respondWithError(w, http.StatusNoContent, fmt.Sprintf("No content: %v", err))
		return
	}

	updatedSubsection, err := h.services.Subsections.UpdateSubsection(r.Context(), subsectionId, newSubsectionParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't update subsection name: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, updatedSubsection)
}
