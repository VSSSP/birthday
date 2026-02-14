package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
)

// UserRepository defines the data access methods for users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

// AuthProviderRepository defines the data access methods for auth provider links.
type AuthProviderRepository interface {
	Create(ctx context.Context, link *domain.AuthProviderLink) error
	GetByProviderUID(ctx context.Context, provider domain.AuthProvider, uid string) (*domain.AuthProviderLink, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.AuthProviderLink, error)
}

// RefreshTokenRepository defines the data access methods for refresh tokens.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshTokenRecord) error
	GetByToken(ctx context.Context, token string) (*domain.RefreshTokenRecord, error)
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
	RevokeByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
}

// RecipientRepository defines the data access methods for recipients.
type RecipientRepository interface {
	Create(ctx context.Context, recipient *domain.Recipient) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipient, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Recipient, error)
	Update(ctx context.Context, recipient *domain.Recipient) error
	Delete(ctx context.Context, id uuid.UUID) error
	BulkDelete(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
}
