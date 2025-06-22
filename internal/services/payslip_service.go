package services

import (
	"errors"
	"hr-system/config"
	"hr-system/internal/models"
)

var (
	ErrFailedListPayslip = errors.New("failed to list payslip")
)

func ListPayslip(ip, reqId string, userID uint) ([]models.Payslip, error) {
	var payslip []models.Payslip
	result := config.DB.Where("employee_id = ?", userID).Find(&payslip)

	if result.Error != nil {
		return nil, ErrFailedListSummary
	}

	return payslip, nil
}
