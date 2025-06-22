package handlers

import (
	"errors"
	"hr-system/internal/dto"
	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"hr-system/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Summary(c *gin.Context) {
	var req dto.PayrollRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ip, reqID := middleware.ExtractData(c)

	summary, err := services.SummaryEmployee(req.PayCycleId, ip, reqID, userID)
	if err != nil {
		utils.LogHandlerError("Summary", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrFailedListSummary):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	utils.LogHandler("Summary", userID, ip, reqID, "✅ Successfully list summary")
	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}
