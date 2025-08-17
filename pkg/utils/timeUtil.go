package utils

import (
	"database/sql"
	"time"
)

func NowUTC() time.Time {
	return time.Now().UTC()
}

func NullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// return null time UTC
func NullTimeNow() sql.NullTime {
	return NullTime(NowUTC())
}
