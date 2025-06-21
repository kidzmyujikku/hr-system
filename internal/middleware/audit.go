package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
)

func ExtractData(c *gin.Context) (userID uint, ip, reqID string) {
	claims := jwt.ExtractClaims(c)
	userID = uint(claims["userID"].(float64))

	ip = c.ClientIP()
	reqID = uuid.New().String()

	return
}
