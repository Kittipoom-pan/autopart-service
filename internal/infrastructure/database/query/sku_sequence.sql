-- name: NextSkuSequence :execresult
INSERT INTO sku_sequence (prefix, last_number)
VALUES (?, LAST_INSERT_ID(1))
ON DUPLICATE KEY UPDATE last_number = LAST_INSERT_ID(last_number + 1);