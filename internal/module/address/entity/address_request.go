package entity

type CreateAddressReq struct {
	Title         *string `json:"title"`
	RecipientName string  `json:"recipient_name"`
	AddressLine   string  `json:"address_line"`
	SubdistrictID int32   `json:"subdistrict_id"`
	PhoneNumber   string  `json:"phone_number"`
	IsDefault     bool    `json:"is_default"`
}

type UpdateAddressReq struct {
	Title         *string `json:"title"`
	RecipientName string  `json:"recipient_name"`
	AddressLine   string  `json:"address_line"`
	SubdistrictID int32   `json:"subdistrict_id"`
	PhoneNumber   string  `json:"phone_number"`
}
