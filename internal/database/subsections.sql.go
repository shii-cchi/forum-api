// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: subsections.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createSubsection = `-- name: CreateSubsection :one
INSERT INTO subsections (name, section_id, author_id)
VALUES ($1, $2, $3)
RETURNING id, name, section_id, author_id
`

type CreateSubsectionParams struct {
	Name      string
	SectionID uuid.UUID
	AuthorID  uuid.UUID
}

func (q *Queries) CreateSubsection(ctx context.Context, arg CreateSubsectionParams) (Subsection, error) {
	row := q.db.QueryRowContext(ctx, createSubsection, arg.Name, arg.SectionID, arg.AuthorID)
	var i Subsection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SectionID,
		&i.AuthorID,
	)
	return i, err
}

const deleteSubsection = `-- name: DeleteSubsection :exec
DELETE FROM subsections
WHERE id = $1
`

func (q *Queries) DeleteSubsection(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteSubsection, id)
	return err
}

const getSubsection = `-- name: GetSubsection :one
SELECT id, name, section_id, author_id
FROM subsections
WHERE id = $1
`

func (q *Queries) GetSubsection(ctx context.Context, id uuid.UUID) (Subsection, error) {
	row := q.db.QueryRowContext(ctx, getSubsection, id)
	var i Subsection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SectionID,
		&i.AuthorID,
	)
	return i, err
}

const getSubsections = `-- name: GetSubsections :many
SELECT id, name, section_id, author_id
FROM subsections
WHERE section_id = $1
`

func (q *Queries) GetSubsections(ctx context.Context, sectionID uuid.UUID) ([]Subsection, error) {
	rows, err := q.db.QueryContext(ctx, getSubsections, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subsection
	for rows.Next() {
		var i Subsection
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.SectionID,
			&i.AuthorID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSubsection = `-- name: UpdateSubsection :one
UPDATE subsections
SET name = $2, section_id = $3
WHERE id = $1
RETURNING id, name, section_id, author_id
`

type UpdateSubsectionParams struct {
	ID        uuid.UUID
	Name      string
	SectionID uuid.UUID
}

func (q *Queries) UpdateSubsection(ctx context.Context, arg UpdateSubsectionParams) (Subsection, error) {
	row := q.db.QueryRowContext(ctx, updateSubsection, arg.ID, arg.Name, arg.SectionID)
	var i Subsection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SectionID,
		&i.AuthorID,
	)
	return i, err
}
