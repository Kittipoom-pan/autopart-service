package entity

import "time"

type Image struct {
	ImageID       int64
	ReferenceID   int32
	ReferenceType string
	FileName      string
	ImageURL      string
	AltText       string
	IsPrimary     bool
	SortImage     int32
	CreatedAt     time.Time
}
