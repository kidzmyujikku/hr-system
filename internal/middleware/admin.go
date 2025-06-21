package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		if role, ok := claims["role"].(string); !ok || role != "admin" {
			c.AbortWithStatusJSON(403, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}
