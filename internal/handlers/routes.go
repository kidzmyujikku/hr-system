package handlers

import (
	"hr-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	jwt := middleware.JwtMiddleware()

	r.POST("/login", jwt.LoginHandler)

	auth := r.Group("/")
	auth.Use(jwt.MiddlewareFunc())

	admin := auth.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		admin.POST("/pay-cycle", CreatePayCycle)
		admin.GET("/pay-cycle", ListPayCycles)
		admin.PUT("/pay-cycle/:id", UpdatePayCycle)
		admin.DELETE("/pay-cycle/:id", DeletePayCycle)
	}

	emp := auth.Group("/employee")
	{
		emp.POST("/attendance", SubmitAttendance)
		emp.POST("/overtime", SubmitOvertime)
	}

}
