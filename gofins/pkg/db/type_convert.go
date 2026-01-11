package db

import (
	"github.com/flocko-motion/gofins/pkg/db/generated"
	"github.com/flocko-motion/gofins/pkg/types"
)

// ToUser converts generated.User to types.User
func ToUser(u generated.User) types.User {
	return types.User{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		IsAdmin:   u.IsAdmin,
	}
}

// ToUserPtr converts generated.User to *types.User
func ToUserPtr(u generated.User) *types.User {
	result := ToUser(u)
	return &result
}

// ToUsers converts []generated.User to []types.User
func ToUsers(users []generated.User) []types.User {
	result := make([]types.User, len(users))
	for i, u := range users {
		result[i] = ToUser(u)
	}
	return result
}

// Add more converters as you migrate other types:
// - ToTypesSymbol
// - ToTypesPriceData
// - ToTypesAnalysisPackage
// etc.
