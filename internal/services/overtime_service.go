package services

import (
	"errors"
	"hr-system/config"
	"hr-system/internal/dto"
	"hr-system/internal/models"
	"time"

	"gorm.io/gorm"
)

func ValidateDateFormat(err error) error {
	if err != nil {
		return ErrInvalidDateFormat
	}
	return nil
}

var (
	ErrWeekDayWorkNotFound = errors.New("overtime only can takes after checkout")
	ErrOvertimeExists      = errors.New("overtime already exists")
)

func SubmitOvertime(userID uint, req dto.SubmitOvertimeRequest, ip, reqID string) (string, error) {
	date, err := time.Parse(layoutISO, req.Date)
	if err := ValidateDateFormat(err); err != nil {
		return "", err
	}

	var count int64
	config.DB.Model(&models.Overtime{}).
		Where("user_id = ? AND date = ?", userID, date).
		Count(&count)

	if count > 0 {
		return "", ErrOvertimeExists
	}

	var attendance models.Attendance

	weekday := date.Weekday()
	if weekday != time.Saturday && weekday != time.Sunday {
		err := config.DB.
			Where("user_id = ? AND date = ?", userID, req.Date).
			First(&attendance).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrWeekDayWorkNotFound
		}

		if err != nil {
			return "", err // unexpected DB error
		}
	}

	overtime := models.Overtime{
		UserID:    userID,
		Date:      date,
		Hours:     req.Hours,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&overtime).Error; err != nil {
		return "", err
	}

	CreateAuditLog(userID, "overtime", overtime.ID, "CREATE", "Created overtime", ip, reqID)

	return "overtime created successfully", nil
}
