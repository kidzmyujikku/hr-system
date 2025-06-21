package handlers

import (
	"hr-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	jwt := middleware.JwtMiddleware()

	r.POST("/login", jwt.LoginHandler)
}
