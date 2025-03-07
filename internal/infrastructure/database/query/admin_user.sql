-- name: GetAdminUser :one
SELECT * FROM admin_user WHERE admin_user_id = ?;

-- name: ListAdminUsers :many
SELECT * FROM admin_user ORDER BY username;