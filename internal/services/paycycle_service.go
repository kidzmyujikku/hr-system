package services

import (
	"errors"
	"time"

	"hr-system/config"
	"hr-system/internal/models"

	"gorm.io/gorm"
)

const layoutISO = "2006-01-02"

var (
	ErrFailedListPayCycle  = errors.New("failed to list pay cycles")
	ErrInvalidDateFormat   = errors.New("invalid date format")
	ErrDateRange           = errors.New("end date must be after start date")
	ErrPayCycleExists      = errors.New("a pay cycle already exists in this range")
	ErrPayCycleNotFound    = errors.New("pay cycle not found")
	ErrPayCycleIsProcessed = errors.New("cannot delete a processed pay cycle")
	ErrFailedDelete        = errors.New("failed to delete pay cycle")
)

func ValidateFormat(err1, err2 error) error {
	if err1 != nil || err2 != nil {
		return ErrInvalidDateFormat
	}
	return nil
}

func ValidateCycle(startDate, endDate time.Time) error {
	if endDate.Before(startDate) {
		return ErrDateRange
	}
	return nil
}

func ListPayCycles(limit, offset int) (int64, []models.PayCycle, error) {
	var cycles []models.PayCycle
	var total int64

	config.DB.Model(&models.PayCycle{}).Count(&total)

	result := config.DB.
		Preload("Creator").
		Preload("Updater").
		Limit(limit).
		Offset(offset).
		Order("start_date desc").
		Find(&cycles)

	if result.Error != nil {
		return 0, nil, ErrFailedListPayCycle
	}

	return total, cycles, nil
}

func CreatePayCycle(startStr, endStr, ip, reqID string, userID uint) (string, error) {
	startDate, err1 := time.Parse(layoutISO, startStr)
	endDate, err2 := time.Parse(layoutISO, endStr)

	if err := ValidateFormat(err1, err2); err != nil {
		return "", err
	}
	if err := ValidateCycle(startDate, endDate); err != nil {
		return "", err
	}

	var count int64
	config.DB.Model(&models.PayCycle{}).
		Where("start_date <= ? AND end_date >= ?", endDate, startDate).
		Count(&count)

	if count > 0 {
		return "", ErrPayCycleExists
	}

	period := models.PayCycle{
		StartDate:   startDate,
		EndDate:     endDate,
		MonthLabel:  startDate.Format("January 2006"),
		IsProcessed: false,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&period).Error; err != nil {
		return "", err
	}

	CreateAuditLog(userID, "pay_cycle", period.ID, "CREATE", "Created pay cycle", ip, reqID)

	return "pay cycle created successfully", nil
}

func UpdatePayCycle(id, startStr, endStr, ip, reqID string, userID uint) (string, error) {
	var cycle models.PayCycle
	err := config.DB.First(&cycle, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrPayCycleNotFound
		}
		return "", err
	}

	startDate, err1 := time.Parse(layoutISO, startStr)
	endDate, err2 := time.Parse(layoutISO, endStr)

	if err := ValidateFormat(err1, err2); err != nil {
		return "", err
	}
	if err := ValidateCycle(startDate, endDate); err != nil {
		return "", err
	}

	cycle.StartDate = startDate
	cycle.EndDate = endDate
	cycle.MonthLabel = startDate.Format("January 2006")
	cycle.UpdatedBy = &userID

	if err := config.DB.Save(&cycle).Error; err != nil {
		return "", err
	}

	CreateAuditLog(userID, "pay_cycle", cycle.ID, "UPDATE", "Update pay cycle", ip, reqID)

	return "pay cycle updated", nil
}

func DeletePayCycle(id, ip, reqID string, userID uint) (string, error) {
	var cycle models.PayCycle
	if err := config.DB.First(&cycle, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrPayCycleNotFound
		}
		return "", err
	}

	if cycle.IsProcessed {
		return "", ErrPayCycleIsProcessed
	}

	if err := config.DB.Delete(&cycle).Error; err != nil {
		return "", ErrFailedDelete
	}

	CreateAuditLog(userID, "pay_cycle", cycle.ID, "DELETE", "Delete pay cycle", ip, reqID)

	return "pay cycle deleted", nil
}
