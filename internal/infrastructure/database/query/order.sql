-- name: CreateOrder :execresult
INSERT INTO `order` (
    order_number,
    customer_id,
    channel,
    placed_by_admin_user_id,
    customer_payment_method_id,
    cart_id,
    amount,
    address_id,
    discount_id,
    status,
    remark,
    created_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetOrderByID :one
SELECT
    order_id,
    order_number,
    customer_id,
    channel,
    placed_by_admin_user_id,
    customer_payment_method_id,
    cart_id,
    amount,
    address_id,
    discount_id,
    status,
    remark,
    created_at,
    updated_at
FROM `order`
WHERE order_id = ?;

-- name: ListOrders :many
SELECT
    order_id,
    order_number,
    customer_id,
    channel,
    placed_by_admin_user_id,
    customer_payment_method_id,
    cart_id,
    amount,
    address_id,
    discount_id,
    status,
    remark,
    created_at,
    updated_at
FROM `order`
ORDER BY created_at DESC;

-- name: GetOrderByOrderNumber :one
SELECT 
    order_id, 
    order_number, 
    customer_id, 
    channel, 
    placed_by_admin_user_id, 
    customer_payment_method_id, 
    cart_id, 
    amount, 
    address_id, 
    discount_id, 
    status, 
    remark, 
    created_at, 
    updated_at
FROM `order`
WHERE order_number = ?;