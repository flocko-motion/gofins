package f

import (
	"database/sql"
	"time"
)

// MaybeInt64ToNullInt64 converts *int64 to sql.NullInt64
func MaybeInt64ToNullInt64(ptr *int64) sql.NullInt64 {
	if ptr == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *ptr, Valid: true}
}

// NullInt64ToMaybeInt64 converts sql.NullInt64 to *int64
func NullInt64ToMaybeInt64(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return &n.Int64
}

// MaybeTimeToNullTime converts *time.Time to sql.NullTime
func MaybeTimeToNullTime(ptr *time.Time) sql.NullTime {
	if ptr == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *ptr, Valid: true}
}

// NullTimeToMaybeTime converts sql.NullTime to *time.Time
func NullTimeToMaybeTime(n sql.NullTime) *time.Time {
	if !n.Valid {
		return nil
	}
	return &n.Time
}
