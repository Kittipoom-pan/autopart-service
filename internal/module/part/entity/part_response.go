package entity

type PartRes struct {
	PartID      int32   `json:"part_id" db:"part_id"`
	CarBrandID  int32   `json:"car_brand_id" db:"car_brand_id"`
	PartBrandID int32   `json:"part_brand_id" db:"part_brand_id"`
	PartTypeID  int32   `json:"part_type_id" db:"part_type_id"`
	Name        string  `json:"name" db:"name"`
	SKU         string  `json:"sku" db:"sku"`
	Description *string `json:"description,omitempty" db:"description"`
	Price       int32   `json:"price" db:"price"`
	Quantity    int32   `json:"quantity" db:"quantity"`
	WeightG     *int32  `json:"weight_g" db:"weight_g"`
}
