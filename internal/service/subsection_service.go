package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"reflect"
)

type SubsectionService struct {
	queries *database.Queries
	cfg     *config.Config

	tokenService Tokens
}

func NewSubsectionService(q *database.Queries, c *config.Config, t *TokenService) *SubsectionService {
	return &SubsectionService{
		queries:      q,
		cfg:          c,
		tokenService: t,
	}
}

func (s SubsectionService) GetSubsection(ctx context.Context, subsectionId uuid.UUID) (database.Subsection, error) {
	subsection, err := s.queries.GetSubsection(ctx, subsectionId)

	if err != nil {
		return database.Subsection{}, err
	}

	if reflect.DeepEqual(subsection, database.Subsection{}) {
		return database.Subsection{}, errors.New("No subsection with this id")
	}

	return subsection, nil
}

func (s SubsectionService) GetSubsectionList(ctx context.Context, sectionId uuid.UUID) ([]database.Subsection, error) {
	subsectionList, err := s.queries.GetSubsections(ctx, sectionId)

	if err != nil {
		return nil, err
	}

	return subsectionList, nil
}

func (s SubsectionService) CreateSubsection(ctx context.Context, name string, sectionId uuid.UUID, token string) (database.Subsection, error) {
	authorId, err := s.tokenService.GetIdFromToken(s.cfg.AccessSigningKey, token)

	if err != nil {
		return database.Subsection{}, err
	}

	subsection, err := s.queries.CreateSubsection(ctx, database.CreateSubsectionParams{
		Name:      name,
		SectionID: sectionId,
		AuthorID:  authorId,
	})

	if err != nil {
		return database.Subsection{}, err
	}

	return subsection, nil
}

func (s SubsectionService) DeleteSubsection(ctx context.Context, subsectionId uuid.UUID) error {
	err := s.queries.DeleteSubsection(ctx, subsectionId)

	if err != nil {
		return err
	}

	return nil
}

func (s SubsectionService) UpdateSubsection(ctx context.Context, subsectionId uuid.UUID, newSubsectionParams *dto.SubsectionDto) (database.Subsection, error) {
	updatedSubsection, err := s.queries.UpdateSubsection(ctx, database.UpdateSubsectionParams{
		ID:        subsectionId,
		Name:      newSubsectionParams.Name,
		SectionID: newSubsectionParams.SectionID,
	})

	if err != nil {
		return database.Subsection{}, err
	}

	return updatedSubsection, nil
}
