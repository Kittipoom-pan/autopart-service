-- Table: customer
CREATE TABLE customer (
  customer_id INT AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(40),
  last_name VARCHAR(40),
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(150),
  birth_date DATE,
  phone_number VARCHAR(13) UNIQUE,
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  updated_at TIMESTAMP NULL DEFAULT NULL,
  updated_by VARCHAR(40),
  is_active TINYINT(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: admin_user
CREATE TABLE admin_user (
  admin_user_id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(40) NOT NULL UNIQUE,
  email VARCHAR(60) NOT NULL UNIQUE,
  password VARCHAR(150) NOT NULL,
  role ENUM('super_admin', 'staff') NOT NULL,
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  updated_at TIMESTAMP NULL DEFAULT NULL,
  updated_by VARCHAR(40),
  is_active TINYINT(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: discount
CREATE TABLE discount (
  discount_id INT AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(60) NOT NULL UNIQUE,
  discount_type ENUM('percentage', 'fixed'),
  amount INT,
  description VARCHAR(255),
  min_order INT,
  expiration_date DATE,
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  updated_at TIMESTAMP NULL DEFAULT NULL,
  updated_by VARCHAR(40)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: part_type
CREATE TABLE part_type (
  part_type_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: part_brand
CREATE TABLE part_brand (
  part_brand_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: car_brand
CREATE TABLE car_brand (
  car_brand_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: part
CREATE TABLE part (
  part_id INT AUTO_INCREMENT PRIMARY KEY,
  car_brand_id INT NOT NULL,
  part_brand_id INT NOT NULL,
  part_type_id INT NOT NULL,
  name VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  price INT,
  quantity INT,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  updated_at TIMESTAMP NULL DEFAULT NULL,
  updated_by VARCHAR(40),
  CONSTRAINT fk_part_car_brand FOREIGN KEY (car_brand_id) REFERENCES car_brand(car_brand_id),
  CONSTRAINT fk_part_part_brand FOREIGN KEY (part_brand_id) REFERENCES part_brand(part_brand_id),
  CONSTRAINT fk_part_part_type FOREIGN KEY (part_type_id) REFERENCES part_type(part_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: customer_payment_method
CREATE TABLE customer_payment_method (
  customer_payment_method_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  account_name VARCHAR(60),
  payment_method_type ENUM('credit_card', 'bank_account', 'paypal') NOT NULL,
  card_token VARCHAR(255) UNIQUE,
  is_default TINYINT(1) DEFAULT 1,
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  updated_at TIMESTAMP NULL DEFAULT NULL,
  updated_by VARCHAR(40),
  CONSTRAINT fk_cpm_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: token
CREATE TABLE token (
  id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL UNIQUE,
  token VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NULL DEFAULT NULL,
  created_at TIMESTAMP NULL DEFAULT NULL,
  CONSTRAINT fk_token_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: address
CREATE TABLE address (
  address_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  title VARCHAR(50),
  address_line_1 VARCHAR(150),
  address_line_2 VARCHAR(150),
  country VARCHAR(50),
  city VARCHAR(50),
  postal_code VARCHAR(8),
  phone_number VARCHAR(13),
  is_default BIT(1),
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  updated_at TIMESTAMP NULL DEFAULT NULL,
  updated_by VARCHAR(40),
  CONSTRAINT fk_address_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: cart
CREATE TABLE cart (
  cart_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  amount INT,
  is_checkout BIT(1),
  created_at TIMESTAMP NULL DEFAULT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL,
  CONSTRAINT fk_cart_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: cart_item
CREATE TABLE cart_item (
  cart_item_id INT AUTO_INCREMENT PRIMARY KEY,
  cart_id INT NOT NULL,
  part_id INT NOT NULL,
  quantity INT,
  created_at TIMESTAMP NULL DEFAULT NULL,
  updated_at TIMESTAMP NULL DEFAULT NULL,
  CONSTRAINT fk_cart_item_cart FOREIGN KEY (cart_id) REFERENCES cart(cart_id),
  CONSTRAINT fk_cart_item_part FOREIGN KEY (part_id) REFERENCES part(part_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: order
CREATE TABLE `order` (
  order_id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT NOT NULL,
  customer_payment_method_id INT,
  cart_id INT,
  amount INT,
  address_id INT,
  discount_id INT,
  status ENUM('pending', 'completed', 'cancelled'),
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  updated_at TIMESTAMP NULL DEFAULT NULL,
  updated_by VARCHAR(40),
  CONSTRAINT fk_order_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
  CONSTRAINT fk_order_payment_method FOREIGN KEY (customer_payment_method_id) REFERENCES customer_payment_method(customer_payment_method_id),
  CONSTRAINT fk_order_cart FOREIGN KEY (cart_id) REFERENCES cart(cart_id),
  CONSTRAINT fk_order_address FOREIGN KEY (address_id) REFERENCES address(address_id),
  CONSTRAINT fk_order_discount FOREIGN KEY (discount_id) REFERENCES discount(discount_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: order_item
CREATE TABLE order_item (
  order_item_id INT AUTO_INCREMENT PRIMARY KEY,
  order_id INT NOT NULL,
  part_id INT,
  quantity INT,
  price INT,
  CONSTRAINT fk_order_item_order FOREIGN KEY (order_id) REFERENCES `order`(order_id),
  CONSTRAINT fk_order_item_part FOREIGN KEY (part_id) REFERENCES part(part_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: payment
CREATE TABLE payment (
  payment_id INT AUTO_INCREMENT PRIMARY KEY,
  order_id INT NOT NULL,
  amount INT,
  payment_method ENUM('credit_card', 'bank_transfer', 'paypal', 'cod'),
  status ENUM('pending', 'completed', 'failed'),
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  updated_at TIMESTAMP NULL DEFAULT NULL,
  updated_by VARCHAR(40),
  CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES `order`(order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table: stock_movement
CREATE TABLE stock_movement (
  stock_movement_id INT AUTO_INCREMENT PRIMARY KEY,
  part_id INT NOT NULL,
  part_brand_id INT,
  quantity_change INT,
  price INT,
  event_type ENUM('in', 'out'),
  remark VARCHAR(255),
  created_at TIMESTAMP NULL DEFAULT NULL,
  created_by VARCHAR(40),
  CONSTRAINT fk_stock_movement_part FOREIGN KEY (part_id) REFERENCES part(part_id),
  CONSTRAINT fk_stock_movement_part_brand FOREIGN KEY (part_brand_id) REFERENCES part_brand(part_brand_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
