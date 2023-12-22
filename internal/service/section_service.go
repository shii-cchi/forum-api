package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/database"
	"reflect"
)

type SectionService struct {
	queries *database.Queries
}

func NewSectionService(q *database.Queries) *SectionService {
	return &SectionService{
		queries: q,
	}
}

func (s SectionService) GetSection(ctx context.Context, sectionId uuid.UUID) (database.Section, error) {
	section, err := s.queries.GetSection(ctx, sectionId)

	if err != nil {
		return database.Section{}, err
	}

	if reflect.DeepEqual(section, database.Section{}) {
		return database.Section{}, errors.New("No section with this id")
	}

	return section, nil
}

func (s SectionService) GetSectionList(ctx context.Context) ([]database.Section, error) {
	sectionList, err := s.queries.GetSections(ctx)

	if err != nil {
		return nil, err
	}

	return sectionList, nil
}

func (s SectionService) CreateSection(ctx context.Context, name string) (database.Section, error) {
	section, err := s.queries.CreateSection(ctx, name)

	if err != nil {
		return database.Section{}, err
	}

	return section, nil
}

func (s SectionService) DeleteSection(ctx context.Context, sectionId uuid.UUID) error {
	err := s.queries.DeleteSection(ctx, sectionId)

	if err != nil {
		return err
	}

	return nil
}

func (s SectionService) UpdateSectionName(ctx context.Context, sectionId uuid.UUID, newName string) (database.Section, error) {
	updatedSection, err := s.queries.UpdateSectionName(ctx, database.UpdateSectionNameParams{
		ID:   sectionId,
		Name: newName,
	})

	if err != nil {
		return database.Section{}, err
	}

	return updatedSection, nil
}
