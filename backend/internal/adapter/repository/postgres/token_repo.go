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

// RefreshTokenRepository implements port.RefreshTokenRepository with PostgreSQL.
type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository.
func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

// Create inserts a new refresh token record.
func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshTokenRecord) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, query,
		token.ID, token.UserID, token.Token, token.ExpiresAt, token.Revoked, token.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

// GetByToken retrieves a refresh token record by its token string.
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshTokenRecord, error) {
	query := `
		SELECT id, user_id, token, expires_at, revoked, created_at
		FROM refresh_tokens WHERE token = $1`

	record := &domain.RefreshTokenRecord{}
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&record.ID, &record.UserID, &record.Token, &record.ExpiresAt,
		&record.Revoked, &record.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return record, nil
}

// RevokeByUserID revokes all refresh tokens for a user.
func (r *RefreshTokenRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1 AND revoked = FALSE`
	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh tokens by user: %w", err)
	}
	return nil
}

// RevokeByToken revokes a specific refresh token.
func (r *RefreshTokenRepository) RevokeByToken(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1`
	_, err := r.pool.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

// DeleteExpired removes all expired refresh tokens.
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}
	return nil
}
