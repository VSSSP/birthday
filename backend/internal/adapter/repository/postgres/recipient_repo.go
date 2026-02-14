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

// RecipientRepository implements port.RecipientRepository with PostgreSQL.
type RecipientRepository struct {
	pool *pgxpool.Pool
}

// NewRecipientRepository creates a new RecipientRepository.
func NewRecipientRepository(pool *pgxpool.Pool) *RecipientRepository {
	return &RecipientRepository{pool: pool}
}

// Create inserts a new recipient.
func (r *RecipientRepository) Create(ctx context.Context, recipient *domain.Recipient) error {
	query := `
		INSERT INTO recipients (id, user_id, name, age, gender, min_budget, max_budget, keywords, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		recipient.ID, recipient.UserID, recipient.Name, recipient.Age, recipient.Gender,
		recipient.MinBudget, recipient.MaxBudget, recipient.Keywords,
		recipient.CreatedAt, recipient.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create recipient: %w", err)
	}
	return nil
}

// GetByID retrieves a recipient by ID.
func (r *RecipientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipient, error) {
	query := `
		SELECT id, user_id, name, age, gender, min_budget, max_budget, keywords, created_at, updated_at
		FROM recipients WHERE id = $1`

	rec := &domain.Recipient{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rec.ID, &rec.UserID, &rec.Name, &rec.Age, &rec.Gender,
		&rec.MinBudget, &rec.MaxBudget, &rec.Keywords,
		&rec.CreatedAt, &rec.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get recipient: %w", err)
	}
	return rec, nil
}

// ListByUserID returns all recipients belonging to a user.
func (r *RecipientRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Recipient, error) {
	query := `
		SELECT id, user_id, name, age, gender, min_budget, max_budget, keywords, created_at, updated_at
		FROM recipients WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list recipients: %w", err)
	}
	defer rows.Close()

	var recipients []domain.Recipient
	for rows.Next() {
		var rec domain.Recipient
		if err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.Name, &rec.Age, &rec.Gender,
			&rec.MinBudget, &rec.MaxBudget, &rec.Keywords,
			&rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan recipient: %w", err)
		}
		recipients = append(recipients, rec)
	}
	return recipients, nil
}

// Update modifies a recipient's fields.
func (r *RecipientRepository) Update(ctx context.Context, recipient *domain.Recipient) error {
	query := `
		UPDATE recipients
		SET name = $2, age = $3, gender = $4, min_budget = $5, max_budget = $6,
		    keywords = $7, updated_at = $8
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		recipient.ID, recipient.Name, recipient.Age, recipient.Gender,
		recipient.MinBudget, recipient.MaxBudget, recipient.Keywords,
		recipient.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update recipient: %w", err)
	}
	return nil
}

// Delete removes a recipient by ID.
func (r *RecipientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM recipients WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipient: %w", err)
	}
	return nil
}

// BulkDelete removes multiple recipients belonging to a user.
func (r *RecipientRepository) BulkDelete(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	query := `DELETE FROM recipients WHERE user_id = $1 AND id = ANY($2)`
	_, err := r.pool.Exec(ctx, query, userID, ids)
	if err != nil {
		return fmt.Errorf("failed to bulk delete recipients: %w", err)
	}
	return nil
}
