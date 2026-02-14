package response

import (
	"encoding/json"
	"net/http"
)

// errorResponse is the standard error payload.
type errorResponse struct {
	Error string `json:"error"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error writes a JSON error response.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, errorResponse{Error: message})
}
