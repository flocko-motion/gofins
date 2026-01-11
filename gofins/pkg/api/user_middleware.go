package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flocko-motion/gofins/pkg/db"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

// userMiddleware extracts user from --user flag (dev) or X-Remote-User header (production)
func (s *Server) userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username string

		// 1. Check if dev user override is set (--user flag)
		if s.devUser != "" {
			username = s.devUser
		} else if user := r.Header.Get("X-Remote-User"); user != "" {
			// 2. Check X-Remote-User header (from Apache .htaccess auth)
			fmt.Printf("[API] X-Remote-User: %s\n", user)
			username = user
		} else {
			// 3. No authentication found - this should not happen in production
			fmt.Printf("[API] No authentication: no --user flag and no X-Remote-User header\n")
			_ = db.Db().LogError("api.user_middleware", "error", "No authentication provided", map[string]interface{}{"path": r.URL.Path})
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// 4. Get or create user (auto-creates on first access)
		user, err := db.GetUser(r.Context(), username)
		if err != nil {
			fmt.Printf("[API] Error getting user '%s': %v\n", username, err)
			_ = db.Db().LogError("api.user_middleware", "error", "Failed to get user", map[string]interface{}{"username": username, "error": err.Error()})
			http.Error(w, "Authentication error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			// User doesn't exist, create them
			user, err = db.CreateUser(r.Context(), username)
			if err != nil {
				fmt.Printf("[API] Error creating user '%s': %v\n", username, err)
				_ = db.Db().LogError("api.user_middleware", "error", "Failed to create user", map[string]interface{}{"username": username, "error": err.Error()})
				http.Error(w, "Authentication error", http.StatusInternalServerError)
				return
			}
		}

		// 5. Store user ID in context
		ctx := context.WithValue(r.Context(), userIDKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getUserID extracts the user ID from the request context
func getUserID(r *http.Request) uuid.UUID {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		// This should never happen if middleware is properly configured
		return uuid.Nil
	}
	return userID
}

// adminOnlyMiddleware restricts access to admin users (checks is_admin from database)
func (s *Server) adminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get current user ID from context
		userID := getUserID(r)

		// Get user from database to check is_admin field
		user, err := db.GetUserByID(r.Context(), userID)
		if err != nil {
			fmt.Printf("[API] Error getting user for admin check: %v\n", err)
			_ = db.Db().LogError("api.admin_middleware", "error", "Failed to get user", map[string]interface{}{"user_id": userID.String(), "error": err.Error()})
			http.Error(w, "Authentication error", http.StatusInternalServerError)
			return
		}

		if user == nil {
			fmt.Printf("[API] User not found for admin check: %s\n", userID)
			http.Error(w, "Authentication error", http.StatusInternalServerError)
			return
		}

		// Check if user is admin
		if !user.IsAdmin {
			http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
