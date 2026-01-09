package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/flocko-motion/gofins/pkg/analysis"
	"github.com/flocko-motion/gofins/pkg/db"
	"github.com/flocko-motion/gofins/pkg/types"
	"github.com/go-chi/chi/v5"
)

// UpdateAnalysisRequest represents the request body for updating an analysis
type UpdateAnalysisRequest struct {
	Name string `json:"name"`
}

// handleAnalyses handles REST operations on /api/analyses
// GET  /api/analyses - List all analyses
// POST /api/analyses - Create new analysis
func (s *Server) handleAnalyses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleListAnalyses(w, r)
	case http.MethodPost:
		s.handleCreateAnalysis(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleListAnalyses lists all analysis packages
func (s *Server) handleListAnalyses(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	packages, err := analysis.ListPackages(userID)
	if err != nil {
		http.Error(w, "Failed to list analyses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ensure we return an empty array instead of null
	if packages == nil {
		packages = []types.AnalysisPackage{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}

// handleGetAnalysis retrieves a single analysis package
// GET /api/analysis/{id}
func (s *Server) handleGetAnalysis(w http.ResponseWriter, r *http.Request) {
	packageID := chi.URLParam(r, "id")
	if packageID == "" {
		http.Error(w, "Package ID required", http.StatusBadRequest)
		return
	}
	userID := getUserID(r)
	pkg, err := analysis.GetPackage(userID, packageID)
	if err != nil {
		http.Error(w, "Failed to get analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if pkg == nil {
		http.Error(w, "Analysis not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkg)
}

// handleUpdateAnalysis updates an analysis package (currently just name)
// PUT /api/analysis/{id}
func (s *Server) handleUpdateAnalysis(w http.ResponseWriter, r *http.Request) {
	packageID := chi.URLParam(r, "id")
	if packageID == "" {
		http.Error(w, "Package ID required", http.StatusBadRequest)
		return
	}
	var req UpdateAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	userID := getUserID(r)
	pkg, err := analysis.UpdatePackageName(userID, packageID, req.Name)
	if err != nil {
		http.Error(w, "Failed to update analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if pkg == nil {
		http.Error(w, "Analysis not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkg)
}

// handleDeleteAnalysis deletes an analysis package
// DELETE /api/analysis/{id}
func (s *Server) handleDeleteAnalysis(w http.ResponseWriter, r *http.Request) {
	packageID := chi.URLParam(r, "id")
	if packageID == "" {
		http.Error(w, "Package ID required", http.StatusBadRequest)
		return
	}
	userID := getUserID(r)
	err := analysis.DeletePackage(userID, packageID)
	if err != nil {
		http.Error(w, "Failed to delete analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleAnalysisResults retrieves all results for an analysis package
// GET /api/analysis/{id}/results
func (s *Server) handleAnalysisResults(w http.ResponseWriter, r *http.Request) {
	packageID := chi.URLParam(r, "id")
	if packageID == "" {
		http.Error(w, "Package ID required", http.StatusBadRequest)
		return
	}

	userID := getUserID(r)
	results, err := db.GetAnalysisResults(userID, packageID)
	if err != nil {
		http.Error(w, "Failed to get results: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return empty array instead of null
	if results == nil {
		results = []types.AnalysisResult{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// handleAnalysisChart serves PNG chart images
// GET /api/analysis/{id}/chart/{ticker}
func (s *Server) handleAnalysisChart(w http.ResponseWriter, r *http.Request) {
	packageID := chi.URLParam(r, "id")
	ticker := chi.URLParam(r, "ticker")

	if packageID == "" || ticker == "" {
		http.Error(w, "Package ID and ticker required", http.StatusBadRequest)
		return
	}

	packageID = strings.ReplaceAll(packageID, "..", "")
	ticker = strings.ReplaceAll(ticker, "..", "")

	// Construct chart path
	chartPath := analysis.PathPlot(packageID, analysis.PlotTypeChart, ticker)

	// Check if file exists
	if _, err := os.Stat(chartPath); os.IsNotExist(err) {
		http.Error(w, "Chart not found", http.StatusNotFound)
		return
	}

	// Serve the PNG file
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, chartPath)
}

// handleAnalysisHistogram serves PNG histogram images
// GET /api/analysis/{id}/histogram/{ticker}
func (s *Server) handleAnalysisHistogram(w http.ResponseWriter, r *http.Request) {
	packageID := chi.URLParam(r, "id")
	ticker := chi.URLParam(r, "ticker")

	if packageID == "" || ticker == "" {
		http.Error(w, "Package ID and ticker required", http.StatusBadRequest)
		return
	}

	packageID = strings.ReplaceAll(packageID, "..", "")
	ticker = strings.ReplaceAll(ticker, "..", "")

	// Construct histogram path
	histogramPath := analysis.PathPlot(packageID, analysis.PlotTypeHistogram, ticker)

	// Check if file exists
	if _, err := os.Stat(histogramPath); os.IsNotExist(err) {
		http.Error(w, "Histogram not found", http.StatusNotFound)
		return
	}

	// Serve the PNG file
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, histogramPath)
}

// handleSymbolProfile retrieves profile information for a symbol
// GET /api/analysis/{id}/profile/{ticker}
func (s *Server) handleSymbolProfile(w http.ResponseWriter, r *http.Request) {
	packageID := chi.URLParam(r, "id")
	ticker := chi.URLParam(r, "ticker")

	if packageID == "" || ticker == "" {
		http.Error(w, "Package ID and ticker required", http.StatusBadRequest)
		return
	}

	packageID = strings.ReplaceAll(packageID, "..", "")
	ticker = strings.ReplaceAll(ticker, "..", "")

	// Get symbol profile from database
	symbol, err := db.GetSymbol(ticker)
	if err != nil {
		http.Error(w, "Failed to get symbol profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if symbol == nil {
		http.Error(w, "Symbol not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(symbol)
}
