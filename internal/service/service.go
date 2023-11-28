package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"github.com/shii-cchi/forum-api/internal/models"
	"github.com/shii-cchi/forum-api/pkg/auth"
	"github.com/shii-cchi/forum-api/pkg/hash"
	"os"
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

func (s UserService) Logout(ctx context.Context, accessToken string) error {
	signingKey := os.Getenv("ACCESS_SIGNING_KEY")
	if signingKey == "" {
		return errors.New("ACCESS_SIGNING_KEY is not found in the environment")
	}

	m, err := auth.NewManager(signingKey)

	if err != nil {
		return err
	}

	userIdStr, err := m.Parse(accessToken)

	if err != nil {
		return err
	}

	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		return err
	}

	err = s.queries.LogoutUser(ctx, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s UserService) Login(ctx context.Context, checkedUser *dto.UserDto) (interface{}, string, error) {
	user, err := s.queries.CheckDataToLogin(ctx, database.CheckDataToLoginParams{
		Email:    checkedUser.Email,
		Password: checkedUser.Password,
		Login:    checkedUser.Login,
	})

	if err != nil {
		return nil, "", err
	}

	if user.ID == uuid.Nil {
		return nil, "", errors.New("Wrong credentials")
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
