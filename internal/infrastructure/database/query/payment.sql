-- name: CreatePayment :execresult
INSERT INTO payment (
    order_id,
    amount,
    customer_payment_method_id,
    payment_method_type_id,
    status,
    paid_at,
    created_by_admin_id
)
VALUES (?, ?, ?, ?, ?, ?, ?);
