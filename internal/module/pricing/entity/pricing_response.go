package entity

type PricingResponse struct {
	Items    []CalculatedItem
	Subtotal int32
	//DiscountAmount int32
	ShippingFee int32
	TotalAmount int32
}

type PricingItemInput struct {
	PartID   int32
	Name     string
	SKU      string
	Quantity int32
	Price    int32
	WeightG  *int32 // โหลดจาก DB เท่านั้น
}

type CalculatedItem struct {
	PartID    int32
	Name      string
	SKU       string
	Quantity  int32
	UnitPrice int32
	Subtotal  int32
	//DiscountAmount int32
	TotalPrice int32
}
