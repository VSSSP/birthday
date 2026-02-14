package domain

import (
	"time"

	"github.com/google/uuid"
)

// AuthProvider represents the authentication method used.
type AuthProvider string

const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
	AuthProviderApple  AuthProvider = "apple"
)

// User represents a registered user.
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash *string   `json:"-"`
	AvatarURL    *string   `json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AuthProviderLink represents a link between a user and an auth provider.
type AuthProviderLink struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"user_id"`
	Provider    AuthProvider `json:"provider"`
	ProviderUID string       `json:"provider_uid"`
	CreatedAt   time.Time    `json:"created_at"`
}
