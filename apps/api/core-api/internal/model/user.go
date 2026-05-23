package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	KeycloakID string
	Email      string
	Role       string
	FirstName  *string
	LastName   *string
	Phone      *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Borrower struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	INN          string
	BIN          *string
	OrgName      *string
	ActivityType *string
	FarmID       *uuid.UUID
	CreatedAt    time.Time
}

// RolePriority returns a numeric weight for role comparison.
// Higher = more privileged.
func RolePriority(role string) int {
	switch role {
	case "admin":
		return 4
	case "expert":
		return 3
	case "employee":
		return 2
	case "borrower":
		return 1
	default:
		return 0
	}
}

// HighestRole picks the most privileged role from a slice of Keycloak roles.
func HighestRole(roles []string) string {
	best := "public"
	for _, r := range roles {
		if RolePriority(r) > RolePriority(best) {
			best = r
		}
	}
	return best
}
