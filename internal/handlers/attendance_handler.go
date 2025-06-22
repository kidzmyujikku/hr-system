package handlers

import (
	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SubmitAttendance(c *gin.Context) {
	userID, ip, reqID := middleware.ExtractData(c)

	message, err := services.SubmitAttendance(userID, ip, reqID)
	if err != nil {
		logHandlerError("SubmitAttendance", userID, ip, reqID, err)
		switch {
		case errors.Is(err, services.ErrWeekendSubmission):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrAlreadyCheckedOut):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	log.WithFields(log.Fields{
		"action": "SubmitAttendance",
		"userID": userID,
		"ip":     ip,
		"reqID":  reqID,
	}).Info("✅ Attendance submitted")

	c.JSON(http.StatusOK, gin.H{"message": message})
}
