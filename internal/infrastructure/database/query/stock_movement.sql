-- name: DeductPartStock :execresult
UPDATE part
SET
    quantity = quantity - sqlc.arg(quantity),
    updated_at = NOW(),
    updated_by = sqlc.arg(updated_by)
WHERE part_id = sqlc.arg(part_id)
  AND is_active = 1
  AND quantity >= sqlc.arg(quantity);

-- name: CreateStockMovement :execresult
INSERT INTO stock_movement (
    part_id,
    quantity_change,
    price,
    event_type,
    reason,
    order_id,
    remark,
    created_by
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);
