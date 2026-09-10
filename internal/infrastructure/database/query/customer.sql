-- name: GetCustomer :one
SELECT customer_id, uuid, first_name, last_name, username, email, birth_date, phone_number, is_active, created_at FROM customer WHERE customer_id = ? and is_active = 1;

-- name: GetCustomerByUsername :one
SELECT customer_id, uuid, username, password FROM customer WHERE username = ? and is_active = 1;

-- name: ListCustomers :many
SELECT customer_id, uuid, first_name, last_name, username, email, birth_date, phone_number, is_active, created_at FROM customer WHERE is_active = 1 ORDER BY created_at;

-- name: CreateCustomer :execresult
INSERT INTO customer (
    uuid,
    first_name,
    last_name,
    username,
    email,
    password,
    birth_date,
    phone_number,
    created_at,
    created_by
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateMyProfile :execresult
UPDATE customer SET
    first_name = ?,
    last_name = ?,
    username = ?,
    email = ?,
    birth_date = ?,
    phone_number = ?,
    updated_at = ?,
    updated_by = ?
WHERE customer_id = ? and is_active = 1;

-- name: UpdateCustomerIsActive :execresult
UPDATE customer SET
    is_active = ?,
    updated_at = ?,
    updated_by = ?
WHERE customer_id = ?;

-- name: UpdateCustomerPassword :execresult
UPDATE customer SET
    password = ?,
    updated_at = ?,
    updated_by = ?
WHERE customer_id = ? and is_active = 1;

-- name: AdminGetCustomerDetail :one
SELECT customer_id, uuid, first_name, last_name, username, email, birth_date, phone_number, is_active, created_at FROM customer WHERE uuid = ? and is_active = 1;

-- name: AdminListCustomers :many
-- Pagination (Limit/Offset) 
SELECT customer_id, uuid, first_name, last_name, username, is_active, created_at 
FROM customer 
ORDER BY created_at DESC 
LIMIT ? OFFSET ?;

-- name: AdminUpdateCustomerStatus :execresult
UPDATE customer SET
    is_active = ?,
    updated_at = ?,
    updated_by = ?
WHERE uuid = ?;
