SET FOREIGN_KEY_CHECKS = 0;


-- 1. Image
DROP TABLE IF EXISTS image;

-- 2. Payment / Order Items
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS order_item;

-- 3. Order related
DROP TABLE IF EXISTS stock_movement;
DROP TABLE IF EXISTS order_shipping_address;
DROP TABLE IF EXISTS `order`;
DROP TABLE IF EXISTS order_sequence;

-- 4. Compatible / Cart Items
DROP TABLE IF EXISTS compatible_car;
DROP TABLE IF EXISTS cart_item;
DROP TABLE IF EXISTS cart;

-- 5. Address hierarchy
DROP TABLE IF EXISTS address;
DROP TABLE IF EXISTS subdistrict;
DROP TABLE IF EXISTS district;
DROP TABLE IF EXISTS province;

-- 6. Shipping
DROP TABLE IF EXISTS shipping_rate;

-- 7. Payment method
DROP TABLE IF EXISTS customer_payment_method;
DROP TABLE IF EXISTS payment_method_type;

-- 8. Token
DROP TABLE IF EXISTS customer_token;
DROP TABLE IF EXISTS admin_token;

-- 9. Part
DROP TABLE IF EXISTS part;

-- 10. Car model
DROP TABLE IF EXISTS car_model;

-- 11. Master tables
DROP TABLE IF EXISTS car_brand;
DROP TABLE IF EXISTS part_brand;
DROP TABLE IF EXISTS part_type;
DROP TABLE IF EXISTS discount;
DROP TABLE IF EXISTS admin_user;
DROP TABLE IF EXISTS customer;