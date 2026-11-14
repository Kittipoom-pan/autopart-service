-- name: GetPartBrandName :one
SELECT name FROM part_brand WHERE part_brand_id = ?;
