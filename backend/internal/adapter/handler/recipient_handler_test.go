package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to register and get access token
func registerAndGetToken(t *testing.T, router http.Handler, email string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": "password123",
		"name":     "Test User",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	return resp["access_token"].(string)
}

func TestCreateRecipient_Success(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "create-rec@example.com")

	body, _ := json.Marshal(map[string]interface{}{
		"name":       "Maria",
		"age":        65,
		"gender":     "female",
		"min_budget": 50.0,
		"max_budget": 200.0,
		"keywords":   []string{"cooking", "reading"},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/recipients", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Maria", resp["name"])
	assert.Equal(t, float64(65), resp["age"])
	assert.NotEmpty(t, resp["id"])
}

func TestCreateRecipient_MissingName(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "no-name@example.com")

	body, _ := json.Marshal(map[string]interface{}{
		"age":    30,
		"gender": "male",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/recipients", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateRecipient_Unauthorized(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)

	body, _ := json.Marshal(map[string]interface{}{"name": "Test"})
	req := httptest.NewRequest(http.MethodPost, "/api/recipients", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListRecipients_Empty(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "list-empty@example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/recipients", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Empty(t, resp)
}

func TestListRecipients_WithData(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "list-data@example.com")

	// Create 2 recipients
	for _, name := range []string{"Alice", "Bob"} {
		body, _ := json.Marshal(map[string]interface{}{
			"name":       name,
			"age":        30,
			"gender":     "other",
			"min_budget": 10,
			"max_budget": 100,
			"keywords":   []string{"test"},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/recipients", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	// List
	req := httptest.NewRequest(http.MethodGet, "/api/recipients", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Len(t, resp, 2)
}

func TestUpdateRecipient_Success(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "update@example.com")

	// Create
	createBody, _ := json.Marshal(map[string]interface{}{
		"name":       "Original",
		"age":        25,
		"gender":     "male",
		"min_budget": 10,
		"max_budget": 50,
		"keywords":   []string{"gaming"},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/recipients", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	json.NewDecoder(w.Body).Decode(&created)
	recipientID := created["id"].(string)

	// Update
	updateBody, _ := json.Marshal(map[string]interface{}{
		"name": "Updated",
		"age":  26,
	})
	req = httptest.NewRequest(http.MethodPut, "/api/recipients/"+recipientID, bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updated map[string]interface{}
	json.NewDecoder(w.Body).Decode(&updated)
	assert.Equal(t, "Updated", updated["name"])
	assert.Equal(t, float64(26), updated["age"])
}

func TestDeleteRecipient_Success(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "delete@example.com")

	// Create
	createBody, _ := json.Marshal(map[string]interface{}{
		"name":       "ToDelete",
		"age":        40,
		"gender":     "female",
		"min_budget": 20,
		"max_budget": 80,
		"keywords":   []string{},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/recipients", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	json.NewDecoder(w.Body).Decode(&created)
	recipientID := created["id"].(string)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/api/recipients/"+recipientID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify it's gone
	req = httptest.NewRequest(http.MethodGet, "/api/recipients/"+recipientID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteRecipients_Success(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "bulk-del@example.com")

	// Create 3 recipients
	var ids []string
	for _, name := range []string{"A", "B", "C"} {
		body, _ := json.Marshal(map[string]interface{}{
			"name":       name,
			"age":        20,
			"gender":     "other",
			"min_budget": 10,
			"max_budget": 50,
			"keywords":   []string{},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/recipients", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var created map[string]interface{}
		json.NewDecoder(w.Body).Decode(&created)
		ids = append(ids, created["id"].(string))
	}

	// Bulk delete first 2
	deleteBody, _ := json.Marshal(map[string]interface{}{
		"ids": ids[:2],
	})
	req := httptest.NewRequest(http.MethodDelete, "/api/recipients", bytes.NewReader(deleteBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRecipient_NotFound(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "notfound@example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/recipients/00000000-0000-0000-0000-000000000099", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetRecipient_InvalidID(t *testing.T) {
	router, _, _, _, _, _ := setupRouter(t)
	token := registerAndGetToken(t, router, "invalid-id@example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/recipients/not-a-uuid", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
