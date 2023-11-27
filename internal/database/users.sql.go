// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addToken = `-- name: AddToken :exec
UPDATE users
SET token = $2
WHERE id = $1
`

type AddTokenParams struct {
	ID    uuid.UUID
	Token string
}

func (q *Queries) AddToken(ctx context.Context, arg AddTokenParams) error {
	_, err := q.db.ExecContext(ctx, addToken, arg.ID, arg.Token)
	return err
}

const checkUserIsExist = `-- name: CheckUserIsExist :one
SELECT COUNT(*)
FROM users
WHERE email = $1 OR login = $2
`

type CheckUserIsExistParams struct {
	Email string
	Login string
}

func (q *Queries) CheckUserIsExist(ctx context.Context, arg CheckUserIsExistParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkUserIsExist, arg.Email, arg.Login)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password, login)
VALUES ($1, $2, $3)
RETURNING id, email, password, login, token
`

type CreateUserParams struct {
	Email    string
	Password string
	Login    string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.Password, arg.Login)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Login,
		&i.Token,
	)
	return i, err
}
