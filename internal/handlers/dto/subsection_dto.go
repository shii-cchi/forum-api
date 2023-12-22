package dto

import "github.com/google/uuid"

type SubsectionDto struct {
	Name      string    `json:"name"`
	SectionID uuid.UUID `json:"sectionID"`
}
