package dto

import "github.com/google/uuid"

type UserDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Login    string `json:"login"`
}

type UserPreviewDto struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Login string    `json:"login"`
	Role  RoleDto   `json:"role"`
	Token string    `json:"token"`
}
