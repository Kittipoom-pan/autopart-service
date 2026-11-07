-- ============================================================
-- Autopart schema (ออกแบบใหม่)
-- ต้องใช้ MySQL 8.0.16+ (CHECK constraint) และ stored generated column (ใช้กับ cart / payment)
--
-- หลักการ audit
--   *_admin_id  = admin_user_id ของแอดมินที่ทำรายการ มี FK ไป admin_user เสมอ
--                 NULL = ไม่ได้ทำโดยแอดมิน (ลูกค้าทำเอง หรือระบบ)
--   ตารางที่ลูกค้า/ระบบเขียนอย่างเดียว (cart, cart_item, token, order_item ...) ไม่มี audit
--   ผู้เปลี่ยนสถานะ order เก็บใน order_status_history (แยกประเภท admin / customer / system)
--
-- หลักการ soft delete + unique (ใช้ UNIQUE ธรรมดา ไม่ใช้ generated column)
--   admin_user : ปิดบัญชีด้วย is_active = 0 (soft delete) username/email ถูกจองตลอด
--   customer   : is_active = 0 = ระงับชั่วคราว (ข้อมูลยังอยู่ เปิดกลับได้ ค่ายังถูกจอง)
--                deleted_at มีค่า = ลบบัญชีถาวร: ล้างข้อมูลส่วนตัวทั้งหมดเป็น NULL
--                UNIQUE ยอมให้ NULL ซ้ำ ค่าเดิม (อีเมล/เบอร์/username) จึงว่างให้สมัครใหม่ได้
--                ขั้นตอนลบดู customer_anonymize.sql
-- ============================================================

-- ================================
-- 1. Admin & master data
-- ================================

CREATE TABLE admin_user (
  admin_user_id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(40) NOT NULL,
  email VARCHAR(60) NULL,
  password VARCHAR(150) NOT NULL,
  role ENUM('super_admin', 'staff') NOT NULL,
  is_active BIT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,                       -- NULL = แอดมินคนแรก (seed)
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  UNIQUE KEY uq_admin_username (username),
  UNIQUE KEY uq_admin_email (email),                  -- email เป็น NULL ได้ (NULL ไม่ชนกัน)
  CONSTRAINT fk_admin_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_admin_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- แอดมินแก้ master data ได้ จึงเก็บ audit
CREATE TABLE part_type (
  part_type_id INT AUTO_INCREMENT PRIMARY KEY,
  parent_id INT NULL,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  CONSTRAINT fk_part_type_parent FOREIGN KEY (parent_id) REFERENCES part_type(part_type_id) ON DELETE SET NULL,
  CONSTRAINT fk_part_type_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_part_type_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE part_brand (
  part_brand_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  CONSTRAINT fk_part_brand_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_part_brand_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE car_brand (
  car_brand_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  CONSTRAINT fk_car_brand_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_car_brand_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE car_model (
  car_model_id INT AUTO_INCREMENT PRIMARY KEY,
  car_brand_id INT NOT NULL,
  name VARCHAR(50) NOT NULL,
  year_from YEAR DEFAULT NULL,
  year_to YEAR DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  CONSTRAINT fk_car_model_car_brand FOREIGN KEY (car_brand_id) REFERENCES car_brand(car_brand_id),
  CONSTRAINT fk_car_model_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_car_model_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE part (
  part_id INT AUTO_INCREMENT PRIMARY KEY,
  part_brand_id INT NOT NULL,
  part_type_id INT NOT NULL,
  name VARCHAR(100) NOT NULL UNIQUE,
  sku VARCHAR(50) NOT NULL UNIQUE,
  description TEXT NULL,
  price INT NOT NULL DEFAULT 0,
  quantity INT NOT NULL DEFAULT 0,                    -- ต้องตรงกับผลรวมของ stock_movement
  weight_g INT NULL,
  is_active BIT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,                       -- การขายไม่ใช่การแก้ part: ประวัติสต็อกดูที่ stock_movement
  CONSTRAINT fk_part_part_brand FOREIGN KEY (part_brand_id) REFERENCES part_brand(part_brand_id),
  CONSTRAINT fk_part_part_type FOREIGN KEY (part_type_id) REFERENCES part_type(part_type_id),
  CONSTRAINT fk_part_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_part_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT chk_part_values CHECK (quantity >= 0 AND price >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE compatible_car (
  compatible_id INT AUTO_INCREMENT PRIMARY KEY,
  part_id INT NOT NULL,
  car_model_id INT NOT NULL,
  year_from YEAR DEFAULT NULL,                        -- NULL = ไม่จำกัดปี
  year_to YEAR DEFAULT NULL,
  CONSTRAINT fk_compatible_car_part FOREIGN KEY (part_id) REFERENCES part(part_id) ON DELETE CASCADE,
  CONSTRAINT fk_compatible_car_model FOREIGN KEY (car_model_id) REFERENCES car_model(car_model_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE discount (
  discount_id INT AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(60) NOT NULL UNIQUE,
  discount_type ENUM('percentage', 'fixed') NOT NULL,
  discount_value INT NOT NULL,
  description VARCHAR(255) NULL,
  min_order_amount INT NOT NULL DEFAULT 0,
  expiration_date DATE NULL,                          -- NULL = ไม่หมดอายุ
  is_active BIT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  CONSTRAINT fk_discount_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_discount_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT chk_discount_values CHECK (
    discount_value > 0
    AND min_order_amount >= 0
    AND (discount_type <> 'percentage' OR discount_value <= 100))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ค่าคงที่ที่ seed ไว้ (credit_card, cod, bank_transfer, cash) ไม่ต้องมี audit
CREATE TABLE payment_method_type (
  payment_method_type_id INT AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(20) NOT NULL UNIQUE,
  name VARCHAR(40) NOT NULL,
  is_active BIT(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ข้อมูลที่อยู่จากทางการ (import ครั้งเดียว) ไม่ต้องมี audit
CREATE TABLE province (
  id INT AUTO_INCREMENT PRIMARY KEY,
  code INT NOT NULL,
  name_in_thai VARCHAR(150) NOT NULL,
  name_in_english VARCHAR(150) NOT NULL,
  UNIQUE KEY ux_province_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE district (
  id INT AUTO_INCREMENT PRIMARY KEY,
  code INT NOT NULL,
  name_in_thai VARCHAR(150) NOT NULL,
  name_in_english VARCHAR(150) NOT NULL,
  province_id INT NOT NULL,
  UNIQUE KEY ux_district_code (code),
  KEY ix_district_province_id (province_id),
  CONSTRAINT fk_district_province FOREIGN KEY (province_id) REFERENCES province(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE subdistrict (
  id INT AUTO_INCREMENT PRIMARY KEY,
  code INT NOT NULL,
  name_in_thai VARCHAR(150) NOT NULL,
  name_in_english VARCHAR(150) NULL,
  district_id INT NOT NULL,
  zip_code CHAR(5) NOT NULL,
  UNIQUE KEY ux_subdistrict_code (code),
  KEY ix_subdistrict_district_id (district_id),
  CONSTRAINT fk_subdistrict_district FOREIGN KEY (district_id) REFERENCES district(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 2. Customer
-- ================================

-- แอดมินสร้าง/แก้ลูกค้าได้ (เช่น คีย์ออเดอร์ทางโทรศัพท์) จึงเก็บ audit ผู้กระทำที่เป็นแอดมิน
-- คอลัมน์ข้อมูลส่วนตัวเป็น NULL ได้ เพื่อให้ล้างทิ้งตอน anonymize
-- แต่ chk_customer_deletedบังคับว่า: บัญชีปกติต้องมีครบ, บัญชีที่ลบแล้วต้องว่างทั้งหมด
-- (ชื่อที่แสดงในประวัติออเดอร์ของบัญชีที่ลบแล้ว ให้แอปแสดงข้อความ เช่น "ลูกค้าที่ถูกลบ")
CREATE TABLE customer (
  customer_id INT AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(40) NULL,
  last_name VARCHAR(40) NULL,
  username VARCHAR(50) NULL,
  email VARCHAR(50) NULL,
  password VARCHAR(150) NULL,
  birth_date DATE NULL,
  phone_number VARCHAR(15) NULL,
  is_active BIT(1) NOT NULL DEFAULT 1,
  deleted_at TIMESTAMP NULL,                       -- NULL = บัญชีปกติ, มีค่า = ลบบัญชีถาวรแล้ว
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,                       -- NULL = สมัครเอง
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,                       -- NULL = ลูกค้าแก้เอง (หรือยังไม่เคยแก้ ดู updated_at)
  UNIQUE KEY uq_customer_username (username),
  UNIQUE KEY uq_customer_email (email),
  UNIQUE KEY uq_customer_phone (phone_number),
  CONSTRAINT fk_customer_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_customer_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT chk_customer_deleted CHECK (
    (deleted_at IS NULL
       AND first_name IS NOT NULL AND last_name IS NOT NULL AND username IS NOT NULL
       AND email IS NOT NULL AND password IS NOT NULL AND phone_number IS NOT NULL)
    OR (deleted_at IS NOT NULL AND is_active = 0
       AND first_name IS NULL AND last_name IS NULL AND username IS NULL
       AND email IS NULL AND password IS NULL AND phone_number IS NULL AND birth_date IS NULL))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE customer_payment_method (
  customer_payment_method_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  account_name VARCHAR(60) NULL,
  payment_method_type_id INT NOT NULL,
  card_token VARCHAR(255) NULL UNIQUE,
  is_default BIT(1) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  CONSTRAINT fk_cpm_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT fk_cpm_payment_type FOREIGN KEY (payment_method_type_id) REFERENCES payment_method_type(payment_method_type_id),
  CONSTRAINT fk_cpm_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_cpm_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE address (
  address_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  title VARCHAR(50) NULL,                             -- ชื่อเรียกที่อยู่ เช่น บ้าน, ที่ทำงาน
  recipient_name VARCHAR(100) NOT NULL,
  address_line VARCHAR(255) NOT NULL,
  subdistrict_id INT NOT NULL,
  phone_number VARCHAR(15) NOT NULL,
  is_default BIT(1) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  INDEX idx_address_customer (customer_id),
  INDEX idx_address_subdistrict (subdistrict_id),
  CONSTRAINT fk_address_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT fk_address_subdistrict FOREIGN KEY (subdistrict_id) REFERENCES subdistrict(id),
  CONSTRAINT fk_address_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_address_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 3. Cart (ลูกค้าเขียนเองเท่านั้น: ไม่มี audit ผู้กระทำ)
-- ================================

CREATE TABLE cart (
  cart_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  subtotal INT NOT NULL DEFAULT 0,                    -- ราคารวมสินค้า (ก่อนหักส่วนลด ยังไม่รวมค่าส่ง)
  is_checkout BIT(1) NOT NULL DEFAULT 0,
  -- ลูกค้ามี active cart ได้ใบเดียว (cart ที่ checkout แล้วได้ NULL จึงมีกี่ใบก็ได้)
  active_customer_id INT GENERATED ALWAYS AS (IF(is_checkout = 0, customer_id, NULL)) STORED,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_cart_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT uq_cart_one_active UNIQUE (active_customer_id),
  CONSTRAINT chk_cart_subtotal CHECK (subtotal >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE cart_item (
  cart_item_id INT AUTO_INCREMENT PRIMARY KEY,
  cart_id INT NOT NULL,
  part_id INT NOT NULL,
  quantity INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_cart_item_cart FOREIGN KEY (cart_id) REFERENCES cart(cart_id) ON DELETE CASCADE,
  CONSTRAINT fk_cart_item_part FOREIGN KEY (part_id) REFERENCES part(part_id),
  CONSTRAINT uq_cart_item_part UNIQUE (cart_id, part_id),
  CONSTRAINT chk_cart_item_quantity CHECK (quantity > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 4. Token
-- ================================

CREATE TABLE admin_token (
  id INT AUTO_INCREMENT PRIMARY KEY,
  admin_user_id INT NOT NULL,
  token VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_revoked BIT(1) NOT NULL DEFAULT 0,
  CONSTRAINT fk_token_admin_user FOREIGN KEY (admin_user_id) REFERENCES admin_user(admin_user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE customer_token (
  id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  token VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_revoked BIT(1) NOT NULL DEFAULT 0,
  CONSTRAINT fk_token_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 5. Order
-- ================================

-- channel: web/mobile = ลูกค้าสั่งเอง, call/line/facebook/walk_in = แอดมินบันทึกให้
-- หน้าร้าน: ให้แอปตั้ง status = 'completed' และ payment.status = 'completed' ทันที
-- ผู้เปลี่ยนสถานะทุกครั้งเก็บใน order_status_history
CREATE TABLE `order` (
  order_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_number VARCHAR(13) NOT NULL UNIQUE,
  customer_id INT NULL,
  channel ENUM('web', 'mobile', 'call', 'line', 'facebook', 'walk_in') NOT NULL DEFAULT 'web',
  placed_by_admin_id INT NULL,                        -- NULL = ลูกค้าสั่งเอง, มีค่า = แอดมินสั่งแทน
  cart_id INT NULL,                                   -- ออเดอร์ที่แอดมินคีย์มักไม่มี cart
  subtotal INT NOT NULL DEFAULT 0,                    -- ก่อนหักส่วนลด ยังไม่รวมค่าส่ง
  discount_amount INT NOT NULL DEFAULT 0,
  shipping_fee INT NOT NULL DEFAULT 0,
  total_amount INT NOT NULL DEFAULT 0,                -- subtotal - discount + shipping
  discount_id INT NULL,
  status ENUM('pending', 'processing', 'shipped', 'delivered', 'cancelled', 'completed') NOT NULL DEFAULT 'pending',
  remark VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_order_customer_created (customer_id, created_at),
  INDEX idx_order_status_created (status, created_at),
  UNIQUE KEY uq_order_cart (cart_id),
  CONSTRAINT fk_order_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT fk_order_placed_by_admin FOREIGN KEY (placed_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_order_cart FOREIGN KEY (cart_id) REFERENCES cart(cart_id),
  CONSTRAINT fk_order_discount FOREIGN KEY (discount_id) REFERENCES discount(discount_id),
  CONSTRAINT chk_order_amounts CHECK (
    subtotal >= 0 AND discount_amount >= 0 AND shipping_fee >= 0
    AND total_amount = subtotal - discount_amount + shipping_fee),
  -- ต้องมี customer_id ทุกช่องทาง ยกเว้น walk_in
  CONSTRAINT chk_order_walkin CHECK (channel = 'walk_in' OR customer_id IS NOT NULL),
  -- web/mobile ต้องไม่มีแอดมิน, ช่องทางอื่นต้องมีแอดมินเป็นผู้บันทึก
  CONSTRAINT chk_order_channel_admin CHECK (
    (channel IN ('web', 'mobile') AND placed_by_admin_id IS NULL)
    OR (channel NOT IN ('web', 'mobile') AND placed_by_admin_id IS NOT NULL))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE order_item (
  order_item_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  part_id INT NOT NULL,
  quantity INT NOT NULL,
  price INT NOT NULL,                                 -- ราคา ณ ตอนสั่งซื้อ
  CONSTRAINT fk_order_item_order FOREIGN KEY (order_id) REFERENCES `order`(order_id) ON DELETE CASCADE,
  CONSTRAINT fk_order_item_part FOREIGN KEY (part_id) REFERENCES part(part_id),
  CONSTRAINT uq_order_item_part UNIQUE (order_id, part_id),
  CONSTRAINT chk_order_item_values CHECK (quantity > 0 AND price >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE order_sequence (
  seq_date DATE NOT NULL PRIMARY KEY,
  `last_value` INT UNSIGNED NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- snapshot ที่อยู่ ณ ตอนสั่ง (ไม่ผูก master เพราะชื่อ/รหัสไปรษณีย์อาจถูกแก้ภายหลัง)
CREATE TABLE order_shipping_address (
  order_id BIGINT PRIMARY KEY,
  recipient_name VARCHAR(100) NOT NULL,
  phone_number VARCHAR(15) NOT NULL,
  address_line VARCHAR(255) NOT NULL,
  subdistrict VARCHAR(100) NOT NULL,
  district VARCHAR(100) NOT NULL,
  province VARCHAR(100) NOT NULL,
  zip_code CHAR(5) NOT NULL,
  CONSTRAINT fk_order_shipping FOREIGN KEY (order_id) REFERENCES `order`(order_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- แอดมินตั้งค่าค่าส่ง จึงเก็บ audit
CREATE TABLE shipping_rate (
  shipping_rate_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  province_id INT NULL COMMENT 'NULL = ใช้ได้ทุกจังหวัด',
  min_weight_g INT NOT NULL DEFAULT 0 COMMENT 'รวมค่านี้ (>=)',
  max_weight_g INT NULL COMMENT 'ไม่รวมค่านี้ (<); NULL = ไม่จำกัด',
  fee INT NOT NULL COMMENT 'ค่าจัดส่ง',
  free_shipping_min_subtotal INT NULL COMMENT 'ยอดสินค้าหลังหักส่วนลดตั้งแต่เท่านี้ส่งฟรี NULL = ไม่มี',
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,
  KEY idx_shipping_lookup (is_active, province_id, min_weight_g),
  CONSTRAINT fk_shipping_rate_province FOREIGN KEY (province_id) REFERENCES province(id),
  CONSTRAINT fk_shipping_rate_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_shipping_rate_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT chk_shipping_fee CHECK (fee >= 0),
  CONSTRAINT chk_shipping_weight CHECK (max_weight_g IS NULL OR max_weight_g > min_weight_g)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 6. Stock & Payment
-- ================================

-- ขาย (ออนไลน์/หน้าร้าน) ต้องผ่าน order เสมอ: reason = 'sale' และมี order_id
-- ยกเลิกออเดอร์/คืนของ: reason = 'return' ผูก order_id เดิม
-- รับของเข้า/นับสต็อก/ของเสีย: ไม่ต้องมี order_id ให้ใส่ remark และ created_by_admin_id
-- ตัดสต็อกในธุรกรรมเดียวกับการสร้าง order:
--   UPDATE part SET quantity = quantity - ? WHERE part_id = ? AND quantity >= ?;
--   (affected rows = 0 แปลว่าสต็อกไม่พอ ให้ rollback)
CREATE TABLE stock_movement (
  stock_movement_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  part_id INT NOT NULL,
  quantity_change INT NOT NULL,
  price INT NOT NULL DEFAULT 0,
  event_type ENUM('in', 'out') NOT NULL,
  reason ENUM('sale', 'return', 'purchase', 'adjustment', 'damaged') NOT NULL,
  order_id BIGINT NULL,
  remark VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,                       -- NULL = ลูกค้า/ระบบ (ขายออนไลน์) ดู order_id
  INDEX idx_stock_part (part_id, created_at),
  CONSTRAINT fk_stock_movement_part FOREIGN KEY (part_id) REFERENCES part(part_id),
  CONSTRAINT fk_stock_movement_order FOREIGN KEY (order_id) REFERENCES `order`(order_id),
  CONSTRAINT fk_stock_movement_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT chk_stock_sale_order CHECK (reason <> 'sale' OR order_id IS NOT NULL),
  CONSTRAINT chk_stock_direction CHECK (
    (event_type = 'in' AND quantity_change > 0) OR (event_type = 'out' AND quantity_change < 0))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE payment (
  payment_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  amount INT NOT NULL,
  customer_payment_method_id INT NULL,                -- ถ้าลูกค้าเลือกตัดบัตรที่บันทึกไว้
  payment_method_type_id INT NOT NULL,
  status ENUM('pending', 'completed', 'failed') NOT NULL DEFAULT 'pending',
  paid_at TIMESTAMP NULL,                             -- มีค่าเมื่อ status = 'completed' เท่านั้น
  gateway_ref VARCHAR(100) NULL UNIQUE,               -- อ้างอิงฝั่ง payment gateway (กัน webhook ยิงซ้ำ)
  -- 1 order มี payment ที่ pending/completed ได้ใบเดียว (failed มีกี่ใบก็ได้ เพื่อให้ลองจ่ายใหม่)
  active_order_id BIGINT GENERATED ALWAYS AS (IF(status IN ('pending', 'completed'), order_id, NULL)) STORED,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,                       
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by_admin_id INT NULL,                      
  UNIQUE KEY uq_payment_one_active_per_order (active_order_id),
  INDEX idx_payment_order (order_id),
  CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES `order`(order_id),
  CONSTRAINT fk_payment_cpm FOREIGN KEY (customer_payment_method_id) REFERENCES customer_payment_method(customer_payment_method_id),
  CONSTRAINT fk_payment_method_type FOREIGN KEY (payment_method_type_id) REFERENCES payment_method_type(payment_method_type_id),
  CONSTRAINT fk_payment_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_payment_updated_by_admin_id FOREIGN KEY (updated_by_admin_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT chk_payment_amount CHECK (amount >= 0),
  CONSTRAINT chk_payment_paid_at CHECK (
    (status = 'completed' AND paid_at IS NOT NULL) OR (status <> 'completed' AND paid_at IS NULL))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 7. Image
-- ================================

CREATE TABLE image (
  image_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  reference_id INT NOT NULL,
  reference_type VARCHAR(20) NOT NULL,                
  file_name VARCHAR(100) NOT NULL,
  image_url VARCHAR(500) NOT NULL,
  alt_text VARCHAR(50) NOT NULL,
  is_primary BIT(1) NOT NULL DEFAULT 0,
  sort_image INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by_admin_id INT NULL,
  INDEX idx_image_reference (reference_type, reference_id),
  CONSTRAINT fk_image_created_by_admin_id FOREIGN KEY (created_by_admin_id) REFERENCES admin_user(admin_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;