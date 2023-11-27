package dto

type UserDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Login    string `json:"login"`
}
