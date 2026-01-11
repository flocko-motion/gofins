package api

import (
	"encoding/json"
	"net/http"

	"github.com/flocko-motion/gofins/pkg/db"
)

func (s *Server) handleListSymbols(w http.ResponseWriter, r *http.Request) {
	tickers, err := db.GetAllTickers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":   len(tickers),
		"tickers": tickers,
	})
}
