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

type Sections interface {
	GetSection(ctx context.Context, sectionId uuid.UUID) (database.Section, error)
	GetSectionList(ctx context.Context) ([]database.Section, error)
	CreateSection(ctx context.Context, name string) (database.Section, error)
	DeleteSection(ctx context.Context, sectionId uuid.UUID) error
	UpdateSectionName(ctx context.Context, sectionId uuid.UUID, newName string) (database.Section, error)
}

type Subsections interface {
	GetSubsection(ctx context.Context, subsectionId uuid.UUID) (database.Subsection, error)
	GetSubsectionList(ctx context.Context, sectionId uuid.UUID) ([]database.Subsection, error)
	CreateSubsection(ctx context.Context, name string, sectionId uuid.UUID, token string) (database.Subsection, error)
	DeleteSubsection(ctx context.Context, subsectionId uuid.UUID) error
	UpdateSubsection(ctx context.Context, subsectionId uuid.UUID, newSubsectionParams *dto.SubsectionDto) (database.Subsection, error)
}

type Themes interface {
	GetTheme(ctx context.Context, themeId uuid.UUID) (database.Theme, error)
	GetThemeList(ctx context.Context, subsectionId uuid.UUID) ([]database.Theme, error)
	CreateTheme(ctx context.Context, name string, subsectionId uuid.UUID, token string) (database.Theme, error)
	DeleteTheme(ctx context.Context, themeId uuid.UUID) error
	UpdateTheme(ctx context.Context, themeId uuid.UUID, newThemeParams *dto.ThemeDto) (database.Theme, error)
}

type Threads interface {
	GetThread(ctx context.Context, threadId uuid.UUID) (database.Thread, error)
	GetThreadList(ctx context.Context, themeId uuid.UUID) ([]database.Thread, error)
	CreateThread(ctx context.Context, name string, themeId uuid.UUID, token string) (database.Thread, error)
	DeleteThread(ctx context.Context, threadId uuid.UUID) error
	UpdateThread(ctx context.Context, threadId uuid.UUID, newThreadParams *dto.ThreadDto) (database.Thread, error)
}

type Messages interface {
	GetMessage(ctx context.Context, messageId uuid.UUID) (database.Message, error)
	GetMessageList(ctx context.Context, threadId uuid.UUID) ([]database.Message, error)
	CreateMessage(ctx context.Context, name string, threadId uuid.UUID, token string) (database.Message, error)
	DeleteMessage(ctx context.Context, messageId uuid.UUID) error
	UpdateMessage(ctx context.Context, messageId uuid.UUID, newMessageParams *dto.MessageDto) (database.Message, error)
}

type Services struct {
	Users            Users
	Tokens           Tokens
	RolesPermissions RolesPermissions
	Sections         Sections
	Subsections      Subsections
	Themes           Themes
	Threads          Threads
	Messages         Messages
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
	sectionService := NewSectionService(deps.Queries)
	subsectionService := NewSubsectionService(deps.Queries, deps.Config, tokenService)
	themesService := NewThemeService(deps.Queries, deps.Config, tokenService)
	threadsService := NewThreadService(deps.Queries, deps.Config, tokenService)
	messagesService := NewMessageService(deps.Queries, deps.Config, tokenService)

	return &Services{
		Users:            userService,
		Tokens:           tokenService,
		RolesPermissions: rolesPermissionsService,
		Sections:         sectionService,
		Subsections:      subsectionService,
		Themes:           themesService,
		Threads:          threadsService,
		Messages:         messagesService,
	}
}
