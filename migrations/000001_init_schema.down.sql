-- 1️⃣ Image table
DROP TABLE IF EXISTS image;

-- 2️⃣ Payment and order_item
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS order_item;

-- 3️⃣ Order table
DROP TABLE IF EXISTS stock_movement;
DROP TABLE IF EXISTS order_shipping_address;
DROP TABLE IF EXISTS `order`;
DROP TABLE IF EXISTS order_sequence;

-- 4️⃣ Compatible / stock / cart_item
DROP TABLE IF EXISTS compatible_car;
DROP TABLE IF EXISTS cart_item;

-- 5️⃣ Cart and related
DROP TABLE IF EXISTS cart;

-- 6️⃣ Address and payment method
DROP TABLE IF EXISTS address;
DROP TABLE IF EXISTS subdistrict;
DROP TABLE IF EXISTS district;
DROP TABLE IF EXISTS province;
DROP TABLE IF EXISTS customer_payment_method;
DROP TABLE IF EXISTS payment_method_type;

-- 7️⃣ Token tables
DROP TABLE IF EXISTS customer_token;
DROP TABLE IF EXISTS admin_token;

-- 8️⃣ Part table
DROP TABLE IF EXISTS part;
DROP TABLE IF EXISTS shipping_rate;

-- 9️⃣ Car model
DROP TABLE IF EXISTS car_model;

-- 10️⃣ Master / independent tables
DROP TABLE IF EXISTS car_brand;
DROP TABLE IF EXISTS part_brand;
DROP TABLE IF EXISTS part_type;
DROP TABLE IF EXISTS discount;
DROP TABLE IF EXISTS admin_user;
DROP TABLE IF EXISTS customer;
