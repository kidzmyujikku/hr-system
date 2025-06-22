package utils

import (
	log "github.com/sirupsen/logrus"
)

func LogHandlerError(action string, userID uint, ip, reqID string, err error) {
	log.WithFields(log.Fields{
		"action": action,
		"userID": userID,
		"ip":     ip,
		"reqID":  reqID,
		"error":  err.Error(),
	}).Error("Handler error")
}

func LogHandler(action string, userID uint, ip, reqID, message string) {
	log.WithFields(log.Fields{
		"action": action,
		"userID": userID,
		"ip":     ip,
		"reqID":  reqID,
	}).Info(message)
}
