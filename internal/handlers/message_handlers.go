package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"net/http"
)

func (h *Handler) fetchMessages(w http.ResponseWriter, r *http.Request) {
	threadIdStr := r.URL.Query().Get("thread_id")

	threadId, err := uuid.Parse(threadIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't parse str to uuid: %v", err))
		return
	}

	messageList, err := h.services.Messages.GetMessageList(r.Context(), threadId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get messages: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, messageList)
}

func (h *Handler) createMessage(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	newMessage := new(dto.MessageDto)
	err := json.NewDecoder(r.Body).Decode(&newMessage)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	message, err := h.services.Messages.CreateMessage(r.Context(), newMessage.Name, newMessage.ThreadID, accessToken)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create message: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, message)
}

func (h *Handler) deleteMessage(w http.ResponseWriter, r *http.Request) {
	messageId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get messageId: %v", err))
		return
	}

	deletedMessage, err := h.services.Messages.GetMessage(r.Context(), messageId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this message: %v", err))
		return
	}

	err = h.services.Messages.DeleteMessage(r.Context(), messageId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't delete this message: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, deletedMessage)
}

func (h *Handler) updateMessage(w http.ResponseWriter, r *http.Request) {
	messageId, err := h.GetIdFromURL(w, r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't get messageId: %v", err))
		return
	}

	oldMessage, err := h.services.Messages.GetMessage(r.Context(), messageId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find this message: %v", err))
		return
	}

	newMessageParams := new(dto.MessageDto)
	err = json.NewDecoder(r.Body).Decode(&newMessageParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if newMessageParams.Name == oldMessage.Name {
		respondWithError(w, http.StatusNoContent, fmt.Sprintf("No content: %v", err))
		return
	}

	updatedMessage, err := h.services.Messages.UpdateMessage(r.Context(), messageId, newMessageParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't update message name: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, updatedMessage)
}
