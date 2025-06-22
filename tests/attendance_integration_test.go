package tests

import (
	"fmt"
	"hr-system/internal/handlers"
	"hr-system/internal/middleware"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubmitAttendance_CheckInAndOut(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	randomUsername := fmt.Sprintf("employee%d", rand.Intn(100)+1)

	router := setupRouter()
	jwt := middleware.JwtMiddleware()
	router.POST("/employee/attendance", jwt.MiddlewareFunc(), handlers.SubmitAttendance)
	token := loginAndGetToken(t, randomUsername, "password123") // implement this

	// First call: should check in
	req, _ := http.NewRequest("POST", "/employee/attendance", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "check-in recorded")

	// Second call: should check out
	req, _ = http.NewRequest("POST", "/employee/attendance", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "check-out recorded")

	// Third call: already checked out
	req, _ = http.NewRequest("POST", "/employee/attendance", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "already checked out")
}
