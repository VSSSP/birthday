package handler

import (
	"net/http"

	"github.com/vsssp/birthday-app/backend/internal/pkg/response"
	"github.com/vsssp/birthday-app/backend/internal/port"
)

// UserHandler handles user profile HTTP requests.
type UserHandler struct {
	userService port.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService port.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetCurrentUser handles GET /api/auth/me.
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}
	response.JSON(w, http.StatusOK, user)
}
