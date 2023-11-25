package service

import "github.com/shii-cchi/forum-api/internal/database"

type UserService struct {
	queries *database.Queries
}

func NewUserService(q *database.Queries) *UserService {
	return &UserService{
		queries: q,
	}
}
