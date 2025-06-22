package services

import (
	"errors"
	"hr-system/config"
	"hr-system/internal/models"
)

var (
	ErrFailedListSummary = errors.New("failed to list summary")
)

func SummaryEmployee(id uint, ip, reqId string, userID uint) ([]models.Payslip, error) {
	var payslip []models.Payslip
	result := config.DB.Where("pay_cycle_id = ?", id).Find(&payslip)

	if result.Error != nil {
		return nil, ErrFailedListSummary
	}

	return payslip, nil
}
