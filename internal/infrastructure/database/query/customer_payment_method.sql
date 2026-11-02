-- name: CreateCustomerPaymentMethod :execresult
INSERT INTO customer_payment_method (
    customer_id,
    account_name,
    payment_method_type_id,
    card_token,
    is_default,
    created_by
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: GetCustomerPaymentMethods :many
SELECT
    customer_payment_method_id,
    customer_id,
    account_name,
    payment_method_type_id,
    card_token,
    is_default,
    created_at,
    created_by,
    updated_at,
    updated_by
FROM customer_payment_method
WHERE customer_id = ?
ORDER BY is_default DESC, customer_payment_method_id DESC;

-- name: GetCustomerPaymentMethodByID :one
SELECT
    customer_payment_method_id,
    customer_id,
    account_name,
    payment_method_type_id,
    card_token,
    is_default,
    created_at,
    created_by,
    updated_at,
    updated_by
FROM customer_payment_method
WHERE customer_payment_method_id = ?
  AND customer_id = ?
LIMIT 1;

-- name: UpdateCustomerPaymentMethod :exec
UPDATE customer_payment_method
SET
    account_name = ?,
    payment_method_type_id = ?,
    card_token = ?,
    updated_by = ?
WHERE customer_payment_method_id = ?
  AND customer_id = ?;

-- name: ClearDefaultCustomerPaymentMethods :exec
UPDATE customer_payment_method
SET
    is_default = 0,
    updated_by = ?
WHERE customer_id = ?
  AND is_default = 1;

-- name: SetDefaultCustomerPaymentMethod :exec
UPDATE customer_payment_method
SET
    is_default = 1,
    updated_by = ?
WHERE customer_payment_method_id = ?
  AND customer_id = ?;

-- name: DeleteCustomerPaymentMethod :exec
DELETE FROM customer_payment_method
WHERE customer_payment_method_id = ?
  AND customer_id = ?;
