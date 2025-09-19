-- ======================================
-- Down migration for part_type seed
-- ======================================

-- 1. ลบ sub-sub-category (ลูกหลานของ sub-category)
DELETE pt3
FROM part_type pt3
JOIN part_type pt2 ON pt3.parent_id = pt2.part_type_id
JOIN part_type pt1 ON pt2.parent_id = pt1.part_type_id;

-- 2. ลบ sub-category (ลูกของ parent)
DELETE pt2
FROM part_type pt2
JOIN part_type pt1 ON pt2.parent_id = pt1.part_type_id;

-- 3. ลบ parent category
DELETE FROM part_type WHERE parent_id IS NULL;
