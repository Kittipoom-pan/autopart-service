-- name: ListAddressesByCustomer :many
SELECT 
    a.address_id,
    a.customer_id,
    a.title,
    a.recipient_name,
    a.address_line,
    a.subdistrict_id,
    a.zip_code,
    a.phone_number,
    a.is_default,
    s.name_in_thai AS subdistrict,
    d.name_in_thai AS district,
    p.name_in_thai AS province
FROM address a
JOIN subdistrict s ON s.id = a.subdistrict_id
JOIN district d    ON d.id = s.district_id
JOIN province p    ON p.id = d.province_id
WHERE a.customer_id = ?;