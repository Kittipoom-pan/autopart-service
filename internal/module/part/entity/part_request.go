package entity

type PartReq struct {
	PartID      int32   `json:"part_id" db:"part_id"`
	CarBrandID  int32   `json:"car_brand_id" db:"car_brand_id"`
	PartBrandID int32   `json:"part_brand_id" db:"part_brand_id"`
	PartTypeID  int32   `json:"part_type_id" db:"part_type_id"`
	Name        string  `json:"name" db:"name"`
	WeightG     *int32  `json:"weight_g,omitempty" db:"weight_g"`
	SKU         string  `json:"sku" db:"sku"`
	OemNumber   *string `json:"oem_number,omitempty" db:"oem_number"`
	Description *string `json:"description,omitempty" db:"description"`
	Price       int     `json:"price" db:"price"`
	IsActive    bool    `json:"is_active" db:"is_active"`
	Quantity    int     `json:"quantity" db:"quantity"`
}
