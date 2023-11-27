package service

import (
	"context"
	"errors"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"github.com/shii-cchi/forum-api/internal/models"
	"github.com/shii-cchi/forum-api/pkg/hash"
)

type UserService struct {
	queries *database.Queries
	hasher  *hash.SHA1Hasher
}

func NewUserService(q *database.Queries, h *hash.SHA1Hasher) *UserService {
	return &UserService{
		queries: q,
		hasher:  h,
	}
}

func (s UserService) CreateUser(ctx context.Context, newUser *dto.UserDto) (interface{}, string, error) {
	count, err := s.queries.CheckUserIsExist(ctx, database.CheckUserIsExistParams{
		Email: newUser.Email,
		Login: newUser.Login,
	})

	if err != nil {
		return nil, "", err
	}

	if count != 0 {
		return nil, "", errors.New("User with this login or email already exists")
	}

	passwordHash, err := s.hasher.Hash(newUser.Password)

	if err != nil {
		return nil, "", err
	}

	user, err := s.queries.CreateUser(ctx, database.CreateUserParams{
		Email:    newUser.Email,
		Password: passwordHash,
		Login:    newUser.Login,
	})

	if err != nil {
		return nil, "", err
	}

	refreshToken, err := CreateToken("REFRESH_SIGNING_KEY", "REFRESH_TTL", user.ID.String())

	if err != nil {
		return nil, "", err
	}

	accessToken, err := CreateToken("ACCESS_SIGNING_KEY", "ACCESS_TTL", user.ID.String())

	if err != nil {
		return nil, "", err
	}

	err = s.queries.AddToken(ctx, database.AddTokenParams{
		ID:    user.ID,
		Token: refreshToken,
	})

	if err != nil {
		return nil, "", err
	}

	return models.UserForResponse{
		ID:    user.ID,
		Email: user.Email,
		Login: user.Login,
		Token: accessToken,
	}, refreshToken, nil
}
