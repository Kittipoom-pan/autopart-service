package entity

import "time"

type Address struct {
	AddressID     int32
	CustomerID    int32
	Title         *string
	RecipientName string
	AddressLine   string
	SubdistrictID int32
	ZipCode       string
	PhoneNumber   string
	IsDefault     bool
	CreatedAt     time.Time
	CreatedBy     *int32
	UpdatedAt     *time.Time
	UpdatedBy     *int32
}

type SubdistrictDetail struct {
	SubdistrictID int32
	Subdistrict   string
	District      string
	ProvinceID    int32
	Province      string
	ZipCode       string
}
