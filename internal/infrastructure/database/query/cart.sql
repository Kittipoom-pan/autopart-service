-- name: GetCartByCustomerID :one
SELECT
    cart_id,
    customer_id,
    subtotal,
    is_checkout,
    created_at,
    updated_at
FROM cart
WHERE customer_id = ? AND is_checkout = 0
ORDER BY cart_id DESC
LIMIT 1;

-- name: GetOrCreateActiveCart :execresult
INSERT INTO cart (customer_id)
VALUES (?)
ON DUPLICATE KEY UPDATE cart_id = LAST_INSERT_ID(cart_id);

-- name: RecalcCartSubtotal :execresult
UPDATE cart
SET subtotal = (
    SELECT COALESCE(SUM(ci.quantity * p.price), 0)
    FROM cart_item ci
    INNER JOIN part p ON p.part_id = ci.part_id
    WHERE ci.cart_id = cart.cart_id
)
WHERE cart.cart_id = ?;

-- name: ListCartItemsWithPart :many
SELECT
    ci.cart_item_id,
    ci.cart_id,
    ci.part_id,
    ci.quantity,
    p.name AS part_name,
    p.sku,
    p.price,
    ci.created_at,
    ci.updated_at
FROM cart_item ci
INNER JOIN part p ON p.part_id = ci.part_id
WHERE ci.cart_id = ?
ORDER BY ci.created_at;

-- name: GetCartItemByCartAndPart :one
SELECT
    cart_item_id,
    cart_id,
    part_id,
    quantity,
    created_at,
    updated_at
FROM cart_item
WHERE cart_id = ? AND part_id = ?;

-- name: CreateCartItem :execresult
INSERT INTO cart_item (cart_id, part_id, quantity)
VALUES (?, ?, ?);

-- name: UpdateCartItemQuantity :execresult
UPDATE cart_item
SET quantity = ?, updated_at = NOW()
WHERE cart_item_id = ? AND cart_id = ?;

-- name: GetCartItemByID :one
SELECT
    cart_item_id,
    cart_id,
    part_id,
    quantity,
    created_at,
    updated_at
FROM cart_item
WHERE cart_item_id = ?;

-- name: UpdateCartCheckoutByID :execresult
UPDATE cart
SET is_checkout = 1, updated_at = NOW()
WHERE cart_id = ? AND is_checkout = 0;
