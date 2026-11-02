package entity

// ปลายทางจัดส่ง ส่ง nil ได้ตอน preview ที่ยังไม่เลือกที่อยู่ (จะใช้ค่า default)
type ShippingDestination struct {
	ProvinceID *int32
}

type ShippingRate struct {
	Fee                     int32
	FreeShippingMinSubtotal *int32
}
