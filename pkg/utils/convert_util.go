package utils

import (
	"database/sql"
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
)

// NullString → string
func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// NullInt64 → int64
func NullInt64ToInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

// string → NullString
func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// int64 → NullInt64
func Int64ToNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

// time.Time → NullTime
func TimeToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t.UTC(), Valid: true}
}

// sql.NullXXX → pointer
func NullStringToStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func ToInt32Ptr(ni sql.NullInt32) *int {
	if ni.Valid {
		v := int(ni.Int32)
		return &v
	}
	return nil
}

func NullInt32ToInt32Ptr(ni sql.NullInt32) *int32 {
	if ni.Valid {
		return &ni.Int32
	}
	return nil
}

func Int32ToInt32Ptr(ni int32) *int32 {
	if ni != 0 {
		return &ni
	}
	return nil
}

func ToInt64Ptr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func NullTimeToTimePtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

// *string → sql.NullString
func StringPtrToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

// *int → sql.NullInt32
func IntNullToNullInt32(i *int) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: int32(*i), Valid: true}
	}
	return sql.NullInt32{}
}

// *int32 → sql.NullInt32
func Int32PtrToNullInt32(i *int32) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: *i, Valid: true}
	}
	return sql.NullInt32{}
}

// int → sql.NullInt32
func IntToNullInt32(i int) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}
}

// *int → sql.NullInt64
func IntToNullInt64(i *int) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: int64(*i), Valid: true}
	}
	return sql.NullInt64{}
}

func NullOrderStatusToStatus(ns db.NullOrderStatus) *db.OrderStatus {
	if !ns.Valid {
		return nil
	}
	status := ns.OrderStatus
	return &status
}

func ParseBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case []byte:
		return len(val) > 0 && val[0] == 1
	default:
		return false
	}
}
