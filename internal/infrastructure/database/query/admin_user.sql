-- name: GetAdminUser :one
SELECT admin_user_id, username, role, email, is_active FROM admin_user WHERE admin_user_id = ? and is_active = 1;

-- name: GetAdminByUsername :one
SELECT admin_user_id, username, password, role FROM admin_user WHERE username = ? and is_active = 1;

-- name: ListAdminUsers :many
SELECT admin_user_id, username, role, email, is_active FROM admin_user WHERE is_active = 1 ORDER BY created_at;

-- name: CreateAdmin :execresult
INSERT INTO admin_user (username, role, email, password, created_at, created_by)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateAdmin :execresult
UPDATE admin_user SET
    username = ?,
    password = ?,
    email = ?,
    role = ?,
    updated_at = ?,
    updated_by = ?
WHERE admin_user_id = ? and is_active = 1;

-- name: UpdateAdminIsActive :execresult
UPDATE admin_user SET
    is_active = ?,
    updated_at = ?,
    updated_by = ?
WHERE admin_user_id = ?;
