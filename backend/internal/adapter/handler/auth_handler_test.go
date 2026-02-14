package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vsssp/birthday-app/backend/internal/adapter/handler"
	jwtpkg "github.com/vsssp/birthday-app/backend/internal/pkg/jwt"
	"github.com/vsssp/birthday-app/backend/internal/usecase"
)

func setupRouter(t *testing.T) (*http.ServeMux, *mockUserRepo, *mockAuthProviderRepo, *mockRefreshTokenRepo, *mockRecipientRepo, *jwtpkg.Service) {
	t.Helper()

	userRepo := newMockUserRepo()
	providerRepo := newMockAuthProviderRepo()
	tokenRepo := newMockRefreshTokenRepo()
	recipientRepo := newMockRecipientRepo()

	jwtService := jwtpkg.NewService(
		"test-access-secret-32-chars-long!",
		"test-refresh-secret-32-chars-lo!",
		15*time.Minute,
		7*24*time.Hour,
	)

	socialVerifier := &mockSocialVerifier{}
	authUseCase := usecase.NewAuthUseCase(userRepo, providerRepo, tokenRepo, jwtService, socialVerifier)
	userUseCase := usecase.NewUserUseCase(userRepo)
	recipientUseCase := usecase.NewRecipientUseCase(recipientRepo)

	router := handler.NewRouter(authUseCase, userUseCase, recipientUseCase, jwtService)

	mux := http.NewServeMux()
	mux.Handle("/", router)

	return mux, userRepo, providerRepo, tokenRepo, recipientRepo, jwtService
}

func TestRegister_Success(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["access_token"])
	assert.NotEmpty(t, resp["refresh_token"])
	assert.NotNil(t, resp["expires_at"])
}

func TestRegister_DuplicateEmail(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	jsonBody, _ := json.Marshal(body)

	// First registration
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Second registration with same email
	jsonBody, _ = json.Marshal(body)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_MissingFields(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	body := map[string]string{"email": "test@example.com"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_ShortPassword(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "short",
		"name":     "Test",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	// Register first
	regBody, _ := json.Marshal(map[string]string{
		"email":    "login@example.com",
		"password": "password123",
		"name":     "Login User",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "login@example.com",
		"password": "password123",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.NotEmpty(t, resp["access_token"])
}

func TestLogin_WrongPassword(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	// Register
	regBody, _ := json.Marshal(map[string]string{
		"email":    "wrong@example.com",
		"password": "password123",
		"name":     "User",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Login with wrong password
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "wrong@example.com",
		"password": "wrongpassword",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_NonexistentUser(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	body, _ := json.Marshal(map[string]string{
		"email":    "noone@example.com",
		"password": "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshToken_Success(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	// Register to get tokens
	regBody, _ := json.Marshal(map[string]string{
		"email":    "refresh@example.com",
		"password": "password123",
		"name":     "Refresh User",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var regResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&regResp)
	refreshToken := regResp["refresh_token"].(string)

	// Refresh
	refreshBody, _ := json.Marshal(map[string]string{
		"refresh_token": refreshToken,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.NotEmpty(t, resp["access_token"])
	assert.NotEqual(t, refreshToken, resp["refresh_token"]) // Token rotation
}

func TestRefreshToken_Invalid(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	body, _ := json.Marshal(map[string]string{
		"refresh_token": "invalid-token",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetMe_Success(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	// Register
	regBody, _ := json.Marshal(map[string]string{
		"email":    "me@example.com",
		"password": "password123",
		"name":     "Me User",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var regResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&regResp)
	accessToken := regResp["access_token"].(string)

	// Get me
	req = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var user map[string]interface{}
	json.NewDecoder(w.Body).Decode(&user)
	assert.Equal(t, "me@example.com", user["email"])
	assert.Equal(t, "Me User", user["name"])
}

func TestGetMe_NoToken(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetMe_InvalidToken(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHealthCheck(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}
