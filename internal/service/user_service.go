package service

import (
	"context"
	"errors"
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"github.com/shii-cchi/forum-api/pkg/hash"
)

type UserService struct {
	queries *database.Queries
	hasher  *hash.SHA1Hasher
	cfg     *config.Config

	tokenService            Tokens
	rolesPermissionsService RolesPermissions
}

func NewUserService(q *database.Queries, h *hash.SHA1Hasher, c *config.Config, t *TokenService, r *RolesPermissionsService) *UserService {
	return &UserService{
		queries:                 q,
		hasher:                  h,
		cfg:                     c,
		tokenService:            t,
		rolesPermissionsService: r,
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

	accessToken, refreshToken, err := s.tokenService.CreateTokens(user.ID.String())

	err = s.queries.AddToken(ctx, database.AddTokenParams{
		ID:    user.ID,
		Token: refreshToken,
	})

	if err != nil {
		return dto.UserPreviewDto{}, "", err
	}

	role, err := s.rolesPermissionsService.GetRoleAndPermissions(ctx, user.ID)

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

func (s UserService) Logout(ctx context.Context, accessToken string) (bool, error) {
	ok, err := s.tokenService.IsValidToken(accessToken, s.cfg.AccessSigningKey)

	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	userId, err := s.tokenService.GetIdFromToken(s.cfg.AccessSigningKey, accessToken)

	if err != nil {
		return false, err
	}

	err = s.queries.LogoutUser(ctx, userId)

	if err != nil {
		return false, err
	}

	return true, nil
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

	accessToken, refreshToken, err := s.tokenService.CreateTokens(user.ID.String())

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

	role, err := s.rolesPermissionsService.GetRoleAndPermissions(ctx, user.ID)

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

func IsEmptyUser(user database.User) bool {
	emptyUser := database.User{}

	return emptyUser == user
}
