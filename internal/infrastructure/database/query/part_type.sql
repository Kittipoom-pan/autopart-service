-- name: GetPartTypeCode :one
SELECT code FROM part_type WHERE part_type_id = ?;
