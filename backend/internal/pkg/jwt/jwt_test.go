package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
	svc := NewService("test-access-secret-32-chars-long!", "test-refresh-secret", 15*time.Minute, 7*24*time.Hour)

	userID := uuid.New()
	email := "test@example.com"

	token, expiresAt, err := svc.GenerateAccessToken(userID, email)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))

	claims, err := svc.ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestValidateAccessToken_Invalid(t *testing.T) {
	svc := NewService("test-access-secret-32-chars-long!", "test-refresh-secret", 15*time.Minute, 7*24*time.Hour)

	_, err := svc.ValidateAccessToken("invalid-token")
	assert.Error(t, err)
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	svc1 := NewService("secret-one-32-chars-long-enough!", "refresh", 15*time.Minute, 7*24*time.Hour)
	svc2 := NewService("secret-two-32-chars-long-enough!", "refresh", 15*time.Minute, 7*24*time.Hour)

	token, _, err := svc1.GenerateAccessToken(uuid.New(), "test@example.com")
	require.NoError(t, err)

	_, err = svc2.ValidateAccessToken(token)
	assert.Error(t, err)
}

func TestGenerateRefreshToken(t *testing.T) {
	svc := NewService("access", "refresh", 15*time.Minute, 7*24*time.Hour)

	token, expiresAt, err := svc.GenerateRefreshToken()
	require.NoError(t, err)
	assert.Len(t, token, 64) // 32 bytes = 64 hex chars
	assert.True(t, expiresAt.After(time.Now()))
}

func TestGenerateRefreshToken_Unique(t *testing.T) {
	svc := NewService("access", "refresh", 15*time.Minute, 7*24*time.Hour)

	token1, _, _ := svc.GenerateRefreshToken()
	token2, _, _ := svc.GenerateRefreshToken()
	assert.NotEqual(t, token1, token2)
}
