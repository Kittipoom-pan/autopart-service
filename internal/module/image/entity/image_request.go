package entity

type UploadImageReq struct {
	ReferenceID   int32  `form:"reference_id" json:"reference_id" example:"10"`
	ReferenceType string `form:"reference_type" json:"reference_type" example:"part"`
	AltText       string `form:"alt_text" json:"alt_text" example:"brake pad front"`
	// FilePrefix is used internally to build file names (e.g. part SKU).
	// It is not bound from the request payload.
	FilePrefix string `form:"-" json:"-"`
}
