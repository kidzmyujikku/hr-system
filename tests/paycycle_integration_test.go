package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hr-system/internal/handlers"
	"hr-system/internal/middleware"

	"github.com/stretchr/testify/assert"
)

func TestCreatePayCycle_InvalidDateFormat(t *testing.T) {
	router := setupRouter()
	adminAuth := middleware.JwtMiddleware()
	router.POST("/admin/paycycle", adminAuth.MiddlewareFunc(), handlers.CreatePayCycle)
	token := loginAndGetToken(t, "admin", "admin123") // implement this

	payload := map[string]string{
		"start_date": "bad-date",
		"end_date":   "2025-11-30",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/admin/paycycle", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePayCycle_EndBeforeStart(t *testing.T) {
	router := setupRouter()
	adminAuth := middleware.JwtMiddleware()
	router.POST("/admin/paycycle", adminAuth.MiddlewareFunc(), handlers.CreatePayCycle)

	token := loginAndGetToken(t, "admin", "admin123")

	payload := map[string]string{
		"start_date": "2025-11-30",
		"end_date":   "2025-11-01",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/admin/paycycle", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "end date must be after start date")
}
