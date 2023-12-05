// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"github.com/google/uuid"
)

type Permission struct {
	ID   int64
	Name string
}

type Role struct {
	ID   int64
	Name string
}

type RolesPermission struct {
	RoleID       int64
	PermissionID int64
}

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Login    string
	RoleID   int64
	Token    string
}
