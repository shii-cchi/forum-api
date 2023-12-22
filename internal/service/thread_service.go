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

type ThreadService struct {
	queries *database.Queries
	cfg     *config.Config

	tokenService Tokens
}

func NewThreadService(q *database.Queries, c *config.Config, t *TokenService) *ThreadService {
	return &ThreadService{
		queries:      q,
		cfg:          c,
		tokenService: t,
	}
}

func (s ThreadService) GetThread(ctx context.Context, threadId uuid.UUID) (database.Thread, error) {
	thread, err := s.queries.GetThread(ctx, threadId)

	if err != nil {
		return database.Thread{}, err
	}

	if reflect.DeepEqual(thread, database.Thread{}) {
		return database.Thread{}, errors.New("No thread with this id")
	}

	return thread, nil
}

func (s ThreadService) GetThreadList(ctx context.Context, themeId uuid.UUID) ([]database.Thread, error) {
	threadList, err := s.queries.GetThreads(ctx, themeId)

	if err != nil {
		return nil, err
	}

	return threadList, nil
}

func (s ThreadService) CreateThread(ctx context.Context, name string, themeId uuid.UUID, token string) (database.Thread, error) {
	authorId, err := s.tokenService.GetIdFromToken(s.cfg.AccessSigningKey, token)

	if err != nil {
		return database.Thread{}, err
	}

	thread, err := s.queries.CreateThread(ctx, database.CreateThreadParams{
		Name:     name,
		ThemeID:  themeId,
		AuthorID: authorId,
	})

	if err != nil {
		return database.Thread{}, err
	}

	return thread, nil
}

func (s ThreadService) DeleteThread(ctx context.Context, threadId uuid.UUID) error {
	err := s.queries.DeleteThread(ctx, threadId)

	if err != nil {
		return err
	}

	return nil
}

func (s ThreadService) UpdateThread(ctx context.Context, threadId uuid.UUID, newThreadParams *dto.ThreadDto) (database.Thread, error) {
	updatedThread, err := s.queries.UpdateThread(ctx, database.UpdateThreadParams{
		ID:      threadId,
		Name:    newThreadParams.Name,
		ThemeID: newThreadParams.ThemeID,
	})

	if err != nil {
		return database.Thread{}, err
	}

	return updatedThread, nil
}
