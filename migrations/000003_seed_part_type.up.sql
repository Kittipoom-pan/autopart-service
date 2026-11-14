-- ======================================
-- Up migration: part_type seed (with code)
-- ======================================
-- code ใช้สร้าง SKU: A-Z 0-9 เท่านั้น, ไม่เกิน 8 ตัว, unique

-- ---------- Level 1: root ----------
INSERT INTO part_type (parent_id, code, name) VALUES
  (NULL, 'CARCARE',  'การดูแลรถยนต์และอุปกรณ์'),
  (NULL, 'ENGTRANS', 'เครื่องยนต์และระบบส่งกำลัง'),
  (NULL, 'SUSPBRK',  'ระบบช่วงล่างและเบรก'),
  (NULL, 'ELECTRIC', 'ระบบไฟฟ้าและอุปกรณ์อิเล็กทรอนิกส์'),
  (NULL, 'COOLING',  'ระบบระบายความร้อนและน้ำมัน'),
  (NULL, 'FUELEXH',  'ระบบเชื้อเพลิงและไอเสีย'),
  (NULL, 'AIRCON',   'ระบบปรับอากาศและความเย็น');

-- ---------- Level 2 ----------
-- 1. การดูแลรถยนต์และอุปกรณ์
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'FILTER' AS code, 'หมวดไส้กรอง' AS name
  UNION ALL
  SELECT 'DIYCARE', 'หมวดดูแลรถยนต์ DIY'
  UNION ALL
  SELECT 'ACCESS', 'หมวดอุปกรณ์ดูแลรถยนต์'
) AS v WHERE p.code = 'CARCARE';

-- 2. เครื่องยนต์และระบบส่งกำลัง
-- หมายเหตุ: seed เดิมไม่มี 2 แถวนี้ (อนุมานจากกลุ่มลูก 4 + 2 รายการ) ตรวจชื่อก่อนใช้งานจริง
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'ENGINE' AS code, 'เครื่องยนต์' AS name
  UNION ALL
  SELECT 'TRANS', 'ระบบส่งกำลัง'
) AS v WHERE p.code = 'ENGTRANS';

-- 3. ระบบช่วงล่างและเบรก
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'SUSP' AS code, 'ช่วงล่างรถยนต์' AS name
  UNION ALL
  SELECT 'BRAKE', 'ระบบเบรก'
) AS v WHERE p.code = 'SUSPBRK';

-- 4. ระบบไฟฟ้าและอุปกรณ์อิเล็กทรอนิกส์
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'BATTERY' AS code, 'แบตเตอรี่และชาร์จไฟ' AS name
  UNION ALL
  SELECT 'LIGHTING', 'ระบบไฟส่องสว่างรถยนต์'
  UNION ALL
  SELECT 'SENSOR', 'เซ็นเซอร์และกล่องควบคุม'
) AS v WHERE p.code = 'ELECTRIC';

-- 5. ระบบระบายความร้อนและน้ำมัน
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'RADPUMP' AS code, 'หม้อน้ำและปั๊มน้ำ' AS name
  UNION ALL
  SELECT 'LUBE', 'น้ำมันเครื่องและสารหล่อลื่น'
) AS v WHERE p.code = 'COOLING';

-- 6. ระบบเชื้อเพลิงและไอเสีย
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'FUELSYS' AS code, 'ระบบเชื้อเพลิง' AS name
  UNION ALL
  SELECT 'EXHAUST', 'ระบบไอเสีย'
) AS v WHERE p.code = 'FUELEXH';

-- 7. ระบบปรับอากาศและความเย็น
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'COMPCOIL' AS code, 'คอมเพรสเซอร์และคอยล์' AS name
  UNION ALL
  SELECT 'ACFLTGAS', 'แอร์และไส้กรอง'
) AS v WHERE p.code = 'AIRCON';

-- ---------- Level 3 ----------
-- หมวดไส้กรอง
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'OILFLT' AS code, 'ไส้กรองน้ำมันเครื่อง' AS name
  UNION ALL
  SELECT 'FUELFLT', 'ไส้กรองน้ำมันเชื้อเพลิง'
  UNION ALL
  SELECT 'AIRFLT', 'ไส้กรองอากาศ'
  UNION ALL
  SELECT 'TRANSFLT', 'ไส้กรองเกียร์'
  UNION ALL
  SELECT 'ACFLT', 'ไส้กรองแอร์'
) AS v WHERE p.code = 'FILTER';

-- หมวดดูแลรถยนต์ DIY
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'CARWASH' AS code, 'น้ำยาล้างรถยนต์' AS name
  UNION ALL
  SELECT 'PAINTCT', 'น้ำยาเคลือบสีรถยนต์'
  UNION ALL
  SELECT 'SCRATCH', 'น้ำยาลบรอยขีดข่วนบนรถ'
  UNION ALL
  SELECT 'SEATCARE', 'น้ำยาบำรุงและเคลือบเบาะ'
  UNION ALL
  SELECT 'DIYTOOL', 'อุปกรณ์ DIY รถยนต์'
) AS v WHERE p.code = 'DIYCARE';

-- หมวดอุปกรณ์ดูแลรถยนต์
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'WIPER' AS code, 'ใบปัดน้ำฝนรถยนต์' AS name
  UNION ALL
  SELECT 'GREASE', 'จารบีหล่อลื่น'
  UNION ALL
  SELECT 'CHEMICAL', 'เคมีภัณฑ์ดูแลรถ'
  UNION ALL
  SELECT 'BULB', 'หลอดไฟรถยนต์'
  UNION ALL
  SELECT 'EXTRABAT', 'แบตเตอรี่เสริม'
) AS v WHERE p.code = 'ACCESS';

-- เครื่องยนต์
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'ENGPUMP' AS code, 'ปั๊มน้ำมันเครื่องยนต์' AS name
  UNION ALL
  SELECT 'BELTCHN', 'สายพานและโซ่เครื่องยนต์'
  UNION ALL
  SELECT 'SPARKPLG', 'หัวเทียนเครื่องยนต์'
  UNION ALL
  SELECT 'ENGFLT', 'กรองน้ำมันเครื่องยนต์'
) AS v WHERE p.code = 'ENGINE';

-- ระบบส่งกำลัง
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'CLUTCH' AS code, 'คลัตช์และชิ้นส่วนเกียร์' AS name
  UNION ALL
  SELECT 'SHAFT', 'เพลาและลูกปืนเกียร์'
) AS v WHERE p.code = 'TRANS';

-- ช่วงล่างรถยนต์
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'SHOCK' AS code, 'โช้คอัพรถยนต์' AS name
  UNION ALL
  SELECT 'SPRING', 'สปริงและคอยล์สปริง'
  UNION ALL
  SELECT 'BALLJNT', 'ลูกหมากและบูชช่วงล่าง'
  UNION ALL
  SELECT 'STEERING', 'พวงมาลัยและระบบบังคับเลี้ยว'
) AS v WHERE p.code = 'SUSP';

-- ระบบเบรก
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'DISC' AS code, 'จานเบรก' AS name
  UNION ALL
  SELECT 'PAD', 'ผ้าเบรก'
  UNION ALL
  SELECT 'BRKLINE', 'สายเบรกและน้ำมันเบรก'
  UNION ALL
  SELECT 'ABSSNSR', 'เซนเซอร์ ABS'
) AS v WHERE p.code = 'BRAKE';

-- แบตเตอรี่และชาร์จไฟ
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'ALTERNTR' AS code, 'ไดชาร์จรถยนต์' AS name
  UNION ALL
  SELECT 'WIRING', 'สายไฟและฟิวส์'
) AS v WHERE p.code = 'BATTERY';

-- ระบบไฟส่องสว่างรถยนต์
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'HEADLAMP' AS code, 'หลอดไฟหน้า' AS name
  UNION ALL
  SELECT 'TAILLAMP', 'หลอดไฟท้าย'
  UNION ALL
  SELECT 'SIGNAL', 'ไฟเลี้ยวและสัญญาณ'
) AS v WHERE p.code = 'LIGHTING';

-- หม้อน้ำและปั๊มน้ำ
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'RADIATOR' AS code, 'หม้อน้ำรถยนต์' AS name
  UNION ALL
  SELECT 'WATRPUMP', 'ปั๊มน้ำรถยนต์'
) AS v WHERE p.code = 'RADPUMP';

-- น้ำมันเครื่องและสารหล่อลื่น
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'ENGOIL' AS code, 'น้ำมันเครื่องยนต์' AS name
  UNION ALL
  SELECT 'GEAROIL', 'น้ำมันเกียร์'
) AS v WHERE p.code = 'LUBE';

-- ระบบเชื้อเพลิง
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'FUELTANK' AS code, 'ถังน้ำมัน' AS name
  UNION ALL
  SELECT 'FUELPUMP', 'ปั๊มเชื้อเพลิง'
  UNION ALL
  SELECT 'INJECTOR', 'หัวฉีดเชื้อเพลิง'
) AS v WHERE p.code = 'FUELSYS';

-- คอมเพรสเซอร์และคอยล์
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'ACCOMP' AS code, 'คอมเพรสเซอร์ AC' AS name
  UNION ALL
  SELECT 'ACCONDSR', 'คอยล์ร้อน AC'
  UNION ALL
  SELECT 'ACEVAP', 'คอยล์เย็น AC'
) AS v WHERE p.code = 'COMPCOIL';

-- แอร์และไส้กรอง
INSERT INTO part_type (parent_id, code, name)
SELECT p.part_type_id, v.code, v.name FROM part_type p, (
  SELECT 'CABINFLT' AS code, 'ไส้กรองแอร์รถยนต์' AS name
  UNION ALL
  SELECT 'ACGAS', 'น้ำยาแอร์รถยนต์'
) AS v WHERE p.code = 'ACFLTGAS';