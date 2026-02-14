package domain

import (
	"time"

	"github.com/google/uuid"
)

// Recipient represents a person the user wants to buy a gift for.
type Recipient struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	MinBudget float64   `json:"min_budget"`
	MaxBudget float64   `json:"max_budget"`
	Keywords  []string  `json:"keywords"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRecipientRequest is the payload for creating a recipient.
type CreateRecipientRequest struct {
	Name      string   `json:"name"`
	Age       int      `json:"age"`
	Gender    string   `json:"gender"`
	MinBudget float64  `json:"min_budget"`
	MaxBudget float64  `json:"max_budget"`
	Keywords  []string `json:"keywords"`
}

// UpdateRecipientRequest is the payload for updating a recipient.
type UpdateRecipientRequest struct {
	Name      *string   `json:"name"`
	Age       *int      `json:"age"`
	Gender    *string   `json:"gender"`
	MinBudget *float64  `json:"min_budget"`
	MaxBudget *float64  `json:"max_budget"`
	Keywords  *[]string `json:"keywords"`
}

// BulkDeleteRequest is the payload for deleting multiple recipients.
type BulkDeleteRequest struct {
	IDs []uuid.UUID `json:"ids"`
}
