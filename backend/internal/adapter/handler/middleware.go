package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	jwtpkg "github.com/vsssp/birthday-app/backend/internal/pkg/jwt"
	"github.com/vsssp/birthday-app/backend/internal/pkg/response"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
)

// AuthMiddleware validates JWT tokens on protected routes.
type AuthMiddleware struct {
	jwtService *jwtpkg.Service
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(jwtService *jwtpkg.Service) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// Authenticate is the middleware handler that validates Bearer tokens.
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.Error(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext extracts the user ID from the request context.
func UserIDFromContext(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(UserIDKey).(uuid.UUID)
	return id
}
