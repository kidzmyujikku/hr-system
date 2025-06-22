package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"hr-system/config"
	"hr-system/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func loginAndGetToken(t *testing.T, username, password string) string {
	// Optional: set up a minimal router just for login
	r := gin.Default()
	jwt := middleware.JwtMiddleware()
	r.POST("/login", jwt.LoginHandler)

	// Login payload
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Check login success
	require.Equal(t, http.StatusOK, w.Code)

	// Parse token from response
	var resBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resBody)

	token, ok := resBody["token"].(string)
	require.True(t, ok, "Token not found in login response")

	return token
}

func setupRouter() *gin.Engine {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	r := gin.Default()
	return r
}
