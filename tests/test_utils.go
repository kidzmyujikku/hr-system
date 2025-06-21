package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hr-system/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func loginAndGetAdminToken(t *testing.T) string {
	// Optional: set up a minimal router just for login
	r := gin.Default()
	jwt := middleware.JwtMiddleware()
	r.POST("/login", jwt.LoginHandler)

	// Login payload
	payload := map[string]string{
		"username": "admin",
		"password": "admin123",
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
