-- name: CreateUser :one
INSERT INTO users (email, password, login, token)
VALUES ($1, $2, $3, $4)
RETURNING *;