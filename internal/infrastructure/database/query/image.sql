-- name: ListImagesByReference :many
SELECT image_id, reference_id, reference_type, image_url, is_primary, sort_image, created_at, created_by, updated_at, updated_by
FROM image
WHERE reference_type = ? AND reference_id = ?
ORDER BY is_primary DESC, sort_image ASC, created_at ASC;

-- name: GetPrimaryImageByReference :one
SELECT image_id, reference_id, reference_type, image_url, is_primary, sort_image, created_at, created_by, updated_at, updated_by
FROM image
WHERE reference_type = ? AND reference_id = ? AND is_primary = 1
LIMIT 1;

-- name: AddImage :execresult
INSERT INTO image (reference_id, reference_type, image_url, is_primary, sort_image, created_at, created_by)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?);

-- name: UpdateImage :execresult
UPDATE image
SET image_url = ?,
    is_primary = ?,
    sort_image = ?,
    updated_at = CURRENT_TIMESTAMP,
    updated_by = ?
WHERE image_id = ?;

-- name: DeleteImage :execresult
DELETE FROM image WHERE image_id = ?;

-- name: SetPrimaryImage :execresult
-- ตั้งรูปใดรูปหนึ่งเป็น primary ของ reference
UPDATE image
SET is_primary = CASE WHEN image_id = ? THEN 1 ELSE 0 END
WHERE reference_type = ? AND reference_id = ?;
