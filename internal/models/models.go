package models

import "github.com/google/uuid"

type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Login    string    `json:"login"`
	Role     Role      `json:"role"`
	Token    string    `json:"token"`
}

type UserForResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Login string    `json:"login"`
	Role  Role      `json:"role"`
	Token string    `json:"token"`
}
