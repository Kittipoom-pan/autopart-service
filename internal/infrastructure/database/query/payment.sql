-- name: CreatePayment :execresult
INSERT INTO payment (
    order_id,
    amount,
    customer_payment_method_id,
    payment_method_type_id,
    status,
    created_by
)
VALUES (?, ?, ?, ?, ?, ?);
