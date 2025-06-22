package handlers

import (
	"errors"
	"hr-system/internal/dto"
	"hr-system/internal/middleware"
	"hr-system/internal/services"
	"hr-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func CreateReimbursement(c *gin.Context) {
	userID, ip, reqID := middleware.ExtractData(c)

	var req dto.SubmitReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogHandlerError("CreateReimbursement", userID, ip, reqID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := services.CreateReimbursement(userID, req, ip, reqID)
	if err != nil {
		utils.LogHandlerError("CreateReimbursement", userID, ip, reqID, err)
		switch {
		case errors.Is(err, services.ErrInvalidDateFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	utils.LogHandler("CreateReimbursement", userID, ip, reqID, "✅ Reimburse created")
	c.JSON(http.StatusCreated, gin.H{"message": message})
}

func ListReimbursement(c *gin.Context) {
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

	total, reimburse, err := services.ListReimbursement(limit, offset)
	if err != nil {
		utils.LogHandlerError("ListReimbursement", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrFailedListReimburse):
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
		"returned": len(reimburse),
	}).Info("📄 Reimburse listed")

	c.JSON(http.StatusOK, gin.H{
		"page":          page,
		"limit":         limit,
		"total":         total,
		"reimbursement": reimburse,
	})
}

func UpdateReimbursement(c *gin.Context) {
	id := c.Param("id")

	userID, ip, reqID := middleware.ExtractData(c)

	var req dto.SubmitReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogHandlerError("UpdateReimbursement", userID, ip, reqID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := services.UpdateReimbursement(userID, req, id, ip, reqID)
	if err != nil {
		utils.LogHandlerError("UpdateReimbursement", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrReimbursementNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	utils.LogHandler("UpdateReimbursement", userID, ip, reqID, "✅ Reimburse updated")

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func DeleteReimbursement(c *gin.Context) {
	id := c.Param("id")
	userID, ip, reqID := middleware.ExtractData(c)

	message, err := services.DeleteReimbursement(id, ip, reqID, userID)
	if err != nil {
		utils.LogHandlerError("DeleteReimbursement", userID, ip, reqID, err)

		switch {
		case errors.Is(err, services.ErrReimbursementNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrFailedDeleteReimburse):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	utils.LogHandler("DeleteReimbursement", userID, ip, reqID, "✅ Reimburse deleted")
	c.JSON(http.StatusOK, gin.H{"message": message})
}
