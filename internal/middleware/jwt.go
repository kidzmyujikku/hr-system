package middleware

import (
	"log"
	"os"
	"time"

	"hr-system/config"
	"hr-system/internal/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var IdentityKey = "userID"

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// JWT middleware
func JwtMiddleware() *jwt.GinJWTMiddleware {

	jwtKey := os.Getenv("JWT_SECRET")

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "hr zone",
		Key:         []byte(jwtKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login Login
			if err := c.ShouldBindJSON(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			var user models.User
			if err := config.DB.Where("username = ?", login.Username).First(&user).Error; err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &user, nil
		},

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					IdentityKey: user.ID,
					"username":  user.Username,
					"role":      user.Role,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims[IdentityKey]
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true // allow all authenticated users
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:", err)
	}

	return authMiddleware
}
