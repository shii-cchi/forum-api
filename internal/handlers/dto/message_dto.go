package dto

import "github.com/google/uuid"

type MessageDto struct {
	Name     string    `json:"name"`
	ThreadID uuid.UUID `json:"threadID"`
}
