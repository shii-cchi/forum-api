package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
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

func (s UserService) CreateUser(ctx context.Context, newUser *dto.UserDto) (dto.UserPreviewDto, string, error) {
	count, err := s.queries.CheckUserIsExist(ctx, database.CheckUserIsExistParams{
		Email: newUser.Email,
		Login: newUser.Login,
	})

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	if count != 0 {
		return dto.UserPreviewDto{}, "", errors.New("User with this login or email already exists")
	}

	passwordHash, err := s.hasher.Hash(newUser.Password)

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	user, err := s.queries.CreateUser(ctx, database.CreateUserParams{
		Email:    newUser.Email,
		Password: passwordHash,
		Login:    newUser.Login,
	})

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	accessToken, refreshToken, err := s.createTokens(user.ID.String())

	err = s.queries.AddToken(ctx, database.AddTokenParams{
		ID:    user.ID,
		Token: refreshToken,
	})

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	role, err := s.GetRoleAndPermissions(ctx, user.ID)

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	return dto.UserPreviewDto{
		ID:    user.ID,
		Email: user.Email,
		Login: user.Login,
		Role:  role,
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

func (s UserService) Login(ctx context.Context, checkedUser *dto.UserDto) (dto.UserPreviewDto, string, error) {
	user, err := s.queries.CheckDataToLogin(ctx, database.CheckDataToLoginParams{
		Email: checkedUser.Email,
		Login: checkedUser.Login,
	})

	if IsEmptyUser(user) {
		return dto.UserPreviewDto{}, "", errors.New("Wrong credentials")
	}

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	passwordHash, err := s.hasher.Hash(checkedUser.Password)

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	if passwordHash != user.Password {
		return dto.UserPreviewDto{}, "", errors.New("Wrong credentials")
	}

	accessToken, refreshToken, err := s.createTokens(user.ID.String())

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	err = s.queries.AddToken(ctx, database.AddTokenParams{
		ID:    user.ID,
		Token: refreshToken,
	})

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	role, err := s.GetRoleAndPermissions(ctx, user.ID)

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	return dto.UserPreviewDto{
		ID:    user.ID,
		Email: user.Email,
		Login: user.Login,
		Role:  role,
		Token: accessToken,
	}, refreshToken, nil
}

func (s UserService) Refresh(ctx context.Context, refreshToken string) (dto.UserPreviewDto, string, error) {
	userId, err := s.GetIdFromToken(s.cfg.RefreshSigningKey, refreshToken)

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	accessToken, refreshToken, err := s.createTokens(userId.String())

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	err = s.queries.AddToken(ctx, database.AddTokenParams{
		ID:    userId,
		Token: refreshToken,
	})

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	user, err := s.queries.FindUserById(ctx, userId)

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	role, err := s.GetRoleAndPermissions(ctx, user.ID)

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	return dto.UserPreviewDto{
		ID:    userId,
		Email: user.Email,
		Login: user.Login,
		Role:  role,
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

func (s UserService) IsValidToken(validatedToken string, signingKey string) (bool, error) {
	m, err := auth.NewManager(signingKey)
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

func (s UserService) GetRoleAndPermissions(ctx context.Context, userId uuid.UUID) (dto.RoleDto, error) {
	role, err := s.queries.GetRole(ctx, userId)

	if err != nil {
		return dto.RoleDto{}, err
	}

	permissions, err := s.queries.GetPermissions(ctx, role)

	if err != nil {
		return dto.RoleDto{}, err
	}

	return dto.RoleDto{
		Name:        role,
		Permissions: permissions,
	}, nil
}

func IsEmptyUser(user database.User) bool {
	emptyUser := database.User{}

	return emptyUser == user
}
