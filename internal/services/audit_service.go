package services

import (
	"hr-system/config"
	"hr-system/internal/models"
	"time"

	log "github.com/sirupsen/logrus"
)

func CreateAuditLog(userID uint, entity string, entityId uint, action, description, ip, requestID string) {
	go func() {
		logdata := models.AuditLog{
			UserID:      userID,
			Entity:      entity,
			EntityID:    entityId,
			Action:      action,
			Description: description,
			IPAddress:   ip,
			RequestID:   requestID,
			CreatedAt:   time.Now(),
		}

		// Use a separate DB instance or connection if needed

		if err := config.DB.Create(&logdata).Error; err != nil {
			log.Error("Failed to write audit log: ", err)
		}
	}()
}
