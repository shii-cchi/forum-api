// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: tokens.sql

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
