// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package database

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (name, bio)
VALUES ($1, $2)
RETURNING id, name, bio
`

type CreateAuthorParams struct {
	Name string
	Bio  string
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	ctx, span :=  otel.Tracer("CreateAuthorQuery").Start(ctx, "CreateAuthorQueryContext")
	defer span.End()

	row := q.db.QueryRowContext(ctx, createAuthor, arg.Name, arg.Bio)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "error")
			span.RecordError(err)
		}
	}()
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE
FROM authors
WHERE id = $1
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int64) error {
	ctx, span :=  otel.Tracer("DeleteAuthorQuery").Start(ctx, "DeleteAuthorQueryContext")
	defer span.End()
	_, err := q.db.ExecContext(ctx, deleteAuthor, id)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "error")
			span.RecordError(err)
		}
	}()
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio
FROM authors
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error) {
	ctx, span :=  otel.Tracer("GetAuthorQuery").Start(ctx, "GetAuthorQueryContext")
	defer span.End()
	row := q.db.QueryRowContext(ctx, getAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "error")
			span.RecordError(err)
		}
	}()
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, bio
FROM authors
ORDER BY name
`

func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
	ctx, span :=  otel.Tracer("ListAuthorQuery").Start(ctx, "ListAuthorQueryContext")
	defer span.End()
	rows, err := q.db.QueryContext(ctx, listAuthors)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "error")
			span.RecordError(err)
		}
	}()
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name, &i.Bio); err != nil {
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

const partialUpdateAuthor = `-- name: PartialUpdateAuthor :one
UPDATE authors
SET name = CASE WHEN $1::boolean THEN $2::VARCHAR(32) ELSE name END,
    bio  = CASE WHEN $3::boolean THEN $4::TEXT ELSE bio END
WHERE id = $5
RETURNING id, name, bio
`

type PartialUpdateAuthorParams struct {
	UpdateName bool
	Name       string
	UpdateBio  bool
	Bio        string
	ID         int64
}

func (q *Queries) PartialUpdateAuthor(ctx context.Context, arg PartialUpdateAuthorParams) (Author, error) {
	ctx, span :=  otel.Tracer("PartialUpdateAuthorQuery").Start(ctx, "PartialUpdateAuthorQueryContext")
	defer span.End()
	row := q.db.QueryRowContext(ctx, partialUpdateAuthor,
		arg.UpdateName,
		arg.Name,
		arg.UpdateBio,
		arg.Bio,
		arg.ID,
	)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "error")
			span.RecordError(err)
		}
	}()
	return i, err
}

const truncateAuthor = `-- name: TruncateAuthor :exec
TRUNCATE authors
`

func (q *Queries) TruncateAuthor(ctx context.Context) error {
	ctx, span :=  otel.Tracer("TruncateAuthorQuery").Start(ctx, "TruncateAuthorQueryContext")
	defer span.End()
	_, err := q.db.ExecContext(ctx, truncateAuthor)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "error")
			span.RecordError(err)
		}
	}()
	return err
}

const updateAuthor = `-- name: UpdateAuthor :one
UPDATE authors
SET name = $2,
    bio  = $3
WHERE id = $1
RETURNING id, name, bio
`

type UpdateAuthorParams struct {
	ID   int64
	Name string
	Bio  string
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error) {
	ctx, span :=  otel.Tracer("UpdateAuthorQuery").Start(ctx, "UpdateAuthorQueryContext")
	defer span.End()
	row := q.db.QueryRowContext(ctx, updateAuthor, arg.ID, arg.Name, arg.Bio)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "error")
			span.RecordError(err)
		}
	}()
	return i, err
}
