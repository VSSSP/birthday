package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
	"github.com/vsssp/birthday-app/backend/internal/pkg/response"
	"github.com/vsssp/birthday-app/backend/internal/port"
	"github.com/vsssp/birthday-app/backend/internal/usecase"
)

// RecipientHandler handles recipient CRUD HTTP requests.
type RecipientHandler struct {
	recipientService port.RecipientService
}

// NewRecipientHandler creates a new RecipientHandler.
func NewRecipientHandler(recipientService port.RecipientService) *RecipientHandler {
	return &RecipientHandler{recipientService: recipientService}
}

// Create handles POST /api/recipients.
func (h *RecipientHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	var req domain.CreateRecipientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	recipient, err := h.recipientService.Create(r.Context(), userID, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create recipient")
		return
	}
	response.JSON(w, http.StatusCreated, recipient)
}

// List handles GET /api/recipients.
func (h *RecipientHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	recipients, err := h.recipientService.ListByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list recipients")
		return
	}
	response.JSON(w, http.StatusOK, recipients)
}

// GetByID handles GET /api/recipients/{id}.
func (h *RecipientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	recipientID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid recipient id")
		return
	}

	recipient, err := h.recipientService.GetByID(r.Context(), userID, recipientID)
	if err != nil {
		handleRecipientError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, recipient)
}

// Update handles PUT /api/recipients/{id}.
func (h *RecipientHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	recipientID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid recipient id")
		return
	}

	var req domain.UpdateRecipientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	recipient, err := h.recipientService.Update(r.Context(), userID, recipientID, req)
	if err != nil {
		handleRecipientError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, recipient)
}

// Delete handles DELETE /api/recipients/{id}.
func (h *RecipientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	recipientID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid recipient id")
		return
	}

	if err := h.recipientService.Delete(r.Context(), userID, recipientID); err != nil {
		handleRecipientError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "recipient deleted"})
}

// BulkDelete handles DELETE /api/recipients.
func (h *RecipientHandler) BulkDelete(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	var req domain.BulkDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.IDs) == 0 {
		response.Error(w, http.StatusBadRequest, "ids are required")
		return
	}

	if err := h.recipientService.BulkDelete(r.Context(), userID, req.IDs); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete recipients")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "recipients deleted"})
}

func handleRecipientError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecase.ErrRecipientNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, usecase.ErrForbidden):
		response.Error(w, http.StatusForbidden, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
