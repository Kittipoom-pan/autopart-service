package entitie

type PartReq struct {
	PartID      uint32  `json:"part_id" db:"part_id"`
	CarBrandID  uint32  `json:"car_brand_id" db:"car_brand_id"`
	PartBrandID uint32  `json:"part_brand_id" db:"part_brand_id"`
	PartTypeID  uint32  `json:"part_type_id" db:"part_type_id"`
	Name        string  `json:"name" db:"name"`
	SKU         string  `json:"sku" db:"sku"`
	Description *string `json:"description,omitempty" db:"description"`
	Price       *int    `json:"price,omitempty" db:"price"`
	IsActive    bool    `json:"is_active" db:"is_active"`
	Quantity    *int    `json:"quantity,omitempty" db:"quantity"`
}
