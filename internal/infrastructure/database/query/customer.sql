-- name: GetCustomer :one
SELECT * FROM customer WHERE customer_id = ?;

-- name: ListCustomers :many
SELECT * FROM customer ORDER BY last_name;

-- name: CreateCustomer :execresult
INSERT INTO customer (first_name, last_name, username, email, password, birth_date, phone_number, created_at, created_by, updated_at, updated_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);