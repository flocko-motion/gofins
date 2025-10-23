package api

import (
	"encoding/json"
	"net/http"

	"github.com/flocko-motion/gofins/pkg/fmp"
	"github.com/go-chi/chi/v5"
)

// handleFMPVerbose toggles verbose FMP request logging
// POST /api/debug/fmp-verbose/true
// POST /api/debug/fmp-verbose/false
func (s *Server) handleFMPVerbose(w http.ResponseWriter, r *http.Request) {
	enabledStr := chi.URLParam(r, "enabled")
	enabled := enabledStr == "true" || enabledStr == "1" || enabledStr == "on"
	
	fmp.EnableVerboseLogging(enabled)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"verbose_logging": enabled,
		"message": "FMP verbose logging updated",
	})
}
