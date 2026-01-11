package api

import (
	"encoding/json"
	"net/http"

	"github.com/flocko-motion/gofins/pkg/db"
	"github.com/flocko-motion/gofins/pkg/f"
)

type HealthResponse struct {
	Status           string  `json:"status"`
	TotalSymbols     int     `json:"total_symbols"`
	ActivelyTrading  int     `json:"actively_trading"`
	StaleProfiles    int     `json:"stale_profiles"`
	StalePrices      int     `json:"stale_prices"`
	OldestProfile    *string `json:"oldest_profile"`
	OldestPrice      *string `json:"oldest_price"`
	ProfileThreshold string  `json:"profile_threshold"`
	PriceThreshold   string  `json:"price_threshold"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	const timeFormat = "2006-01-02T15:04:05Z"

	response := HealthResponse{
		Status:           "ok",
		TotalSymbols:     f.First(db.CountSymbols(r.Context())),
		ActivelyTrading:  f.First(db.CountActivelyTrading(r.Context())),
		StaleProfiles:    f.First(db.CountStaleProfiles(r.Context())),
		StalePrices:      f.First(db.CountStalePrices(r.Context())),
		OldestProfile:    f.MaybeDateToMaybeString(f.First(db.GetOldestProfileUpdate(r.Context())), timeFormat),
		OldestPrice:      f.MaybeDateToMaybeString(f.First(db.GetOldestPriceUpdate(r.Context())), timeFormat),
		ProfileThreshold: db.GetProfileThreshold().Format(timeFormat),
		PriceThreshold:   db.GetPriceThreshold().Format(timeFormat),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
