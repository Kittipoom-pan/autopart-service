-- ======================================
-- Up migration: part_type seed (unique names)
-- ======================================

-- 1. การดูแลรถยนต์และอุปกรณ์
INSERT INTO part_type (parent_id, name) VALUES (NULL, 'การดูแลรถยนต์และอุปกรณ์');

INSERT INTO part_type (parent_id, name) VALUES (1, 'หมวดไส้กรอง');
INSERT INTO part_type (parent_id, name) VALUES (1, 'หมวดดูแลรถยนต์ DIY');
INSERT INTO part_type (parent_id, name) VALUES (1, 'หมวดอุปกรณ์ดูแลรถยนต์');

INSERT INTO part_type (parent_id, name) VALUES (2, 'ไส้กรองน้ำมันเครื่อง');
INSERT INTO part_type (parent_id, name) VALUES (2, 'ไส้กรองน้ำมันเชื้อเพลิง');
INSERT INTO part_type (parent_id, name) VALUES (2, 'ไส้กรองอากาศ');
INSERT INTO part_type (parent_id, name) VALUES (2, 'ไส้กรองเกียร์');
INSERT INTO part_type (parent_id, name) VALUES (2, 'ไส้กรองแอร์');

INSERT INTO part_type (parent_id, name) VALUES (3, 'น้ำยาล้างรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (3, 'น้ำยาเคลือบสีรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (3, 'น้ำยาลบรอยขีดข่วนบนรถ');
INSERT INTO part_type (parent_id, name) VALUES (3, 'น้ำยาบำรุงและเคลือบเบาะ');
INSERT INTO part_type (parent_id, name) VALUES (3, 'อุปกรณ์ DIY รถยนต์');

INSERT INTO part_type (parent_id, name) VALUES (4, 'ใบปัดน้ำฝนรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (4, 'จารบีหล่อลื่น');
INSERT INTO part_type (parent_id, name) VALUES (4, 'เคมีภัณฑ์ดูแลรถ');
INSERT INTO part_type (parent_id, name) VALUES (4, 'หลอดไฟรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (4, 'แบตเตอรี่เสริม');

-- 2. เครื่องยนต์และระบบส่งกำลัง
INSERT INTO part_type (parent_id, name) VALUES (NULL, 'เครื่องยนต์และระบบส่งกำลัง');

INSERT INTO part_type (parent_id, name) VALUES (6, 'ปั๊มน้ำมันเครื่องยนต์');
INSERT INTO part_type (parent_id, name) VALUES (6, 'สายพานและโซ่เครื่องยนต์');
INSERT INTO part_type (parent_id, name) VALUES (6, 'หัวเทียนเครื่องยนต์');
INSERT INTO part_type (parent_id, name) VALUES (6, 'กรองน้ำมันเครื่องยนต์');

INSERT INTO part_type (parent_id, name) VALUES (7, 'คลัตช์และชิ้นส่วนเกียร์');
INSERT INTO part_type (parent_id, name) VALUES (7, 'เพลาและลูกปืนเกียร์');

-- 3. ระบบช่วงล่างและเบรก
INSERT INTO part_type (parent_id, name) VALUES (NULL, 'ระบบช่วงล่างและเบรก');

INSERT INTO part_type (parent_id, name) VALUES (8, 'ช่วงล่างรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (8, 'ระบบเบรก');

INSERT INTO part_type (parent_id, name) VALUES (9, 'โช้คอัพรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (9, 'สปริงและคอยล์สปริง');
INSERT INTO part_type (parent_id, name) VALUES (9, 'ลูกหมากและบูชช่วงล่าง');
INSERT INTO part_type (parent_id, name) VALUES (9, 'พวงมาลัยและระบบบังคับเลี้ยว');

INSERT INTO part_type (parent_id, name) VALUES (10, 'จานเบรก');
INSERT INTO part_type (parent_id, name) VALUES (10, 'ผ้าเบรก');
INSERT INTO part_type (parent_id, name) VALUES (10, 'สายเบรกและน้ำมันเบรก');
INSERT INTO part_type (parent_id, name) VALUES (10, 'เซนเซอร์ ABS');

-- 4. ระบบไฟฟ้าและอุปกรณ์อิเล็กทรอนิกส์
INSERT INTO part_type (parent_id, name) VALUES (NULL, 'ระบบไฟฟ้าและอุปกรณ์อิเล็กทรอนิกส์');

INSERT INTO part_type (parent_id, name) VALUES (11, 'แบตเตอรี่และชาร์จไฟ');
INSERT INTO part_type (parent_id, name) VALUES (11, 'ระบบไฟส่องสว่างรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (11, 'เซ็นเซอร์และกล่องควบคุม');

INSERT INTO part_type (parent_id, name) VALUES (12, 'ไดชาร์จรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (12, 'สายไฟและฟิวส์');

INSERT INTO part_type (parent_id, name) VALUES (13, 'หลอดไฟหน้า');
INSERT INTO part_type (parent_id, name) VALUES (13, 'หลอดไฟท้าย');
INSERT INTO part_type (parent_id, name) VALUES (13, 'ไฟเลี้ยวและสัญญาณ');

-- 5. ระบบระบายความร้อนและน้ำมัน
INSERT INTO part_type (parent_id, name) VALUES (NULL, 'ระบบระบายความร้อนและน้ำมัน');

INSERT INTO part_type (parent_id, name) VALUES (26, 'หม้อน้ำและปั๊มน้ำ');
INSERT INTO part_type (parent_id, name) VALUES (26, 'น้ำมันเครื่องและสารหล่อลื่น');

INSERT INTO part_type (parent_id, name) VALUES (27, 'หม้อน้ำรถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (27, 'ปั๊มน้ำรถยนต์');

INSERT INTO part_type (parent_id, name) VALUES (28, 'น้ำมันเครื่องยนต์');
INSERT INTO part_type (parent_id, name) VALUES (28, 'น้ำมันเกียร์');

-- 6. ระบบเชื้อเพลิงและไอเสีย
INSERT INTO part_type (parent_id, name) VALUES (NULL, 'ระบบเชื้อเพลิงและไอเสีย');

INSERT INTO part_type (parent_id, name) VALUES (29, 'ระบบเชื้อเพลิง');
INSERT INTO part_type (parent_id, name) VALUES (29, 'ระบบไอเสีย');

INSERT INTO part_type (parent_id, name) VALUES (30, 'ถังน้ำมัน');
INSERT INTO part_type (parent_id, name) VALUES (30, 'ปั๊มเชื้อเพลิง');
INSERT INTO part_type (parent_id, name) VALUES (30, 'หัวฉีดเชื้อเพลิง');

-- 7. ระบบปรับอากาศและความเย็น
INSERT INTO part_type (parent_id, name) VALUES (NULL, 'ระบบปรับอากาศและความเย็น');

INSERT INTO part_type (parent_id, name) VALUES (36, 'คอมเพรสเซอร์และคอยล์');
INSERT INTO part_type (parent_id, name) VALUES (36, 'แอร์และไส้กรอง');

INSERT INTO part_type (parent_id, name) VALUES (37, 'คอมเพรสเซอร์ AC');
INSERT INTO part_type (parent_id, name) VALUES (37, 'คอยล์ร้อน AC');
INSERT INTO part_type (parent_id, name) VALUES (37, 'คอยล์เย็น AC');

INSERT INTO part_type (parent_id, name) VALUES (38, 'ไส้กรองแอร์รถยนต์');
INSERT INTO part_type (parent_id, name) VALUES (38, 'น้ำยาแอร์รถยนต์');
