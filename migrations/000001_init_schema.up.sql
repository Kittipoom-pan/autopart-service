-- ================================
-- 1. Master / independent tables
-- ================================

-- Table: customer
CREATE TABLE customer (
  customer_id INT AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(40) NOT NULL,
  last_name VARCHAR(40) NOT NULL,
  username VARCHAR(50) NOT NULL,
  email VARCHAR(50) NOT NULL,
  password VARCHAR(150) NOT NULL,
  birth_date DATE NULL,
  phone_number VARCHAR(15) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  is_active BIT(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE UNIQUE INDEX unique_active_phone ON customer (phone_number, is_active);
CREATE UNIQUE INDEX unique_active_username ON customer (username, is_active);
CREATE UNIQUE INDEX unique_active_email ON customer (email, is_active);

-- Table: admin_user
CREATE TABLE admin_user (
  admin_user_id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(40) NOT NULL,
  email VARCHAR(60) NULL,
  password VARCHAR(150) NOT NULL,
  role ENUM('super_admin', 'staff') NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  is_active BIT(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE UNIQUE INDEX unique_active_admin_username ON admin_user (username, is_active);
CREATE UNIQUE INDEX unique_active_admin_email ON admin_user (email, is_active);

-- Table: discount
CREATE TABLE discount (
  discount_id INT AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(60) NOT NULL UNIQUE,
  discount_type ENUM('percentage', 'fixed') NOT NULL,
  discount_value INT NOT NULL,
  description VARCHAR(255) NULL,
  min_order_amount INT NOT NULL DEFAULT 0,        
  expiration_date DATE NULL,                      -- NULL = ไม่หมดอายุ
  is_active BIT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: part_type
CREATE TABLE part_type (
  part_type_id INT AUTO_INCREMENT PRIMARY KEY,
  parent_id INT NULL,                            
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255) NULL,
  CONSTRAINT fk_part_type_parent FOREIGN KEY (parent_id) REFERENCES part_type(part_type_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: part_brand
CREATE TABLE part_brand (
  part_brand_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: car_brand
CREATE TABLE car_brand (
  car_brand_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 2. Dependent tables
-- ================================

-- Table: car_model
CREATE TABLE car_model (
  car_model_id INT AUTO_INCREMENT PRIMARY KEY,
  car_brand_id INT NOT NULL,
  name VARCHAR(50) NOT NULL,
  year_from YEAR DEFAULT NULL,
  year_to YEAR DEFAULT NULL,
  CONSTRAINT fk_car_model_car_brand FOREIGN KEY (car_brand_id) REFERENCES car_brand(car_brand_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: part
CREATE TABLE part (
  part_id INT AUTO_INCREMENT PRIMARY KEY,
  part_brand_id INT NOT NULL,
  part_type_id INT NOT NULL,
  name VARCHAR(100) NOT NULL UNIQUE,
  sku VARCHAR(50) NOT NULL UNIQUE,
  description TEXT NULL,
  price INT NOT NULL DEFAULT 0,
  quantity INT NOT NULL DEFAULT 0,                -- ต้องตรงกับผลรวมของ stock_movement
  weight_g INT NULL,
  is_active BIT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  CONSTRAINT fk_part_part_brand FOREIGN KEY (part_brand_id) REFERENCES part_brand(part_brand_id),
  CONSTRAINT fk_part_part_type FOREIGN KEY (part_type_id) REFERENCES part_type(part_type_id),
  CONSTRAINT chk_part_values  CHECK (quantity >= 0 AND price >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 3. Customer-related tables
-- ================================

CREATE TABLE payment_method_type (
    payment_method_type_id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE, -- 'credit_card', 'cod', 'bank_transfer', 'cash'
    name VARCHAR(40) NOT NULL,       -- ชื่อที่แสดงผล เช่น 'Credit Card', 'Cash on Delivery'
    is_active BIT(1) NOT NULL DEFAULT 1
);

-- Table: customer_payment_method
CREATE TABLE customer_payment_method (
  customer_payment_method_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  account_name VARCHAR(60) NULL,
  payment_method_type_id INT NOT NULL,
  card_token VARCHAR(255) NULL UNIQUE,
  is_default BIT(1) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  CONSTRAINT fk_cpm_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT fk_cpm_payment_type FOREIGN KEY (payment_method_type_id) REFERENCES payment_method_type(payment_method_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `province` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` int(11) NOT NULL,
  `name_in_thai` varchar(150) NOT NULL,
  `name_in_english` varchar(150) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ux_province_code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=78 DEFAULT CHARSET=utf8;

CREATE TABLE `district` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` int(11) NOT NULL,
  `name_in_thai` varchar(150) NOT NULL,
  `name_in_english` varchar(150) NOT NULL,
  `province_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ux_district_code` (`code`) USING BTREE,
  KEY `ix_district_province_id` (`province_id`) USING BTREE,
  CONSTRAINT `fk_district_province` FOREIGN KEY (`province_id`) REFERENCES `province` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=929 DEFAULT CHARSET=utf8;

CREATE TABLE `subdistrict` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` int(11) NOT NULL,
  `name_in_thai` varchar(150) NOT NULL,
  `name_in_english` varchar(150) NULL,
  `district_id` int(11) NOT NULL,
  zip_code CHAR(5) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ux_subdistrict_id` (`id`) USING BTREE,
  KEY `ix_subdistrict_district_id` (`district_id`) USING BTREE,
  CONSTRAINT `fk_subdistrict_district` FOREIGN KEY (`district_id`) REFERENCES `district` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Table: address
CREATE TABLE address (
  address_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  title VARCHAR(50) NULL, -- ชื่อเรียกที่อยู่ เช่น บ้าน, ที่ทำงาน
  recipient_name VARCHAR(100) NOT NULL, -- recipient_name ชื่อผู้รับพัสดุ
  address_line VARCHAR(255) NOT NULL,
  subdistrict_id INT NOT NULL,   
  phone_number VARCHAR(15) NOT NULL,
  is_default BIT(1) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  INDEX idx_address_subdistrict (subdistrict_id),
  CONSTRAINT fk_address_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT fk_address_subdistrict FOREIGN KEY (subdistrict_id) REFERENCES subdistrict(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: cart
CREATE TABLE cart (
  cart_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  subtotal INT NOT NULL DEFAULT 0,                -- ราคารวมสินค้าทั้งหมด (ก่อนหักส่วนลดและยังไม่รวมค่าส่ง)
  is_checkout BIT(1) NOT NULL DEFAULT 0,
  -- จำกัดให้ลูกค้ามี Active Cart ได้เพียง 1 ใบ
  active_customer_id INT
    GENERATED ALWAYS AS (
        IF(is_checkout = 0, customer_id, NULL)
    ) STORED,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_cart_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT uq_cart_one_active UNIQUE (active_customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: cart_item
CREATE TABLE cart_item (
  cart_item_id INT AUTO_INCREMENT PRIMARY KEY,
  cart_id INT NOT NULL,
  part_id INT NOT NULL,
  quantity INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_cart_item_cart FOREIGN KEY (cart_id) REFERENCES cart(cart_id) ON DELETE CASCADE,
  CONSTRAINT fk_cart_item_part FOREIGN KEY (part_id) REFERENCES part(part_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 4. Token tables
-- ================================

-- Table: admin_token
CREATE TABLE admin_token (
  id INT AUTO_INCREMENT PRIMARY KEY,
  admin_user_id INT NOT NULL,
  token VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_revoked BIT(1) NOT NULL DEFAULT 0,
  CONSTRAINT fk_token_admin_user FOREIGN KEY (admin_user_id) REFERENCES admin_user(admin_user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: customer_token
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
-- 5. Compatible car table
-- ================================

-- Table: compatible_car
CREATE TABLE compatible_car (
  compatible_id INT AUTO_INCREMENT PRIMARY KEY,
  part_id INT NOT NULL,
  car_model_id INT NOT NULL,
  year_from YEAR DEFAULT NULL,                    -- NULL = ไม่จำกัดปี
  year_to YEAR DEFAULT NULL,
  CONSTRAINT fk_compatible_car_part FOREIGN KEY (part_id) REFERENCES part(part_id) ON DELETE CASCADE,
  CONSTRAINT fk_compatible_car_model FOREIGN KEY (car_model_id) REFERENCES car_model(car_model_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 6. Order-related tables
-- ================================

-- Table: `order`
-- channel: web/mobile = ลูกค้าสั่งเอง, call/line/facebook/walk_in = แอดมินบันทึกให้
-- หน้าร้าน: ให้แอปตั้ง status = 'delivered' และ payment.status = 'completed' ทันที
-- ผู้เปลี่ยนสถานะ (ระบบ/ลูกค้า/แอดมิน) เก็บใน order_status_history
CREATE TABLE `order` (
  order_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_number VARCHAR(13) NOT NULL UNIQUE,
  customer_id INT NULL,
  channel ENUM('web', 'mobile', 'call', 'line', 'facebook', 'walk_in') NOT NULL DEFAULT 'web',
  placed_by_admin_user_id INT NULL,               -- NULL = ลูกค้าสั่งเอง, มีค่า = แอดมินสั่งแทน
  cart_id INT NULL,                               -- ออเดอร์ที่แอดมินคีย์มักไม่มี cart
  subtotal INT NOT NULL DEFAULT 0,                -- ราคารวมสินค้าทั้งหมด (ก่อนหักส่วนลดและยังไม่รวมค่าส่ง)
  discount_amount INT NOT NULL DEFAULT 0,         -- จำนวนเงินส่วนลดที่ใช้จริง (บาท)
  shipping_fee INT NOT NULL DEFAULT 0,            -- ค่าจัดส่ง
  total_amount INT NOT NULL DEFAULT 0,            -- ยอดสุทธิที่ต้องจ่ายจริง (Subtotal - Discount + Shipping)
  discount_id INT NULL,                           -- NULL = ไม่ใช้ส่วนลด
  status ENUM('pending', 'processing', 'shipped', 'delivered', 'cancelled', 'completed') NOT NULL DEFAULT 'pending',
  remark VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_order_customer (customer_id),
  CONSTRAINT fk_order_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT fk_order_placed_by_admin FOREIGN KEY (placed_by_admin_user_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_order_cart FOREIGN KEY (cart_id) REFERENCES cart(cart_id),
  CONSTRAINT fk_order_discount FOREIGN KEY (discount_id) REFERENCES discount(discount_id),
  UNIQUE KEY uq_order_cart (cart_id),
  CONSTRAINT chk_order_amounts CHECK (
    subtotal >= 0 AND discount_amount >= 0 AND shipping_fee >= 0
    AND total_amount = subtotal - discount_amount + shipping_fee),
  -- ต้องมี customer_id ทุกช่องทาง ยกเว้น walk_in
  CONSTRAINT chk_order_walkin CHECK (
    channel = 'walk_in' OR customer_id IS NOT NULL
  ),
  -- web/mobile ต้องไม่มีแอดมิน, ช่องทางอื่นต้องมีแอดมินเป็นผู้บันทึก
  CONSTRAINT chk_order_channel_admin CHECK (
    (channel IN ('web', 'mobile') AND placed_by_admin_user_id IS NULL)
    OR (channel NOT IN ('web', 'mobile') AND placed_by_admin_user_id IS NOT NULL)
  ),
  INDEX idx_order_status_created (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: order_item
CREATE TABLE order_item (
  order_item_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  part_id INT NOT NULL,
  quantity INT NOT NULL,
  price INT NOT NULL,                             -- ราคา ณ ตอนสั่งซื้อ
  CONSTRAINT fk_order_item_order FOREIGN KEY (order_id) REFERENCES `order`(order_id) ON DELETE CASCADE,
  CONSTRAINT fk_order_item_part FOREIGN KEY (part_id) REFERENCES part(part_id),
  UNIQUE KEY uq_order_item_part (order_id, part_id),
  CONSTRAINT chk_order_item_values CHECK (quantity > 0 AND price >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: order_sequence
CREATE TABLE order_sequence (
  seq_date DATE NOT NULL PRIMARY KEY,
  `last_value` INT UNSIGNED NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `order_shipping_address` (
  `order_id` BIGINT PRIMARY KEY,
  `recipient_name` VARCHAR(100) NOT NULL,
  `phone_number` VARCHAR(15) NOT NULL,
  `address_line` VARCHAR(255) NOT NULL,
  `subdistrict` VARCHAR(100) NOT NULL,
  `district` VARCHAR(100) NOT NULL,
  `province` VARCHAR(100) NOT NULL,
  `zip_code` CHAR(5) NOT NULL,
  CONSTRAINT fk_order_shipping FOREIGN KEY (`order_id`) REFERENCES `order`(`order_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE shipping_rate (
    shipping_rate_id           INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name                       VARCHAR(100) NOT NULL,
    province_id                INT          NULL COMMENT 'จังหวัดที่ใช้ค่าส่งนี้; NULL = ใช้ได้ทุกจังหวัด',
    min_weight_g               INT          NOT NULL DEFAULT 0 COMMENT 'รวมค่านี้ (>=)',
    max_weight_g               INT          NULL COMMENT 'ไม่รวมค่านี้ (<); NULL = ไม่จำกัด ไว้ใช้กำหนดค่าส่งตามช่วงน้ำหนัก เช่น 0–1,000 และ 1,000–5,000 กรัม',
    fee                        INT          NOT NULL COMMENT 'ค่าจัดส่ง',
    free_shipping_min_subtotal INT          NULL COMMENT 'ยอดสินค้าหลังหักส่วนลดตั้งแต่เท่านี้ส่งฟรี NULL = ไม่มี',
    is_active                  TINYINT(1)   NOT NULL DEFAULT 1,
    created_at                 TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP     NULL,
    CONSTRAINT fk_shipping_rate_province  FOREIGN KEY (province_id) REFERENCES province(id),
    CONSTRAINT chk_shipping_fee    CHECK (fee >= 0),
    CONSTRAINT chk_shipping_weight CHECK (max_weight_g IS NULL OR max_weight_g > min_weight_g),
    KEY idx_shipping_lookup (is_active, province_id, min_weight_g)
) ENGINE=InnoDB;

-- Table: stock_movement
-- ขาย (ออนไลน์/หน้าร้าน) ต้องผ่าน order เสมอ: reason = 'sale' และมี order_id
-- ยกเลิกออเดอร์/คืนของ: reason = 'return' ผูก order_id เดิม
-- รับของเข้า/นับสต็อก/ของเสีย: ไม่ต้องมี order_id ให้ใส่ remark และ created_by (admin_user_id)
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
  created_by INT NULL,
  INDEX idx_stock_part (part_id, created_at),
  CONSTRAINT fk_stock_movement_part FOREIGN KEY (part_id) REFERENCES part(part_id),
  CONSTRAINT fk_stock_movement_order FOREIGN KEY (order_id) REFERENCES `order`(order_id),
  CONSTRAINT chk_stock_sale_order CHECK (reason <> 'sale' OR order_id IS NOT NULL),
  CONSTRAINT chk_stock_direction CHECK ((event_type = 'in'  AND quantity_change > 0) OR (event_type = 'out' AND quantity_change < 0))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: payment
CREATE TABLE payment (
  payment_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  amount INT NOT NULL,
  customer_payment_method_id INT NULL, -- ถ้าลูกค้าเลือกตัดบัตรที่บันทึกไว้ (เชื่อมกับตารางเก็บบัตร)
  payment_method_type_id INT NOT NULL, -- ประเภทการจ่ายเงิน เช่น บัตรเครดิต, COD, โอนเงิน
  status ENUM('pending', 'completed', 'failed') NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES `order`(order_id),
  CONSTRAINT fk_payment_cpm FOREIGN KEY (customer_payment_method_id) REFERENCES customer_payment_method(customer_payment_method_id),
  CONSTRAINT fk_payment_method_type FOREIGN KEY (payment_method_type_id) REFERENCES payment_method_type(payment_method_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 7. Image table
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
  created_by INT NULL,
  INDEX idx_image_reference (reference_type, reference_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;