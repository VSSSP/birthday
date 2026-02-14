package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vsssp/birthday-app/backend/internal/domain"
	"github.com/vsssp/birthday-app/backend/internal/pkg/response"
	"github.com/vsssp/birthday-app/backend/internal/port"
	"github.com/vsssp/birthday-app/backend/internal/usecase"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	authService port.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}
	if len(req.Password) < 8 {
		response.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	tokens, err := h.authService.Register(r.Context(), req)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, tokens)
}

// Login handles POST /api/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	tokens, err := h.authService.Login(r.Context(), req)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, tokens)
}

// GoogleLogin handles POST /api/auth/google.
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var req domain.SocialLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.IDToken == "" {
		response.Error(w, http.StatusBadRequest, "id_token is required")
		return
	}

	tokens, err := h.authService.GoogleLogin(r.Context(), req.IDToken)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, tokens)
}

// AppleLogin handles POST /api/auth/apple.
func (h *AuthHandler) AppleLogin(w http.ResponseWriter, r *http.Request) {
	var req domain.SocialLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.IDToken == "" {
		response.Error(w, http.StatusBadRequest, "id_token is required")
		return
	}

	tokens, err := h.authService.AppleLogin(r.Context(), req.IDToken)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, tokens)
}

// RefreshToken handles POST /api/auth/refresh.
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, tokens)
}

func handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecase.ErrEmailAlreadyExists):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, usecase.ErrInvalidCredentials):
		response.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, usecase.ErrInvalidToken):
		response.Error(w, http.StatusUnauthorized, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
