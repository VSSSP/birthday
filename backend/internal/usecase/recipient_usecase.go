package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
	"github.com/vsssp/birthday-app/backend/internal/port"
)

var (
	ErrRecipientNotFound = errors.New("recipient not found")
	ErrForbidden         = errors.New("access denied")
)

// RecipientUseCase implements port.RecipientService.
type RecipientUseCase struct {
	recipientRepo port.RecipientRepository
}

// NewRecipientUseCase creates a new RecipientUseCase.
func NewRecipientUseCase(recipientRepo port.RecipientRepository) *RecipientUseCase {
	return &RecipientUseCase{recipientRepo: recipientRepo}
}

// Create adds a new recipient for the authenticated user.
func (uc *RecipientUseCase) Create(ctx context.Context, userID uuid.UUID, req domain.CreateRecipientRequest) (*domain.Recipient, error) {
	now := time.Now()
	recipient := &domain.Recipient{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      req.Name,
		Age:       req.Age,
		Gender:    req.Gender,
		MinBudget: req.MinBudget,
		MaxBudget: req.MaxBudget,
		Keywords:  req.Keywords,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if recipient.Keywords == nil {
		recipient.Keywords = []string{}
	}

	if err := uc.recipientRepo.Create(ctx, recipient); err != nil {
		return nil, err
	}
	return recipient, nil
}

// GetByID retrieves a recipient, ensuring it belongs to the requesting user.
func (uc *RecipientUseCase) GetByID(ctx context.Context, userID, recipientID uuid.UUID) (*domain.Recipient, error) {
	recipient, err := uc.recipientRepo.GetByID(ctx, recipientID)
	if err != nil {
		return nil, err
	}
	if recipient == nil {
		return nil, ErrRecipientNotFound
	}
	if recipient.UserID != userID {
		return nil, ErrForbidden
	}
	return recipient, nil
}

// ListByUserID returns all recipients for the authenticated user.
func (uc *RecipientUseCase) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Recipient, error) {
	recipients, err := uc.recipientRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if recipients == nil {
		recipients = []domain.Recipient{}
	}
	return recipients, nil
}

// Update modifies a recipient's fields.
func (uc *RecipientUseCase) Update(ctx context.Context, userID, recipientID uuid.UUID, req domain.UpdateRecipientRequest) (*domain.Recipient, error) {
	recipient, err := uc.recipientRepo.GetByID(ctx, recipientID)
	if err != nil {
		return nil, err
	}
	if recipient == nil {
		return nil, ErrRecipientNotFound
	}
	if recipient.UserID != userID {
		return nil, ErrForbidden
	}

	if req.Name != nil {
		recipient.Name = *req.Name
	}
	if req.Age != nil {
		recipient.Age = *req.Age
	}
	if req.Gender != nil {
		recipient.Gender = *req.Gender
	}
	if req.MinBudget != nil {
		recipient.MinBudget = *req.MinBudget
	}
	if req.MaxBudget != nil {
		recipient.MaxBudget = *req.MaxBudget
	}
	if req.Keywords != nil {
		recipient.Keywords = *req.Keywords
	}
	recipient.UpdatedAt = time.Now()

	if err := uc.recipientRepo.Update(ctx, recipient); err != nil {
		return nil, err
	}
	return recipient, nil
}

// Delete removes a recipient, ensuring it belongs to the requesting user.
func (uc *RecipientUseCase) Delete(ctx context.Context, userID, recipientID uuid.UUID) error {
	recipient, err := uc.recipientRepo.GetByID(ctx, recipientID)
	if err != nil {
		return err
	}
	if recipient == nil {
		return ErrRecipientNotFound
	}
	if recipient.UserID != userID {
		return ErrForbidden
	}
	return uc.recipientRepo.Delete(ctx, recipientID)
}

// BulkDelete removes multiple recipients belonging to the authenticated user.
func (uc *RecipientUseCase) BulkDelete(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	return uc.recipientRepo.BulkDelete(ctx, userID, ids)
}
