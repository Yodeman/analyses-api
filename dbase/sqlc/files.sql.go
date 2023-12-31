// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: files.sql

package db

import (
	"context"

	"github.com/lib/pq"
)

const createFile = `-- name: CreateFile :one
INSERT INTO files (
    username,
    data
) VALUES (
    $1, $2
)
RETURNING id, username, data, created_at
`

type CreateFileParams struct {
	Username string        `json:"username"`
	Data     []interface{} `json:"data"`
}

func (q *Queries) CreateFile(ctx context.Context, arg CreateFileParams) (File, error) {
	row := q.db.QueryRowContext(ctx, createFile, arg.Username, pq.Array(arg.Data))
	var i File
	err := row.Scan(
		&i.ID,
		&i.Username,
		pq.Array(&i.Data),
		&i.CreatedAt,
	)
	return i, err
}

const getFile = `-- name: GetFile :one
SELECT id, username, data, created_at FROM files
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetFile(ctx context.Context, username string) (File, error) {
	row := q.db.QueryRowContext(ctx, getFile, username)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Username,
		pq.Array(&i.Data),
		&i.CreatedAt,
	)
	return i, err
}

const updateFile = `-- name: UpdateFile :one
UPDATE files
SET data = $1
WHERE username = $2
RETURNING id, username, data, created_at
`

type UpdateFileParams struct {
	Data     []interface{} `json:"data"`
	Username string        `json:"username"`
}

func (q *Queries) UpdateFile(ctx context.Context, arg UpdateFileParams) (File, error) {
	row := q.db.QueryRowContext(ctx, updateFile, pq.Array(arg.Data), arg.Username)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Username,
		pq.Array(&i.Data),
		&i.CreatedAt,
	)
	return i, err
}
