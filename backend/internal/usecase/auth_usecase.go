package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
	"github.com/vsssp/birthday-app/backend/internal/pkg/hash"
	jwtpkg "github.com/vsssp/birthday-app/backend/internal/pkg/jwt"
	"github.com/vsssp/birthday-app/backend/internal/port"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// AuthUseCase implements port.AuthService.
type AuthUseCase struct {
	userRepo     port.UserRepository
	providerRepo port.AuthProviderRepository
	tokenRepo    port.RefreshTokenRepository
	jwtService   *jwtpkg.Service
	social       port.SocialVerifier
}

// NewAuthUseCase creates a new AuthUseCase.
func NewAuthUseCase(
	userRepo port.UserRepository,
	providerRepo port.AuthProviderRepository,
	tokenRepo port.RefreshTokenRepository,
	jwtService *jwtpkg.Service,
	social port.SocialVerifier,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		providerRepo: providerRepo,
		tokenRepo:    tokenRepo,
		jwtService:   jwtService,
		social:       social,
	}
}

// Register creates a new user with email and password.
func (uc *AuthUseCase) Register(ctx context.Context, req domain.RegisterRequest) (*domain.TokenPair, error) {
	existing, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashed, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: &hashed,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	link := &domain.AuthProviderLink{
		ID:          uuid.New(),
		UserID:      user.ID,
		Provider:    domain.AuthProviderEmail,
		ProviderUID: user.Email,
		CreatedAt:   now,
	}
	if err := uc.providerRepo.Create(ctx, link); err != nil {
		return nil, err
	}

	return uc.generateTokenPair(ctx, user)
}

// Login authenticates a user with email and password.
func (uc *AuthUseCase) Login(ctx context.Context, req domain.LoginRequest) (*domain.TokenPair, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}
	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}
	if !hash.CheckPassword(req.Password, *user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
	return uc.generateTokenPair(ctx, user)
}

// GoogleLogin authenticates a user via Google ID token.
func (uc *AuthUseCase) GoogleLogin(ctx context.Context, idToken string) (*domain.TokenPair, error) {
	email, name, sub, err := uc.social.VerifyGoogleToken(ctx, idToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return uc.socialLogin(ctx, domain.AuthProviderGoogle, sub, email, name)
}

// AppleLogin authenticates a user via Apple identity token.
func (uc *AuthUseCase) AppleLogin(ctx context.Context, identityToken string) (*domain.TokenPair, error) {
	email, sub, err := uc.social.VerifyAppleToken(ctx, identityToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return uc.socialLogin(ctx, domain.AuthProviderApple, sub, email, "")
}

// RefreshToken generates a new token pair using a valid refresh token.
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	record, err := uc.tokenRepo.GetByToken(ctx, refreshToken)
	if err != nil || record == nil || record.Revoked || record.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}

	_ = uc.tokenRepo.RevokeByToken(ctx, refreshToken)

	user, err := uc.userRepo.GetByID(ctx, record.UserID)
	if err != nil || user == nil {
		return nil, ErrInvalidToken
	}
	return uc.generateTokenPair(ctx, user)
}

func (uc *AuthUseCase) socialLogin(
	ctx context.Context,
	provider domain.AuthProvider,
	providerUID, email, name string,
) (*domain.TokenPair, error) {
	link, _ := uc.providerRepo.GetByProviderUID(ctx, provider, providerUID)
	if link != nil {
		user, err := uc.userRepo.GetByID(ctx, link.UserID)
		if err != nil {
			return nil, err
		}
		return uc.generateTokenPair(ctx, user)
	}

	now := time.Now()
	user, _ := uc.userRepo.GetByEmail(ctx, email)
	if user == nil {
		user = &domain.User{
			ID:        uuid.New(),
			Email:     email,
			Name:      name,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := uc.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	newLink := &domain.AuthProviderLink{
		ID:          uuid.New(),
		UserID:      user.ID,
		Provider:    provider,
		ProviderUID: providerUID,
		CreatedAt:   now,
	}
	if err := uc.providerRepo.Create(ctx, newLink); err != nil {
		return nil, err
	}

	return uc.generateTokenPair(ctx, user)
}

func (uc *AuthUseCase) generateTokenPair(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	accessToken, expiresAt, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiry, err := uc.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	record := &domain.RefreshTokenRecord{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshExpiry,
		CreatedAt: time.Now(),
	}
	if err := uc.tokenRepo.Create(ctx, record); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
