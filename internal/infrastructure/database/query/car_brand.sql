-- name: GetCarBrand :one
SELECT car_brand_id, name, description
FROM car_brand
WHERE car_brand_id = ?;

-- name: ListCarBrands :many
SELECT car_brand_id, name, description
FROM car_brand
ORDER BY name;

-- name: CreateCarBrand :execresult
INSERT INTO car_brand (name, description)
VALUES (?, ?);

-- name: UpdateCarBrand :execresult
UPDATE car_brand SET
    name = ?,
    description = ?
WHERE car_brand_id = ?;

-- name: DeleteCarBrand :execresult
DELETE FROM car_brand WHERE car_brand_id = ?;
