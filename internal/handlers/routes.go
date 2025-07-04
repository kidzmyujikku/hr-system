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

		admin.POST("/payroll/run", RunPayroll)

		admin.POST("/summary", Summary)
	}

	emp := auth.Group("/employee")
	{
		emp.POST("/attendance", SubmitAttendance)
		emp.POST("/overtime", SubmitOvertime)

		emp.POST("/reimbursement", CreateReimbursement)
		emp.GET("/reimbursement", ListReimbursement)
		emp.PUT("/reimbursement/:id", UpdateReimbursement)
		emp.DELETE("/reimbursement/:id", DeleteReimbursement)

		emp.GET("/payslip", ListPayslip)
	}

}
