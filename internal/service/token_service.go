package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"github.com/shii-cchi/forum-api/pkg/auth"
	"time"
)

type TokenService struct {
	queries                 *database.Queries
	cfg                     *config.Config
	rolesPermissionsService *RolesPermissionsService
}

func NewTokenService(q *database.Queries, c *config.Config, r *RolesPermissionsService) *TokenService {
	return &TokenService{
		queries:                 q,
		cfg:                     c,
		rolesPermissionsService: r,
	}
}

func (t TokenService) Refresh(ctx context.Context, refreshToken string) (dto.UserPreviewDto, string, bool, error) {
	ok, err := t.IsValidToken(refreshToken, t.cfg.RefreshSigningKey)

	if err != nil {
		return dto.UserPreviewDto{}, "", false, err
	}

	if !ok {
		return dto.UserPreviewDto{}, "", false, nil
	}

	userId, err := t.GetIdFromToken(t.cfg.RefreshSigningKey, refreshToken)

	if err != nil {
		return dto.UserPreviewDto{}, "", false, err
	}

	accessToken, refreshToken, err := t.CreateTokens(userId.String())

	if err != nil {
		return dto.UserPreviewDto{}, "", false, err
	}

	err = t.queries.AddToken(ctx, database.AddTokenParams{
		ID:    userId,
		Token: refreshToken,
	})

	if err != nil {
		return dto.UserPreviewDto{}, "", false, err
	}

	user, err := t.queries.FindUserById(ctx, userId)

	if err != nil {
		return dto.UserPreviewDto{}, "", false, err
	}

	role, err := t.rolesPermissionsService.GetRoleAndPermissions(ctx, user.ID)

	if err != nil {
		return dto.UserPreviewDto{}, "", false, err
	}

	return dto.UserPreviewDto{
		ID:    userId,
		Email: user.Email,
		Login: user.Login,
		Role:  role,
		Token: accessToken,
	}, refreshToken, true, nil
}

func (t TokenService) CreateToken(signingKey string, ttl string, userId string) (string, error) {
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

func (t TokenService) CreateTokens(userId string) (string, string, error) {
	accessToken, err := t.CreateToken(t.cfg.AccessSigningKey, t.cfg.AccessTTL, userId)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.CreateToken(t.cfg.RefreshSigningKey, t.cfg.RefreshTTL, userId)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (t TokenService) IsValidToken(validatedToken string, signingKey string) (bool, error) {
	m, err := auth.NewManager(signingKey)
	if err != nil {
		return false, err
	}

	ok, err := m.IsValidToken(validatedToken)

	return ok, err
}

func (t TokenService) GetIdFromToken(signingKey string, token string) (uuid.UUID, error) {
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
