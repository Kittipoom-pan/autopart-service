package utils

import "github.com/google/uuid"

func NewUUIDBytes() []byte {
	id := uuid.New()
	return id[:]
}
