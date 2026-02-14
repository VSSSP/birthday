package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
)

// Service handles JWT token generation and validation.
type Service struct {
	accessSecret  []byte
	refreshSecret []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewService creates a new JWT service.
func NewService(accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration) *Service {
	return &Service{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateAccessToken creates a signed JWT access token.
func (s *Service) GenerateAccessToken(userID uuid.UUID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.accessExpiry)
	claims := domain.AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
		UserID: userID,
		Email:  email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.accessSecret)
	return signedToken, expiresAt, err
}

// GenerateRefreshToken creates a cryptographically secure opaque refresh token.
func (s *Service) GenerateRefreshToken() (string, time.Time, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}
	token := hex.EncodeToString(bytes)
	expiresAt := time.Now().Add(s.refreshExpiry)
	return token, expiresAt, nil
}

// ValidateAccessToken verifies and parses an access token.
func (s *Service) ValidateAccessToken(tokenStr string) (*domain.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*domain.AccessClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
