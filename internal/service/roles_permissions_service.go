package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers/dto"
)

type RolesPermissionsService struct {
	queries *database.Queries
}

func NewRolesPermissionsService(q *database.Queries) *RolesPermissionsService {
	return &RolesPermissionsService{
		queries: q,
	}
}

func (r RolesPermissionsService) GetRoleAndPermissions(ctx context.Context, userId uuid.UUID) (dto.RoleDto, error) {
	role, err := r.queries.GetRole(ctx, userId)

	if err != nil {
		return dto.RoleDto{}, err
	}

	permissions, err := r.queries.GetPermissions(ctx, role)

	if err != nil {
		return dto.RoleDto{}, err
	}

	return dto.RoleDto{
		Name:        role,
		Permissions: permissions,
	}, nil
}
