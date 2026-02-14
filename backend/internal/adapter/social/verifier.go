package social

import "context"

// CompositeVerifier combines Google and Apple verifiers into a single interface.
type CompositeVerifier struct {
	google *GoogleVerifier
	apple  *AppleVerifier
}

// NewCompositeVerifier creates a new CompositeVerifier.
func NewCompositeVerifier(google *GoogleVerifier, apple *AppleVerifier) *CompositeVerifier {
	return &CompositeVerifier{google: google, apple: apple}
}

// VerifyGoogleToken delegates to the Google verifier.
func (v *CompositeVerifier) VerifyGoogleToken(ctx context.Context, idToken string) (email, name, sub string, err error) {
	return v.google.VerifyGoogleToken(ctx, idToken)
}

// VerifyAppleToken delegates to the Apple verifier.
func (v *CompositeVerifier) VerifyAppleToken(ctx context.Context, identityToken string) (email, sub string, err error) {
	return v.apple.VerifyAppleToken(ctx, identityToken)
}
