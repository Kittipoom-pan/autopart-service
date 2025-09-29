-- name: GetCompatibleCar :one
SELECT compatible_id, part_id, car_model_id, year_from, year_to
FROM compatible_car
WHERE compatible_id = ?;

-- name: ListCompatibleCarsByPart :many
SELECT cc.compatible_id, cc.part_id, cc.car_model_id,
       cm.name AS car_model_name, cm.year_from AS model_year_from, cm.year_to AS model_year_to,
       cb.name AS car_brand_name,
       cc.year_from AS override_year_from, cc.year_to AS override_year_to
FROM compatible_car cc
JOIN car_model cm ON cc.car_model_id = cm.car_model_id
JOIN car_brand cb ON cm.car_brand_id = cb.car_brand_id
WHERE cc.part_id = ?
ORDER BY cb.name, cm.name, cm.year_from;

-- name: CreateCompatibleCar :execresult
INSERT INTO compatible_car (part_id, car_model_id, year_from, year_to)
VALUES (?, ?, ?, ?);

-- name: DeleteCompatibleCar :execresult
DELETE FROM compatible_car WHERE compatible_id = ?;
