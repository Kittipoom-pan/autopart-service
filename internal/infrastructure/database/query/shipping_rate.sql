-- name: FindShippingRate :one
SELECT fee, free_shipping_min_subtotal
FROM shipping_rate
WHERE is_active = 1
  AND (province_id = sqlc.arg(province_id) OR province_id IS NULL)
  AND min_weight_g <= CAST(sqlc.arg(weight_g) AS SIGNED)
  AND (
      max_weight_g IS NULL
      OR CAST(sqlc.arg(weight_g) AS SIGNED) < max_weight_g
  )
ORDER BY (province_id IS NULL) ASC, min_weight_g DESC
LIMIT 1;