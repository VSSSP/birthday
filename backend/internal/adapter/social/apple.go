package social

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const appleKeysURL = "https://appleid.apple.com/auth/keys"

// AppleVerifier validates Apple identity tokens.
type AppleVerifier struct {
	clientID string
}

// NewAppleVerifier creates a new AppleVerifier.
func NewAppleVerifier(clientID string) *AppleVerifier {
	return &AppleVerifier{clientID: clientID}
}

type appleKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type appleKeysResponse struct {
	Keys []appleKey `json:"keys"`
}

// VerifyAppleToken validates an Apple identity token and extracts user info.
func (v *AppleVerifier) VerifyAppleToken(ctx context.Context, identityToken string) (email, sub string, err error) {
	resp, err := http.Get(appleKeysURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch apple keys: %w", err)
	}
	defer resp.Body.Close()

	var keysResp appleKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&keysResp); err != nil {
		return "", "", fmt.Errorf("failed to decode apple keys: %w", err)
	}

	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(identityToken, jwt.MapClaims{})
	if err != nil {
		return "", "", fmt.Errorf("failed to parse apple token: %w", err)
	}

	kid, _ := token.Header["kid"].(string)

	var matchingKey *appleKey
	for _, k := range keysResp.Keys {
		if k.Kid == kid {
			matchingKey = &k
			break
		}
	}
	if matchingKey == nil {
		return "", "", fmt.Errorf("no matching apple key found for kid: %s", kid)
	}

	pubKey, err := jwkToRSAPublicKey(matchingKey)
	if err != nil {
		return "", "", err
	}

	claims := jwt.MapClaims{}
	verifiedToken, err := jwt.ParseWithClaims(identityToken, claims, func(t *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil || !verifiedToken.Valid {
		return "", "", fmt.Errorf("invalid apple identity token: %w", err)
	}

	aud, _ := claims["aud"].(string)
	if aud != v.clientID {
		return "", "", fmt.Errorf("apple token audience mismatch")
	}

	iss, _ := claims["iss"].(string)
	if iss != "https://appleid.apple.com" {
		return "", "", fmt.Errorf("apple token issuer mismatch")
	}

	email, _ = claims["email"].(string)
	sub, _ = claims["sub"].(string)

	if sub == "" {
		return "", "", fmt.Errorf("missing subject in apple token")
	}
	return email, sub, nil
}

func jwkToRSAPublicKey(key *appleKey) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}
