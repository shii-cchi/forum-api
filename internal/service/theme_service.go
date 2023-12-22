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

type ThemeService struct {
	queries *database.Queries
	cfg     *config.Config

	tokenService Tokens
}

func NewThemeService(q *database.Queries, c *config.Config, t *TokenService) *ThemeService {
	return &ThemeService{
		queries:      q,
		cfg:          c,
		tokenService: t,
	}
}

func (s ThemeService) GetTheme(ctx context.Context, themeId uuid.UUID) (database.Theme, error) {
	theme, err := s.queries.GetTheme(ctx, themeId)

	if err != nil {
		return database.Theme{}, err
	}

	if reflect.DeepEqual(theme, database.Theme{}) {
		return database.Theme{}, errors.New("No theme with this id")
	}

	return theme, nil
}

func (s ThemeService) GetThemeList(ctx context.Context, subsectionId uuid.UUID) ([]database.Theme, error) {
	themeList, err := s.queries.GetThemes(ctx, subsectionId)

	if err != nil {
		return nil, err
	}

	return themeList, nil
}

func (s ThemeService) CreateTheme(ctx context.Context, name string, subsectionId uuid.UUID, token string) (database.Theme, error) {
	authorId, err := s.tokenService.GetIdFromToken(s.cfg.AccessSigningKey, token)

	if err != nil {
		return database.Theme{}, err
	}

	theme, err := s.queries.CreateTheme(ctx, database.CreateThemeParams{
		Name:         name,
		SubsectionID: subsectionId,
		AuthorID:     authorId,
	})

	if err != nil {
		return database.Theme{}, err
	}

	return theme, nil
}

func (s ThemeService) DeleteTheme(ctx context.Context, themeId uuid.UUID) error {
	err := s.queries.DeleteTheme(ctx, themeId)

	if err != nil {
		return err
	}

	return nil
}

func (s ThemeService) UpdateTheme(ctx context.Context, themeId uuid.UUID, newThemeParams *dto.ThemeDto) (database.Theme, error) {
	updatedTheme, err := s.queries.UpdateTheme(ctx, database.UpdateThemeParams{
		ID:           themeId,
		Name:         newThemeParams.Name,
		SubsectionID: newThemeParams.SubsectionID,
	})

	if err != nil {
		return database.Theme{}, err
	}

	return updatedTheme, nil
}
