package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"net/http"
)

func (h *Handler) fetchThreads(w http.ResponseWriter, r *http.Request) {
	themeIdStr := r.URL.Query().Get("theme_id")

	themeId, err := uuid.Parse(themeIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't parse str to uuid: %v", err))
		return
	}

	threadList, err := h.services.Threads.GetThreadList(r.Context(), themeId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get threads: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, threadList)
}

func (h *Handler) createThread(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	newThread := new(dto.ThreadDto)
	err := json.NewDecoder(r.Body).Decode(&newThread)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	thread, err := h.services.Threads.CreateThread(r.Context(), newThread.Name, newThread.ThemeID, accessToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create thread: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, thread)
}

func (h *Handler) deleteThread(w http.ResponseWriter, r *http.Request) {
	threadId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get threadId: %v", err))
		return
	}

	deletedThread, err := h.services.Threads.GetThread(r.Context(), threadId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this thread: %v", err))
		return
	}

	err = h.services.Threads.DeleteThread(r.Context(), threadId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't delete this thread: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, deletedThread)
}

func (h *Handler) updateThread(w http.ResponseWriter, r *http.Request) {
	threadId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get threadId: %v", err))
		return
	}

	oldThread, err := h.services.Threads.GetThread(r.Context(), threadId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this thread: %v", err))
		return
	}

	newThreadParams := new(dto.ThreadDto)
	err = json.NewDecoder(r.Body).Decode(&newThreadParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if newThreadParams.Name == oldThread.Name {
		respondWithError(w, http.StatusNoContent, fmt.Sprintf("No content: %v", err))
		return
	}

	updatedThread, err := h.services.Threads.UpdateThread(r.Context(), threadId, newThreadParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't update thread name: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, updatedThread)
}
