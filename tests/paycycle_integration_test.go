package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"hr-system/config"
	"hr-system/internal/handlers"
	"hr-system/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupRouterForPaycycle() *gin.Engine {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	r := gin.Default()
	adminAuth := middleware.JwtMiddleware()
	r.POST("/admin/paycycle", adminAuth.MiddlewareFunc(), handlers.CreatePayCycle)
	return r
}

func TestCreatePayCycle_InvalidDateFormat(t *testing.T) {
	router := setupRouterForPaycycle()
	token := loginAndGetAdminToken(t) // implement this

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
	router := setupRouterForPaycycle()
	token := loginAndGetAdminToken(t)

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
