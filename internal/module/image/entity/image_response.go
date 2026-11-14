package entity

import "time"

type ImageRes struct {
	ImageID       int64     `json:"image_id" example:"1"`
	ReferenceID   int32     `json:"reference_id" example:"10"`
	ReferenceType string    `json:"reference_type" example:"part"`
	FileName      string    `json:"file_name" example:"b3c1f0e2-4a5d-4f2b-9c3a-1234567890ab.webp"`
	ImageURL      string    `json:"image_url" example:"/uploads/part/b3c1f0e2-4a5d-4f2b-9c3a-1234567890ab.webp"`
	AltText       string    `json:"alt_text" example:"brake pad front"`
	IsPrimary     bool      `json:"is_primary" example:"true"`
	SortImage     int32     `json:"sort_image" example:"0"`
	CreatedAt     time.Time `json:"created_at"`
}
