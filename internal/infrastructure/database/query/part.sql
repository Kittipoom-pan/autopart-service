-- name: GetPart :one
SELECT * FROM part WHERE part_id = ?;

-- name: ListParts :many
SELECT * FROM part ORDER BY name;