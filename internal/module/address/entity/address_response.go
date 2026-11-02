package entity

import "time"

type AddressRes struct {
	AddressID     int32      `json:"address_id"`
	CustomerID    int32      `json:"customer_id"`
	Title         *string    `json:"title,omitempty"`
	RecipientName string     `json:"recipient_name"`
	AddressLine   string     `json:"address_line"`
	SubdistrictID int32      `json:"subdistrict_id"`
	Subdistrict   string     `json:"subdistrict,omitempty"`
	District      string     `json:"district,omitempty"`
	ProvinceID    *int32     `json:"province_id,omitempty"`
	Province      string     `json:"province,omitempty"`
	ZipCode       string     `json:"zip_code"`
	PhoneNumber   string     `json:"phone_number"`
	IsDefault     bool       `json:"is_default"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
