package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
	"github.com/vsssp/birthday-app/backend/internal/port"
)

var ErrUserNotFound = errors.New("user not found")

// UserUseCase implements port.UserService.
type UserUseCase struct {
	userRepo port.UserRepository
}

// NewUserUseCase creates a new UserUseCase.
func NewUserUseCase(userRepo port.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// GetByID retrieves a user by their ID.
func (uc *UserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
