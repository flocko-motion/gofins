package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/flocko-motion/gofins/pkg/db/generated"
	"github.com/flocko-motion/gofins/pkg/f"
	"github.com/flocko-motion/gofins/pkg/types"
	"github.com/google/uuid"
)

// CreateUser creates a new user with a UUID derived from their name
// If this is the first user, they are made admin
func CreateUser(ctx context.Context, name string) (*types.User, error) {
	// Check if this is the first user
	userCount, err := genQ().CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	isAdmin := userCount == 0

	// Generate stable UUID from username (hash-based)
	userID := f.StringToUUID(name)

	genUser, err := genQ().CreateUser(ctx, generated.CreateUserParams{
		ID:      userID,
		Name:    name,
		IsAdmin: isAdmin,
	})
	if err != nil {
		return nil, err
	}

	return &types.User{
		ID:        genUser.ID,
		Name:      genUser.Name,
		CreatedAt: genUser.CreatedAt,
		IsAdmin:   genUser.IsAdmin,
	}, nil
}

// GetUser retrieves a user by name
func GetUser(ctx context.Context, name string) (*types.User, error) {
	genUser, err := genQ().GetUser(ctx, name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &types.User{
		ID:        genUser.ID,
		Name:      genUser.Name,
		CreatedAt: genUser.CreatedAt,
		IsAdmin:   genUser.IsAdmin,
	}, nil
}

// GetUserByID retrieves a user by UUID
func GetUserByID(ctx context.Context, id uuid.UUID) (*types.User, error) {
	genUser, err := genQ().GetUserByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &types.User{
		ID:        genUser.ID,
		Name:      genUser.Name,
		CreatedAt: genUser.CreatedAt,
		IsAdmin:   genUser.IsAdmin,
	}, nil
}

// ListUsers returns all users
func ListUsers(ctx context.Context) ([]types.User, error) {
	genUsers, err := genQ().ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]types.User, len(genUsers))
	for i, u := range genUsers {
		users[i] = types.User{
			ID:        u.ID,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			IsAdmin:   u.IsAdmin,
		}
	}
	return users, nil
}

// DeleteUser deletes a user and all their data (ratings, favorites)
func DeleteUser(ctx context.Context, name string) error {
	// Get user to find their UUID
	user, err := GetUser(ctx, name)
	if err != nil {
		return err
	}
	if user == nil {
		return nil // User doesn't exist, nothing to delete
	}

	// Start transaction
	tx, err := Db().conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := genQ().WithTx(tx)

	// Delete user's ratings
	if err := q.DeleteUserRatings(ctx, user.ID); err != nil {
		return err
	}

	// Delete user's favorites
	if err := q.DeleteUserFavorites(ctx, user.ID); err != nil {
		return err
	}

	// Delete user
	if err := q.DeleteUser(ctx, user.ID); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// UpdateUserAdmin updates the admin status of a user by name or UUID
func UpdateUserAdmin(ctx context.Context, nameOrID string, isAdmin bool) (*types.User, error) {
	// Try to parse as UUID first
	userID, err := uuid.Parse(nameOrID)
	var user *types.User

	if err == nil {
		// It's a valid UUID, get by ID
		user, err = GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}
	} else {
		// Not a UUID, treat as name
		user, err = GetUser(ctx, nameOrID)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		return nil, sql.ErrNoRows
	}

	// Update admin status
	genUser, err := genQ().UpdateUserAdmin(ctx, generated.UpdateUserAdminParams{
		IsAdmin: isAdmin,
		ID:      user.ID,
	})
	if err != nil {
		return nil, err
	}

	return &types.User{
		ID:        genUser.ID,
		Name:      genUser.Name,
		CreatedAt: genUser.CreatedAt,
		IsAdmin:   genUser.IsAdmin,
	}, nil
}

// UserRating represents a rating given to a symbol
type UserRating struct {
	ID        int       `json:"id"`
	Ticker    string    `json:"ticker"`
	Rating    int       `json:"rating"` // -5 to +5
	Notes     *string   `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
}

// ToggleFavorite adds or removes a symbol from favorites
func ToggleFavorite(ctx context.Context, userID uuid.UUID, ticker string) (bool, error) {
	// Check if already favorited
	exists, err := genQ().IsFavorite(ctx, generated.IsFavoriteParams{
		UserID: userID,
		Ticker: ticker,
	})
	if err != nil {
		return false, err
	}

	if exists {
		// Remove from favorites
		err = genQ().RemoveFavorite(ctx, generated.RemoveFavoriteParams{
			UserID: userID,
			Ticker: ticker,
		})
		return false, err
	} else {
		// Add to favorites
		err = genQ().AddFavorite(ctx, generated.AddFavoriteParams{
			UserID: userID,
			Ticker: ticker,
		})
		return true, err
	}
}

// IsFavorite checks if a symbol is favorited
func IsFavorite(ctx context.Context, userID uuid.UUID, ticker string) (bool, error) {
	return genQ().IsFavorite(ctx, generated.IsFavoriteParams{
		UserID: userID,
		Ticker: ticker,
	})
}

// GetFavorites returns all favorited tickers
func GetFavorites(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return genQ().GetFavorites(ctx, userID)
}

// AddRating adds a new rating for a symbol
func AddRating(ctx context.Context, userID uuid.UUID, ticker string, rating int, notes *string) (*UserRating, error) {
	if rating < -5 || rating > 5 {
		return nil, sql.ErrNoRows
	}

	var sqlNotes sql.NullString
	if notes != nil {
		sqlNotes = sql.NullString{String: *notes, Valid: true}
	}

	genRating, err := genQ().AddRating(ctx, generated.AddRatingParams{
		UserID: userID,
		Ticker: ticker,
		Rating: int32(rating),
		Notes:  sqlNotes,
	})
	if err != nil {
		return nil, err
	}

	var resultNotes *string
	if genRating.Notes.Valid {
		resultNotes = &genRating.Notes.String
	}

	return &UserRating{
		ID:        int(genRating.ID),
		Ticker:    genRating.Ticker,
		Rating:    int(genRating.Rating),
		Notes:     resultNotes,
		CreatedAt: genRating.CreatedAt,
	}, nil
}

// GetLatestRating returns the most recent rating for a symbol
func GetLatestRating(ctx context.Context, userID uuid.UUID, ticker string) (*UserRating, error) {
	genRating, err := genQ().GetLatestRating(ctx, generated.GetLatestRatingParams{
		UserID: userID,
		Ticker: ticker,
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var resultNotes *string
	if genRating.Notes.Valid {
		resultNotes = &genRating.Notes.String
	}

	return &UserRating{
		ID:        int(genRating.ID),
		Ticker:    genRating.Ticker,
		Rating:    int(genRating.Rating),
		Notes:     resultNotes,
		CreatedAt: genRating.CreatedAt,
	}, nil
}

// GetRatingHistory returns all ratings for a symbol
func GetRatingHistory(ctx context.Context, userID uuid.UUID, ticker string) ([]UserRating, error) {
	genRatings, err := genQ().GetRatingHistory(ctx, generated.GetRatingHistoryParams{
		UserID: userID,
		Ticker: ticker,
	})
	if err != nil {
		return nil, err
	}

	result := make([]UserRating, len(genRatings))
	for i, r := range genRatings {
		var resultNotes *string
		if r.Notes.Valid {
			resultNotes = &r.Notes.String
		}
		result[i] = UserRating{
			ID:        int(r.ID),
			Ticker:    r.Ticker,
			Rating:    int(r.Rating),
			Notes:     resultNotes,
			CreatedAt: r.CreatedAt,
		}
	}
	return result, nil
}

// GetAllLatestRatings returns the latest rating for each rated symbol
func GetAllLatestRatings(ctx context.Context, userID uuid.UUID) (map[string]*UserRating, error) {
	genRatings, err := genQ().GetAllLatestRatings(ctx, userID)
	if err != nil {
		return nil, err
	}

	ratings := make(map[string]*UserRating)
	for _, r := range genRatings {
		var resultNotes *string
		if r.Notes.Valid {
			resultNotes = &r.Notes.String
		}
		rating := UserRating{
			ID:        int(r.ID),
			Ticker:    r.Ticker,
			Rating:    int(r.Rating),
			Notes:     resultNotes,
			CreatedAt: r.CreatedAt,
		}
		ratings[rating.Ticker] = &rating
	}
	return ratings, nil
}

// DeleteRating deletes a rating by ID (for the given user)
func DeleteRating(ctx context.Context, userID uuid.UUID, id int) error {
	return genQ().DeleteRating(ctx, generated.DeleteRatingParams{
		UserID: userID,
		ID:     int32(id),
	})
}

// GetAllNotesChronological returns all ratings that have notes, sorted by creation time (newest first)
func GetAllNotesChronological(ctx context.Context, userID uuid.UUID) ([]UserRating, error) {
	genRatings, err := genQ().GetAllNotesChronological(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]UserRating, len(genRatings))
	for i, r := range genRatings {
		var resultNotes *string
		if r.Notes.Valid {
			resultNotes = &r.Notes.String
		}
		result[i] = UserRating{
			ID:        int(r.ID),
			Ticker:    r.Ticker,
			Rating:    int(r.Rating),
			Notes:     resultNotes,
			CreatedAt: r.CreatedAt,
		}
	}
	return result, nil
}
