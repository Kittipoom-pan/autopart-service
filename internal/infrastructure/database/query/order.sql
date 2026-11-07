-- name: CreateOrder :execresult
INSERT INTO `order` (
    order_number,
    customer_id,
    channel,
    placed_by_admin_id,
    cart_id,
    subtotal,
    discount_amount,
    shipping_fee,
    total_amount,
    discount_id,
    status,
    remark
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetOrderByID :one
SELECT
    order_id,
    order_number,
    customer_id,
    channel,
    placed_by_admin_id,
    cart_id,
    subtotal,
    discount_amount,
    shipping_fee,
    total_amount,
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
    placed_by_admin_id,
    cart_id,
    subtotal,
    discount_amount,
    shipping_fee,
    total_amount,
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
    placed_by_admin_id,
    cart_id,
    subtotal,
    discount_amount,
    shipping_fee,
    total_amount,
    discount_id,
    status,
    remark,
    created_at,
    updated_at
FROM `order`
WHERE order_number = ?;

-- name: NextOrderSequence :execresult
INSERT INTO order_sequence (seq_date, `last_value`)
VALUES (?, LAST_INSERT_ID(1))
ON DUPLICATE KEY UPDATE `last_value` = LAST_INSERT_ID(`last_value` + 1);

-- name: CreateOrderItem :execresult
INSERT INTO order_item (order_id, part_id, quantity, price)
VALUES (?, ?, ?, ?);

-- name: CreateOrderShippingAddress :execresult
INSERT INTO order_shipping_address (order_id, recipient_name, phone_number, address_line, subdistrict, district, province, zip_code)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);