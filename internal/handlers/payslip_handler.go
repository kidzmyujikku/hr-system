package handlers

import (
	"errors"
	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"hr-system/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListPayslip(c *gin.Context) {
	userID, ip, reqID := middleware.ExtractData(c)

	payslip, err := services.ListPayslip(ip, reqID, userID)
	if err != nil {
		utils.LogHandlerError("ListPayslip", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrFailedListPayslip):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	utils.LogHandler("Payslip", userID, ip, reqID, "✅ Successfully list payslip")
	c.JSON(http.StatusOK, gin.H{
		"payslip": payslip,
	})
}
