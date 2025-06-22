package handlers

import (
	"hr-system/internal/dto"
	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"hr-system/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunPayroll(c *gin.Context) {
	var req dto.PayrollRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	userID, ip, reqID := middleware.ExtractData(c)

	message, err := services.RunPayroll(req.PayCycleId, ip, reqID, userID)
	if err != nil {
		utils.LogHandlerError("RunPayroll", userID, ip, reqID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.LogHandler("RunPayroll", userID, ip, reqID, "✅ Successfully run payroll")
	c.JSON(http.StatusOK, gin.H{"message": message})
}
