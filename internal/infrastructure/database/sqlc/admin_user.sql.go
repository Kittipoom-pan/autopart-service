// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: admin_user.sql

package db

import (
	"context"
)

const getAdminUser = `-- name: GetAdminUser :one
SELECT admin_user_id, username, email, password, role, created_at, created_by, updated_at, updated_by FROM admin_user WHERE admin_user_id = ?
`

func (q *Queries) GetAdminUser(ctx context.Context, adminUserID int32) (AdminUser, error) {
	row := q.db.QueryRowContext(ctx, getAdminUser, adminUserID)
	var i AdminUser
	err := row.Scan(
		&i.AdminUserID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const listAdminUsers = `-- name: ListAdminUsers :many
SELECT admin_user_id, username, email, password, role, created_at, created_by, updated_at, updated_by FROM admin_user ORDER BY username
`

func (q *Queries) ListAdminUsers(ctx context.Context) ([]AdminUser, error) {
	rows, err := q.db.QueryContext(ctx, listAdminUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AdminUser
	for rows.Next() {
		var i AdminUser
		if err := rows.Scan(
			&i.AdminUserID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Role,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
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
