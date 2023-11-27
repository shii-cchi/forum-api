package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Login    string    `json:"login"`
	Token    string    `json:"token"`
}

type UserForResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Login string    `json:"login"`
	Token string    `json:"token"`
}
