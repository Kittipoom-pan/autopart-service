-- name: GetCustomer :one
SELECT customer_id, first_name, last_name, username, email, birth_date, phone_number, is_active FROM customer WHERE customer_id = ? and is_active = 1;

-- name: ListCustomers :many
SELECT customer_id, first_name, last_name, username, email, birth_date, phone_number, is_active FROM customer WHERE is_active = 1 ORDER BY created_at;

-- name: CreateCustomer :execresult
INSERT INTO customer (first_name, last_name, username, email, password, birth_date, phone_number, created_at, created_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateCustomer :execresult
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
