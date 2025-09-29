-- name: GetPartByID :one
SELECT part_id, part_brand_id, part_type_id, sku, name, description, price, quantity, is_active
FROM part
WHERE part_id = ? AND is_active = 1;

-- name: GetPartBySKU :one
SELECT part_id, part_brand_id, part_type_id, sku, name, description, price, quantity, is_active
FROM part
WHERE sku = ? AND is_active = 1;

-- name: ListParts :many
SELECT part_id, part_brand_id, part_type_id, sku, name, description, price, quantity, is_active
FROM part
WHERE is_active = 1
ORDER BY created_at;

-- name: GetPartWithImages :many
SELECT 
    p.part_id,
    p.sku,
    p.name AS part_name,
    p.description AS part_description,
    p.price,
    p.quantity,
    cm.car_model_id,
    cm.name AS car_model_name,
    COALESCE(cc.year_from, cm.year_from) AS year_from,
    COALESCE(cc.year_to, cm.year_to) AS year_to,
    i.image_id,
    i.image_url,
    i.is_primary,
    i.sort_image
FROM part p
LEFT JOIN compatible_car cc ON p.part_id = cc.part_id
LEFT JOIN car_model cm ON cc.car_model_id = cm.car_model_id
LEFT JOIN image i ON i.reference_type = 'part' AND i.reference_id = p.part_id
WHERE p.part_id = ?
ORDER BY i.is_primary DESC, i.sort_image ASC, i.created_at ASC;

-- name: CreatePart :execresult
INSERT INTO part (part_id, part_brand_id, part_type_id, sku, name, description, price, quantity, created_at, created_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdatePartByID :execresult
UPDATE part SET
    sku = ?,
    name = ?,
    description = ?,
    part_type_id = ?,
    part_brand_id = ?,
    price = ?,
    quantity = ?,
    is_active = ?,
    updated_at = ?,
    updated_by = ?
WHERE part_id = ?;

-- name: UpdatePartStockByID :execresult
UPDATE part SET
    quantity = ?,
    updated_at = ?,
    updated_by = ?
WHERE part_id = ?;

-- name: DeletePartByID :execresult
UPDATE part SET
    is_active = 0,
    updated_at = ?,
    updated_by = ?
WHERE part_id = ?;
