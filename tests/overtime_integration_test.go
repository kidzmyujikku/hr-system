package tests

import (
	"bytes"
	"encoding/json"
	"hr-system/internal/handlers"
	"hr-system/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmitOvertime_HoursMoreThan3(t *testing.T) {
	router := setupRouter()
	jwt := middleware.JwtMiddleware()
	router.POST("/employee/overtime", jwt.MiddlewareFunc(), handlers.SubmitOvertime)
	token := loginAndGetToken(t, "employee6", "password123") // implement this

	payload := map[string]interface{}{
		"date":  "2025-06-22",
		"hours": 3.5,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/employee/overtime", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation for 'Hours' failed on the 'lte' tag")
}
