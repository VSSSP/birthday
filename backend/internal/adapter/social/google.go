package social

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"
)

// GoogleVerifier validates Google ID tokens.
type GoogleVerifier struct {
	clientID string
}

// NewGoogleVerifier creates a new GoogleVerifier.
func NewGoogleVerifier(clientID string) *GoogleVerifier {
	return &GoogleVerifier{clientID: clientID}
}

// VerifyGoogleToken validates a Google ID token and extracts user info.
func (v *GoogleVerifier) VerifyGoogleToken(ctx context.Context, idToken string) (email, name, sub string, err error) {
	payload, err := idtoken.Validate(ctx, idToken, v.clientID)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid google id token: %w", err)
	}

	email, _ = payload.Claims["email"].(string)
	name, _ = payload.Claims["name"].(string)
	sub = payload.Subject

	if email == "" || sub == "" {
		return "", "", "", fmt.Errorf("missing email or subject in google token")
	}
	return email, name, sub, nil
}
