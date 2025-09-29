-- name: GetCarModel :one
SELECT car_model_id, car_brand_id, name, year_from, year_to
FROM car_model
WHERE car_model_id = ?;

-- name: ListCarModelsByBrand :many
SELECT car_model_id, car_brand_id, name, year_from, year_to
FROM car_model
WHERE car_brand_id = ?
ORDER BY name, year_from;

-- name: CreateCarModel :execresult
INSERT INTO car_model (car_brand_id, name, year_from, year_to)
VALUES (?, ?, ?, ?);

-- name: UpdateCarModel :execresult
UPDATE car_model SET
    car_brand_id = ?,
    name = ?,
    year_from = ?,
    year_to = ?
WHERE car_model_id = ?;

-- name: DeleteCarModel :execresult
DELETE FROM car_model WHERE car_model_id = ?;
