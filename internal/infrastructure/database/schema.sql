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
  phone_number VARCHAR(13) NOT NULL,
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
  amount INT NOT NULL,
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
  is_active BIT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  CONSTRAINT fk_part_part_brand FOREIGN KEY (part_brand_id) REFERENCES part_brand(part_brand_id),
  CONSTRAINT fk_part_part_type FOREIGN KEY (part_type_id) REFERENCES part_type(part_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ================================
-- 3. Customer-related tables
-- ================================

-- Table: customer_payment_method
CREATE TABLE customer_payment_method (
  customer_payment_method_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  account_name VARCHAR(60) NULL,
  payment_method_type ENUM('credit_card', 'bank_account', 'paypal') NOT NULL,
  card_token VARCHAR(255) NULL UNIQUE,
  is_default BIT(1) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  CONSTRAINT fk_cpm_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: address
CREATE TABLE address (
  address_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  title VARCHAR(50) NULL,
  address_line_1 VARCHAR(150) NOT NULL,
  address_line_2 VARCHAR(150) NULL,
  country VARCHAR(50) NOT NULL,
  city VARCHAR(50) NOT NULL,
  postal_code VARCHAR(8) NOT NULL,
  phone_number VARCHAR(13) NOT NULL,
  is_default BIT(1) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  CONSTRAINT fk_address_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: cart
CREATE TABLE cart (
  cart_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  amount INT NOT NULL DEFAULT 0,
  is_checkout BIT(1) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_cart_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id)
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
-- customer_id = เจ้าของออเดอร์ ไม่ว่าใครสร้าง
-- channel: web/mobile = ลูกค้าสั่งเอง, call/line/facebook/walk_in = แอดมินบันทึกให้
-- walk_in ที่ไม่ระบุตัวตน: customer_id และ address_id เป็น NULL ได้ (รับของหน้าร้าน)
-- หน้าร้าน: ให้แอปตั้ง status = 'delivered' และ payment.status = 'completed' ทันที
-- ผู้เปลี่ยนสถานะ (ระบบ/ลูกค้า/แอดมิน) เก็บใน order_status_history
CREATE TABLE `order` (
  order_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_number VARCHAR(13) NOT NULL UNIQUE,
  customer_id INT NULL,
  channel ENUM('web', 'mobile', 'call', 'line', 'facebook', 'walk_in') NOT NULL DEFAULT 'web',
  placed_by_admin_user_id INT NULL,               -- NULL = ลูกค้าสั่งเอง, มีค่า = แอดมินสั่งแทน
  customer_payment_method_id INT NULL,            -- NULL = เช่น COD / จ่ายหน้าร้าน
  cart_id INT NULL,                               -- ออเดอร์ที่แอดมินคีย์มักไม่มี cart
  amount INT NOT NULL DEFAULT 0,
  address_id INT NULL,
  discount_id INT NULL,                           -- NULL = ไม่ใช้ส่วนลด
  status ENUM('pending', 'processing', 'shipped', 'delivered', 'cancelled') NOT NULL DEFAULT 'pending',
  remark VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_order_customer (customer_id),
  CONSTRAINT fk_order_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT fk_order_placed_by_admin FOREIGN KEY (placed_by_admin_user_id) REFERENCES admin_user(admin_user_id),
  CONSTRAINT fk_order_payment_method FOREIGN KEY (customer_payment_method_id) REFERENCES customer_payment_method(customer_payment_method_id),
  CONSTRAINT fk_order_cart FOREIGN KEY (cart_id) REFERENCES cart(cart_id),
  CONSTRAINT fk_order_address FOREIGN KEY (address_id) REFERENCES address(address_id),
  CONSTRAINT fk_order_discount FOREIGN KEY (discount_id) REFERENCES discount(discount_id),
  -- ต้องมี customer_id และ address_id ทุกช่องทาง ยกเว้น walk_in
  CONSTRAINT chk_order_walkin CHECK (
    channel = 'walk_in' OR (customer_id IS NOT NULL AND address_id IS NOT NULL)
  ),
  -- web/mobile ต้องไม่มีแอดมิน, ช่องทางอื่นต้องมีแอดมินเป็นผู้บันทึก
  CONSTRAINT chk_order_channel_admin CHECK (
    (channel IN ('web', 'mobile') AND placed_by_admin_user_id IS NULL)
    OR (channel NOT IN ('web', 'mobile') AND placed_by_admin_user_id IS NOT NULL)
  )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: order_item
CREATE TABLE order_item (
  order_item_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  part_id INT NOT NULL,
  quantity INT NOT NULL,
  price INT NOT NULL,                             -- ราคา ณ ตอนสั่งซื้อ
  CONSTRAINT fk_order_item_order FOREIGN KEY (order_id) REFERENCES `order`(order_id) ON DELETE CASCADE,
  CONSTRAINT fk_order_item_part FOREIGN KEY (part_id) REFERENCES part(part_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: order_sequence
CREATE TABLE order_sequence (
  seq_date DATE NOT NULL PRIMARY KEY,
  `last_value` INT UNSIGNED NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


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
  part_brand_id INT NULL,
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
  CONSTRAINT fk_stock_movement_part_brand FOREIGN KEY (part_brand_id) REFERENCES part_brand(part_brand_id),
  CONSTRAINT fk_stock_movement_order FOREIGN KEY (order_id) REFERENCES `order`(order_id),
  CONSTRAINT chk_stock_sale_order CHECK (reason <> 'sale' OR order_id IS NOT NULL)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: payment
CREATE TABLE payment (
  payment_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  amount INT NOT NULL,
  payment_method ENUM('credit_card', 'bank_transfer', 'paypal', 'cod', 'cash') NOT NULL,
  status ENUM('pending', 'completed', 'failed') NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  updated_by INT NULL,
  CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES `order`(order_id)
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