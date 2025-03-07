-- name: CreateOrder :execresult
INSERT INTO `order` (customer_id, customer_payment_method_id, cart_id, amount, address_id, discount_id, status, created_at, created_by, updated_at, updated_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetOrderByID :one
SELECT * FROM `order` WHERE order_id = ?;

-- name: ListOrders :many
SELECT * FROM `order` ORDER BY created_at DESC;