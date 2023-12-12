package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
	"github.com/shii-cchi/forum-api/pkg/hash"
)

type Users interface {
	CreateUser(ctx context.Context, newUser *dto.UserDto) (dto.UserPreviewDto, string, error)
	Logout(ctx context.Context, accessToken string) (bool, error)
	Login(ctx context.Context, checkedUser *dto.UserDto) (dto.UserPreviewDto, string, error)
}

type Tokens interface {
	Refresh(ctx context.Context, refreshToken string) (dto.UserPreviewDto, string, bool, error)
	CreateToken(signingKey string, ttl string, userId string) (string, error)
	CreateTokens(userId string) (string, string, error)
	IsValidToken(validatedToken string, signingKey string) (bool, error)
	GetIdFromToken(signingKey string, token string) (uuid.UUID, error)
}

type RolesPermissions interface {
	GetRoleAndPermissions(ctx context.Context, userId uuid.UUID) (dto.RoleDto, error)
}

type Services struct {
	Users            Users
	Tokens           Tokens
	RolesPermissions RolesPermissions
}

type Deps struct {
	Queries *database.Queries
	Hasher  *hash.SHA1Hasher
	Config  *config.Config
}

func NewServices(deps Deps) *Services {
	rolesPermissionsService := NewRolesPermissionsService(deps.Queries)
	tokenService := NewTokenService(deps.Queries, deps.Config, rolesPermissionsService)
	userService := NewUserService(deps.Queries, deps.Hasher, deps.Config, tokenService, rolesPermissionsService)

	return &Services{
		Users:            userService,
		Tokens:           tokenService,
		RolesPermissions: rolesPermissionsService,
	}
}
