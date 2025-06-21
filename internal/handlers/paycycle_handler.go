package handlers

import (
	"errors"
	"hr-system/internal/dto"
	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Utility for logging errors consistently
func logHandlerError(action string, userID uint, ip, reqID string, err error) {
	log.WithFields(log.Fields{
		"action": action,
		"userID": userID,
		"ip":     ip,
		"reqID":  reqID,
		"error":  err.Error(),
	}).Error("Handler error")
}

func CreatePayCycle(c *gin.Context) {
	var req dto.PayCycleRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithFields(log.Fields{
			"action": "CreatePayCycle",
			"ip":     c.ClientIP(),
			"reqID":  c.GetString("reqID"),
			"error":  err.Error(),
		}).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ip, reqID := middleware.ExtractData(c)

	message, err := services.CreatePayCycle(req.StartDate, req.EndDate, ip, reqID, userID)
	if err != nil {
		logHandlerError("CreatePayCycle", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrInvalidDateFormat),
			errors.Is(err, services.ErrDateRange):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPayCycleExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.WithFields(log.Fields{
		"action": "CreatePayCycle",
		"userID": userID,
		"ip":     ip,
		"reqID":  reqID,
	}).Info("✅ Pay cycle created")
	c.JSON(http.StatusCreated, gin.H{"message": message})
}

func ListPayCycles(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	offset := (page - 1) * limit

	userID, ip, reqID := middleware.ExtractData(c)

	total, cycles, err := services.ListPayCycles(limit, offset)
	if err != nil {
		logHandlerError("ListPayCycles", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrInvalidDateFormat),
			errors.Is(err, services.ErrDateRange):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrFailedListPayCycle):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.WithFields(log.Fields{
		"action":   "ListPayCycles",
		"userID":   userID,
		"ip":       ip,
		"reqID":    reqID,
		"page":     page,
		"limit":    limit,
		"total":    total,
		"returned": len(cycles),
	}).Info("📄 Pay cycles listed")

	c.JSON(http.StatusOK, gin.H{
		"page":   page,
		"limit":  limit,
		"total":  total,
		"cycles": cycles,
	})
}

func UpdatePayCycle(c *gin.Context) {
	id := c.Param("id")

	var req dto.PayCycleRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithFields(log.Fields{
			"action": "UpdatePayCycle",
			"ip":     c.ClientIP(),
			"reqID":  c.GetString("reqID"),
			"error":  err.Error(),
		}).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ip, reqID := middleware.ExtractData(c)

	message, err := services.UpdatePayCycle(id, req.StartDate, req.EndDate, ip, reqID, userID)
	if err != nil {
		logHandlerError("UpdatePayCycle", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrPayCycleNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.WithFields(log.Fields{
		"action": "UpdatePayCycle",
		"userID": userID,
		"ip":     ip,
		"reqID":  reqID,
	}).Info("✅ Pay cycle updated")

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func DeletePayCycle(c *gin.Context) {
	id := c.Param("id")
	userID, ip, reqID := middleware.ExtractData(c)

	message, err := services.DeletePayCycle(id, ip, reqID, userID)
	if err != nil {
		logHandlerError("DeletePayCycle", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrPayCycleNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPayCycleIsProcessed):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrFailedDelete):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.WithFields(log.Fields{
		"action": "DeletePayCycle",
		"userID": userID,
		"ip":     ip,
		"reqID":  reqID,
	}).Info("✅ Pay cycle deleted")

	c.JSON(http.StatusOK, gin.H{"message": message})
}
