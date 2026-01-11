package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flocko-motion/gofins/pkg/db"
)

func (s *Server) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := getUserID(r)

	// Get user from database (includes is_admin from DB)
	user, err := db.GetUserByID(r.Context(), userID)
	if err != nil {
		fmt.Printf("[API] Error getting user by ID '%s': %v\n", userID, err)
		_ = db.Db().LogError("api.get_current_user", "error", "Failed to get user by ID", map[string]interface{}{"user_id": userID.String(), "error": err.Error()})
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		fmt.Printf("[API] User not found for ID '%s'\n", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
