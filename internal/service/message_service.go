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

type MessageService struct {
	queries *database.Queries
	cfg     *config.Config

	tokenService Tokens
}

func NewMessageService(q *database.Queries, c *config.Config, t *TokenService) *MessageService {
	return &MessageService{
		queries:      q,
		cfg:          c,
		tokenService: t,
	}
}

func (s MessageService) GetMessage(ctx context.Context, messageId uuid.UUID) (database.Message, error) {
	message, err := s.queries.GetMessage(ctx, messageId)

	if err != nil {
		return database.Message{}, err
	}

	if reflect.DeepEqual(message, database.Message{}) {
		return database.Message{}, errors.New("No message with this id")
	}

	return message, nil
}

func (s MessageService) GetMessageList(ctx context.Context, threadId uuid.UUID) ([]database.Message, error) {
	messageList, err := s.queries.GetMessages(ctx, threadId)

	if err != nil {
		return nil, err
	}

	return messageList, nil
}

func (s MessageService) CreateMessage(ctx context.Context, name string, threadId uuid.UUID, token string) (database.Message, error) {
	authorId, err := s.tokenService.GetIdFromToken(s.cfg.AccessSigningKey, token)

	if err != nil {
		return database.Message{}, err
	}

	message, err := s.queries.CreateMessage(ctx, database.CreateMessageParams{
		Name:     name,
		ThreadID: threadId,
		AuthorID: authorId,
	})

	if err != nil {
		return database.Message{}, err
	}

	return message, nil
}

func (s MessageService) DeleteMessage(ctx context.Context, messageId uuid.UUID) error {
	err := s.queries.DeleteMessage(ctx, messageId)

	if err != nil {
		return err
	}

	return nil
}

func (s MessageService) UpdateMessage(ctx context.Context, messageId uuid.UUID, newMessageParams *dto.MessageDto) (database.Message, error) {
	updatedMessage, err := s.queries.UpdateMessage(ctx, database.UpdateMessageParams{
		ID:       messageId,
		Name:     newMessageParams.Name,
		ThreadID: newMessageParams.ThreadID,
	})

	if err != nil {
		return database.Message{}, err
	}

	return updatedMessage, nil
}
