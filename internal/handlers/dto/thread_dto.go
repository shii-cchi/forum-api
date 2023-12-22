package dto

import "github.com/google/uuid"

type ThreadDto struct {
	Name    string    `json:"name"`
	ThemeID uuid.UUID `json:"themeID"`
}
