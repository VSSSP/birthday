package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
)

// AuthService defines the business logic for authentication.
type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.TokenPair, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.TokenPair, error)
	GoogleLogin(ctx context.Context, idToken string) (*domain.TokenPair, error)
	AppleLogin(ctx context.Context, identityToken string) (*domain.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
}

// UserService defines the business logic for user operations.
type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

// RecipientService defines the business logic for recipient operations.
type RecipientService interface {
	Create(ctx context.Context, userID uuid.UUID, req domain.CreateRecipientRequest) (*domain.Recipient, error)
	GetByID(ctx context.Context, userID, recipientID uuid.UUID) (*domain.Recipient, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Recipient, error)
	Update(ctx context.Context, userID, recipientID uuid.UUID, req domain.UpdateRecipientRequest) (*domain.Recipient, error)
	Delete(ctx context.Context, userID, recipientID uuid.UUID) error
	BulkDelete(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
}

// SocialVerifier defines the interface for verifying social login tokens.
type SocialVerifier interface {
	VerifyGoogleToken(ctx context.Context, idToken string) (email, name, sub string, err error)
	VerifyAppleToken(ctx context.Context, identityToken string) (email, sub string, err error)
}
