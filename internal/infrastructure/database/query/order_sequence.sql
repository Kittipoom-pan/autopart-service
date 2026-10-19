-- name: NextOrderSequence :execresult
INSERT INTO order_sequence (seq_date, `last_value`)
VALUES (?, LAST_INSERT_ID(1))
ON DUPLICATE KEY UPDATE `last_value` = LAST_INSERT_ID(`last_value` + 1);