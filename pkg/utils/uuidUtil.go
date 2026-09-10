package utils

import "github.com/google/uuid"

func NewUUIDBytes() []byte {
	id := uuid.New()
	return id[:]
}

func UUIDBytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	id, err := uuid.FromBytes(b)
	if err != nil {
		return ""
	}
	return id.String()
}
