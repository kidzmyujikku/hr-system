package tests

import (
	"bytes"
	"encoding/json"
	"hr-system/config"
	"hr-system/internal/middleware"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	jwt := middleware.JwtMiddleware()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/login", jwt.LoginHandler)
	return r
}

func TestLogin_Success(t *testing.T) {
	r := setupTestRouter()

	payload := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLogin_Failure(t *testing.T) {
	r := setupTestRouter()

	payload := map[string]string{
		"username": "admin",
		"password": "wrong",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "incorrect Username or Password")
}
