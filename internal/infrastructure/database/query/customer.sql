-- name: GetCustomer :one
SELECT customer_id, first_name, last_name, username, email, birth_date, phone_number, is_active, created_at FROM customer WHERE customer_id = ? and is_active = 1;

-- name: GetCustomerByUsername :one
SELECT customer_id, username, password FROM customer WHERE username = ? and is_active = 1;

-- name: ListCustomers :many
SELECT customer_id, first_name, last_name, username, email, birth_date, phone_number, is_active, created_at FROM customer WHERE is_active = 1 ORDER BY created_at;

-- name: CreateCustomer :execresult
INSERT INTO customer (
    first_name,
    last_name,
    username,
    email,
    password,
    birth_date,
    phone_number,
    created_at,
    created_by_admin_id
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateCustomer :execresult
UPDATE customer SET
    first_name = ?,
    last_name = ?,
    username = ?,
    email = ?,
    birth_date = ?,
    phone_number = ?,
    updated_at = ?,
    updated_by_admin_id = ?
WHERE customer_id = ? and is_active = 1;

-- name: UpdateCustomerIsActive :execresult
UPDATE customer SET
    is_active = ?,
    updated_at = ?,
    updated_by_admin_id = ?
WHERE customer_id = ? and is_active = 1;

-- name: UpdateCustomerPassword :execresult
UPDATE customer SET
    password = ?,
    updated_at = ?,
    updated_by_admin_id = ?
WHERE customer_id = ? and is_active = 1;

-- name: AdminListCustomers :many
-- Pagination (Limit/Offset) 
SELECT customer_id, first_name, last_name, username, is_active, created_at 
FROM customer 
ORDER BY created_at DESC 
LIMIT ? OFFSET ?;

-- name: AdminUpdateCustomerStatus :execresult
UPDATE customer SET
    is_active = ?,
    updated_at = ?,
    updated_by_admin_id = ?
WHERE customer_id = ?;

-- ============================================================
-- Below are queries for delete customer data. 

-- name: LockCustomerForDelete :one
SELECT customer_id FROM customer 
WHERE customer_id = ? AND is_active = 1 
FOR UPDATE;

-- name: DeleteCustomer :execresult
UPDATE customer
SET first_name = NULL,
    last_name = NULL,
    username = NULL,
    email = NULL,
    password = NULL,
    phone_number = NULL,
    birth_date = NULL,
    is_active = 0,
    deleted_at = CURRENT_TIMESTAMP,
    updated_by_admin_id = ?
WHERE customer_id = ? AND is_active = 1;

-- name: DeleteCustomerAddresses :exec
DELETE FROM address WHERE customer_id = ?;

-- name: DeleteCustomerPaymentMethods :exec
UPDATE customer_payment_method
SET account_name = NULL,
    card_token = NULL,
    is_default = 0
WHERE customer_id = ?;

-- name: DeleteCustomerTokens :exec
DELETE FROM customer_token WHERE customer_id = ?;
-- ============================================================

-- ============================================================
-- สิ่งที่ SQL ข้างบนทำแทนไม่ได้ (ต้องจัดการนอก DB)
--   * card_token ที่ฝั่ง payment gateway: ต้องเรียก API ของ gateway ให้ลบบัตรด้วย
--   * order_shipping_address (ชื่อผู้รับ/เบอร์/ที่อยู่ ณ ตอนสั่ง) ยังค้างในออเดอร์เก่า
--     ต้องตัดสินใจว่าเก็บไว้นานเท่าไรตามข้อกำหนดบัญชี/ภาษี แล้วค่อยล้างหลังพ้นกำหนด
--     (ระยะเวลาที่ต้องเก็บ ให้ผู้รู้ด้านกฎหมาย/บัญชียืนยัน)
--   * ตะกร้าที่ยัง active ไม่มีข้อมูลส่วนตัว ปล่อยไว้ได้
-- ============================================================
