package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"github.com/shii-cchi/forum-api/internal/models"
	"github.com/shii-cchi/forum-api/pkg/auth"
	"github.com/shii-cchi/forum-api/pkg/hash"
	"time"
)

type UserService struct {
	queries *database.Queries
	hasher  *hash.SHA1Hasher
	cfg     *config.Config
}

func NewUserService(q *database.Queries, h *hash.SHA1Hasher, c *config.Config) *UserService {
	return &UserService{
		queries: q,
		hasher:  h,
		cfg:     c,
	}
}

func (s UserService) CreateUser(ctx context.Context, newUser *dto.UserDto) (models.UserForResponse, string, error) {
	count, err := s.queries.CheckUserIsExist(ctx, database.CheckUserIsExistParams{
		Email: newUser.Email,
		Login: newUser.Login,
	})

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	if count != 0 {
		return models.UserForResponse{}, "", errors.New("User with this login or email already exists")
	}

	passwordHash, err := s.hasher.Hash(newUser.Password)

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	user, err := s.queries.CreateUser(ctx, database.CreateUserParams{
		Email:    newUser.Email,
		Password: passwordHash,
		Login:    newUser.Login,
	})

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	accessToken, refreshToken, err := s.createTokens(user.ID.String())

	err = s.queries.AddToken(ctx, database.AddTokenParams{
		ID:    user.ID,
		Token: refreshToken,
	})

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	return models.UserForResponse{
		ID:    user.ID,
		Email: user.Email,
		Login: user.Login,
		Token: accessToken,
	}, refreshToken, nil
}

func (s UserService) Logout(ctx context.Context, accessToken string) error {
	userId, err := s.GetIdFromToken(s.cfg.AccessSigningKey, accessToken)

	if err != nil {
		return err
	}

	err = s.queries.LogoutUser(ctx, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s UserService) Login(ctx context.Context, checkedUser *dto.UserDto) (models.UserForResponse, string, error) {
	passwordHash, err := s.hasher.Hash(checkedUser.Password)

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	user, err := s.queries.CheckDataToLogin(ctx, database.CheckDataToLoginParams{
		Email:    checkedUser.Email,
		Password: passwordHash,
		Login:    checkedUser.Login,
	})

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	if user.ID == uuid.Nil {
		return models.UserForResponse{}, "", errors.New("Wrong credentials")
	}

	accessToken, refreshToken, err := s.createTokens(user.ID.String())

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	err = s.queries.AddToken(ctx, database.AddTokenParams{
		ID:    user.ID,
		Token: refreshToken,
	})

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	return models.UserForResponse{
		ID:    user.ID,
		Email: user.Email,
		Login: user.Login,
		Token: accessToken,
	}, refreshToken, nil
}

func (s UserService) Refresh(ctx context.Context, refreshToken string) (models.UserForResponse, string, error) {
	userId, err := s.GetIdFromToken(s.cfg.RefreshSigningKey, refreshToken)

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	accessToken, refreshToken, err := s.createTokens(userId.String())

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	err = s.queries.AddToken(ctx, database.AddTokenParams{
		ID:    userId,
		Token: refreshToken,
	})

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	user, err := s.queries.GetUser(ctx, userId)

	if err != nil {
		return models.UserForResponse{}, "", err
	}

	return models.UserForResponse{
		ID:    userId,
		Email: user.Email,
		Login: user.Login,
		Token: accessToken,
	}, refreshToken, nil
}

func (s UserService) createToken(signingKey string, ttl string, userId string) (string, error) {
	m, err := auth.NewManager(signingKey)

	if err != nil {
		return "", err
	}

	TTL, err := time.ParseDuration(ttl)

	token, err := m.NewJWT(userId, TTL)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s UserService) createTokens(userId string) (string, string, error) {
	accessToken, err := s.createToken(s.cfg.AccessSigningKey, s.cfg.AccessTTL, userId)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.createToken(s.cfg.RefreshSigningKey, s.cfg.RefreshTTL, userId)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s UserService) IsValidToken(validatedToken string) (bool, error) {
	m, err := auth.NewManager(s.cfg.RefreshSigningKey)
	if err != nil {
		return false, err
	}

	ok, err := m.IsValidToken(validatedToken)

	return ok, err
}

func (s UserService) GetIdFromToken(signingKey string, token string) (uuid.UUID, error) {
	m, err := auth.NewManager(signingKey)

	if err != nil {
		return uuid.Nil, err
	}

	userIdStr, err := m.Parse(token)

	if err != nil {
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}
