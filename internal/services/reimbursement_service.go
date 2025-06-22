package services

import (
	"errors"
	"hr-system/config"
	"hr-system/internal/dto"
	"hr-system/internal/models"
	"time"

	"gorm.io/gorm"
)

var (
	ErrReimbursementNotFound = errors.New("reimburse not found")
	ErrFailedListReimburse   = errors.New("failed to list reimburse")
	ErrFailedDeleteReimburse = errors.New("failed to delete reimburse")
)

func ListReimbursement(limit, offset int) (int64, []models.Reimbursement, error) {
	var reimburse []models.Reimbursement
	var total int64

	config.DB.Model(&models.Reimbursement{}).Count(&total)

	result := config.DB.
		Preload("Creator").
		Preload("Updater").
		Limit(limit).
		Offset(offset).
		Order("date desc").
		Find(&reimburse)

	if result.Error != nil {
		return 0, nil, ErrFailedListReimburse
	}

	return total, reimburse, nil
}

func CreateReimbursement(userID uint, req dto.SubmitReimbursementRequest, ip, reqID string) (string, error) {
	date, err := time.Parse(layoutISO, req.Date)
	if err := ValidateDateFormat(err); err != nil {
		return "", err
	}

	reimbursement := models.Reimbursement{
		UserID:      userID,
		Date:        date,
		Amount:      req.Amount,
		Description: req.Description,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&reimbursement).Error; err != nil {
		return "", err
	}

	CreateAuditLog(userID, "reimbursement", reimbursement.ID, "CREATE", "Created reimburse", ip, reqID)

	return "reimburse created successfully", nil
}

func UpdateReimbursement(userID uint, req dto.SubmitReimbursementRequest, id, ip, reqID string) (string, error) {
	var reimburse models.Reimbursement
	err := config.DB.First(&reimburse, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrReimbursementNotFound
		}
		return "", err
	}

	date, err := time.Parse(layoutISO, req.Date)
	if err := ValidateDateFormat(err); err != nil {
		return "", err
	}

	reimburse.Date = date
	reimburse.Amount = req.Amount
	reimburse.Description = req.Description
	reimburse.UpdatedBy = &userID

	if err := config.DB.Save(&reimburse).Error; err != nil {
		return "", err
	}

	CreateAuditLog(userID, "reimbursement", reimburse.ID, "UPDATE", "Update reimburse", ip, reqID)

	return "reimburse updated", nil
}

func DeleteReimbursement(id, ip, reqID string, userID uint) (string, error) {
	var reimburse models.Reimbursement
	if err := config.DB.First(&reimburse, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrReimbursementNotFound
		}
		return "", err
	}

	if err := config.DB.Delete(&reimburse).Error; err != nil {
		return "", ErrFailedDeleteReimburse
	}

	CreateAuditLog(userID, "reimbursement", reimburse.ID, "DELETE", "Delete reimburse", ip, reqID)

	return "reimburse deleted", nil
}
