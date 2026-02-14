package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
)

// AuthProviderRepository implements port.AuthProviderRepository with PostgreSQL.
type AuthProviderRepository struct {
	pool *pgxpool.Pool
}

// NewAuthProviderRepository creates a new AuthProviderRepository.
func NewAuthProviderRepository(pool *pgxpool.Pool) *AuthProviderRepository {
	return &AuthProviderRepository{pool: pool}
}

// Create inserts a new auth provider link.
func (r *AuthProviderRepository) Create(ctx context.Context, link *domain.AuthProviderLink) error {
	query := `
		INSERT INTO auth_providers (id, user_id, provider, provider_uid, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query,
		link.ID, link.UserID, link.Provider, link.ProviderUID, link.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create auth provider link: %w", err)
	}
	return nil
}

// GetByProviderUID finds a link by provider and external UID.
func (r *AuthProviderRepository) GetByProviderUID(ctx context.Context, provider domain.AuthProvider, uid string) (*domain.AuthProviderLink, error) {
	query := `
		SELECT id, user_id, provider, provider_uid, created_at
		FROM auth_providers WHERE provider = $1 AND provider_uid = $2`

	link := &domain.AuthProviderLink{}
	err := r.pool.QueryRow(ctx, query, provider, uid).Scan(
		&link.ID, &link.UserID, &link.Provider, &link.ProviderUID, &link.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get auth provider link: %w", err)
	}
	return link, nil
}

// GetByUserID returns all provider links for a user.
func (r *AuthProviderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.AuthProviderLink, error) {
	query := `
		SELECT id, user_id, provider, provider_uid, created_at
		FROM auth_providers WHERE user_id = $1`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list auth provider links: %w", err)
	}
	defer rows.Close()

	var links []domain.AuthProviderLink
	for rows.Next() {
		var link domain.AuthProviderLink
		if err := rows.Scan(&link.ID, &link.UserID, &link.Provider, &link.ProviderUID, &link.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan auth provider link: %w", err)
		}
		links = append(links, link)
	}
	return links, nil
}
