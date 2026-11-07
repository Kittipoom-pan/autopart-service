-- name: GetAddressByIDAndCustomer :one
SELECT
    a.address_id,
    a.customer_id,
    a.title,
    a.recipient_name,
    a.address_line,
    a.subdistrict_id,
    s.zip_code,
    a.phone_number,
    a.is_default,
    a.created_at,
    a.created_by_admin_id,
    a.updated_at,
    a.updated_by_admin_id,
    s.name_in_thai AS subdistrict,
    d.name_in_thai AS district,
    d.province_id AS province_id,
    p.name_in_thai AS province
FROM address a
JOIN subdistrict s ON s.id = a.subdistrict_id
JOIN district d ON d.id = s.district_id
JOIN province p ON p.id = d.province_id
WHERE a.address_id = ? AND a.customer_id = ?;

-- name: GetSubdistrictDetailByID :one
SELECT
    s.id AS subdistrict_id,
    s.name_in_thai AS subdistrict,
    d.name_in_thai AS district,
    d.province_id AS province_id,
    p.name_in_thai AS province,
    s.zip_code
FROM subdistrict s
JOIN district d ON d.id = s.district_id
JOIN province p ON p.id = d.province_id
WHERE s.id = ?;

-- name: ListAddressesByCustomerID :many
SELECT
    a.address_id,
    a.customer_id,
    a.title,
    a.recipient_name,
    a.address_line,
    a.subdistrict_id,
    s.zip_code,
    a.phone_number,
    a.is_default,
    a.created_at,
    a.created_by_admin_id,
    a.updated_at,
    a.updated_by_admin_id
FROM address a
JOIN subdistrict s ON s.id = a.subdistrict_id
WHERE a.customer_id = ?
ORDER BY a.is_default DESC, a.created_at DESC;

-- name: GetDefaultAddressByCustomerID :one
SELECT
    a.address_id,
    a.customer_id,
    a.title,
    a.recipient_name,
    a.address_line,
    a.subdistrict_id,
    s.zip_code,
    a.phone_number,
    a.is_default,
    a.created_at,
    a.created_by_admin_id,
    a.updated_at,
    a.updated_by_admin_id
FROM address a
JOIN subdistrict s ON s.id = a.subdistrict_id
WHERE a.customer_id = ? AND a.is_default = 1
LIMIT 1;

-- name: CreateAddress :execresult
INSERT INTO address (
    customer_id,
    title,
    recipient_name,
    address_line,
    subdistrict_id,
    phone_number,
    is_default,
    created_at,
    created_by_admin_id
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateAddress :execresult
UPDATE address SET
    title = ?,
    recipient_name = ?,
    address_line = ?,
    subdistrict_id = ?,
    phone_number = ?,
    is_default = ?,
    updated_at = ?,
    updated_by_admin_id = ?
WHERE address_id = ? AND customer_id = ?;

-- name: ClearDefaultAddress :execresult
UPDATE address SET
    is_default = 0,
    updated_at = ?,
    updated_by_admin_id = ?
WHERE customer_id = ?;

-- name: SetDefaultAddress :execresult
UPDATE address SET
    is_default = 1,
    updated_at = ?,
    updated_by_admin_id = ?
WHERE address_id = ? AND customer_id = ?;

-- name: DeleteAddress :execresult
DELETE FROM address
WHERE address_id = ? AND customer_id = ?;
