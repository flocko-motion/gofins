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

// MaybeStringToNullString converts *string to sql.NullString
func MaybeStringToNullString(ptr *string) sql.NullString {
	if ptr == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *ptr, Valid: true}
}

// NullStringToMaybeString converts sql.NullString to *string
func NullStringToMaybeString(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	return &n.String
}

// StringToNullString converts string to sql.NullString
func StringToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

// NullBoolToMaybeBool converts sql.NullBool to *bool
func NullBoolToMaybeBool(n sql.NullBool) *bool {
	if !n.Valid {
		return nil
	}
	return &n.Bool
}

// MaybeBoolToNullBool converts *bool to sql.NullBool
func MaybeBoolToNullBool(ptr *bool) sql.NullBool {
	if ptr == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *ptr, Valid: true}
}

// NullFloat64ToMaybeFloat64 converts sql.NullFloat64 to *float64
func NullFloat64ToMaybeFloat64(n sql.NullFloat64) *float64 {
	if !n.Valid {
		return nil
	}
	return &n.Float64
}

// MaybeFloat64ToNullFloat64 converts *float64 to sql.NullFloat64
func MaybeFloat64ToNullFloat64(ptr *float64) sql.NullFloat64 {
	if ptr == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *ptr, Valid: true}
}
