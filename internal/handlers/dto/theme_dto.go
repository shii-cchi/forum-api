package dto

import "github.com/google/uuid"

type ThemeDto struct {
	Name         string    `json:"name"`
	SubsectionID uuid.UUID `json:"sectionID"`
}
