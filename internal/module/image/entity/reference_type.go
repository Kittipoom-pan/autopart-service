package entity

type ReferenceType string

const (
	ReferenceTypePart      ReferenceType = "part"
	ReferenceTypeCustomer  ReferenceType = "customer"
	ReferenceTypeOrder     ReferenceType = "order"
	ReferenceTypeCarBrand  ReferenceType = "car_brand"
	ReferenceTypeCarModel  ReferenceType = "car_model"
	ReferenceTypePartBrand ReferenceType = "part_brand"
	ReferenceTypePartType  ReferenceType = "part_type"
)

func (r ReferenceType) IsValid() bool {
	switch r {
	case ReferenceTypePart,
		ReferenceTypeCustomer,
		ReferenceTypeOrder,
		ReferenceTypeCarBrand,
		ReferenceTypeCarModel,
		ReferenceTypePartBrand,
		ReferenceTypePartType:
		return true
	default:
		return false
	}
}
