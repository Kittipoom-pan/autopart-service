-- name: GetPartWithCompatibilityByID :one
SELECT p.part_id, p.sku, p.name AS part_name, p.description, p.price, p.quantity,
       cm.car_model_id, cm.name AS car_model_name, cm.year_from, cm.year_to,
       cb.car_brand_id, cb.name AS car_brand_name
FROM part p
JOIN compatible_car cc ON p.part_id = cc.part_id
JOIN car_model cm ON cc.car_model_id = cm.car_model_id
JOIN car_brand cb ON cm.car_brand_id = cb.car_brand_id
WHERE p.part_id = ?;

-- name: SearchPartsByCar :many
SELECT DISTINCT p.part_id, p.sku, p.name AS part_name, p.description, p.price, p.quantity,
       cb.name AS car_brand_name, cm.name AS car_model_name
FROM part p
JOIN compatible_car cc ON p.part_id = cc.part_id
JOIN car_model cm ON cc.car_model_id = cm.car_model_id
JOIN car_brand cb ON cm.car_brand_id = cb.car_brand_id
WHERE cb.name LIKE CONCAT('%', ?, '%')
  AND cm.name LIKE CONCAT('%', ?, '%')
  AND (? IS NULL OR cm.year_from <= ?)
  AND (? IS NULL OR cm.year_to >= ?)
ORDER BY p.name;

-- -- name: SearchPartsByCar :many
-- SELECT DISTINCT 
--     p.part_id, p.sku, p.name AS part_name, p.description, p.price, p.stock_quantity,
--     cb.name AS car_brand_name, cm.name AS car_model_name,
--     COALESCE(cc.year_from, cm.year_from) AS year_from,
--     COALESCE(cc.year_to, cm.year_to) AS year_to
-- FROM part p
-- JOIN compatible_car cc ON p.part_id = cc.part_id
-- JOIN car_model cm ON cc.car_model_id = cm.car_model_id
-- JOIN car_brand cb ON cm.car_brand_id = cb.car_brand_id
-- WHERE cb.name LIKE CONCAT('%', ?, '%')
--   AND cm.name LIKE CONCAT('%', ?, '%')
--   AND (
--        -- ตรวจสอบปี โดยเลือกใช้ override ถ้ามี
--        (? IS NULL OR (? BETWEEN COALESCE(cc.year_from, cm.year_from) AND COALESCE(cc.year_to, cm.year_to)))
--       )
-- ORDER BY p.name;

-- -- name: GetPartWithCompatibility :one
-- SELECT p.part_id, p.sku, p.name AS part_name, p.description, p.price, p.stock_quantity,
--        cm.car_model_id, cm.name AS car_model_name,
--        cb.car_brand_id, cb.name AS car_brand_name,
--        COALESCE(cc.year_from, cm.year_from) AS year_from,
--        COALESCE(cc.year_to, cm.year_to) AS year_to
-- FROM part p
-- JOIN compatible_car cc ON p.part_id = cc.part_id
-- JOIN car_model cm ON cc.car_model_id = cm.car_model_id
-- JOIN car_brand cb ON cm.car_brand_id = cb.car_brand_id
-- WHERE p.part_id = ?;