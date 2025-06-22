package handlers

import (
	"errors"
	"hr-system/internal/dto"
	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"net/http"

	"hr-system/internal/utils"

	"github.com/gin-gonic/gin"
)

func SubmitOvertime(c *gin.Context) {
	var req dto.SubmitOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ip, reqID := middleware.ExtractData(c)

	message, err := services.SubmitOvertime(userID, req, ip, reqID)
	if err != nil {
		utils.LogHandlerError("SubmitOvertime", userID, ip, reqID, err)
		switch {
		case errors.Is(err, services.ErrWeekDayWorkNotFound):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrOvertimeExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	utils.LogHandler("SubmitOvertime", userID, ip, reqID, "✅ overtime submitted")
	c.JSON(http.StatusOK, gin.H{"message": message})
}
