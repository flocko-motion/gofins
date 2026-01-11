package types

import "time"

// PriceData represents price data (daily, weekly, or monthly)
// Stores both original currency and USD-converted values
type PriceData struct {
	Date         time.Time
	Open         float64  // USD converted
	High         float64  // USD converted
	Low          float64  // USD converted
	Avg          float64  // USD converted
	Close        float64  // USD converted
	YoY          *float64 // Percentage, currency-independent
	SymbolTicker string
	// Original currency values (before USD conversion)
	OpenOrig  *float64 // nil if already in USD
	HighOrig  *float64
	LowOrig   *float64
	AvgOrig   *float64
	CloseOrig *float64
}

// PriceInterval represents the time interval for price data
type PriceInterval string

const (
	IntervalMonthly PriceInterval = "monthly"
	IntervalWeekly  PriceInterval = "weekly"
)

// Price status constants (mirror updater statuses)
const (
	StatusOK       = "ok"
	StatusNotFound = "not_found"
	StatusFailed   = "failed"
)
