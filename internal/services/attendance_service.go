package services

import (
	"errors"
	"hr-system/config"
	"hr-system/internal/models"
	"time"

	"gorm.io/gorm"
)

var (
	ErrWeekendSubmission = errors.New("cannot submit attendance on weekends")
	ErrAlreadyCheckedOut = errors.New("already checked out")
)

func SubmitAttendance(userID uint, ip, reqID string) (string, error) {
	loc, err := time.LoadLocation("Asia/Jakarta") // Adjust based on your app's timezone
	if err != nil {
		return "", err
	}

	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// // Reject weekend submissions
	// weekday := now.Weekday()
	// if weekday == time.Saturday || weekday == time.Sunday {
	// 	return "", ErrWeekendSubmission
	// }

	var attendance models.Attendance
	err = config.DB.
		Where("user_id = ? AND date = ?", userID, today).
		First(&attendance).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Check-in
		newAttendance := models.Attendance{
			UserID:    userID,
			Date:      today,
			CheckIn:   now,
			CreatedAt: now,
			CreatedBy: userID,
		}
		if err := config.DB.Create(&newAttendance).Error; err != nil {

			return "", err
		}

		CreateAuditLog(userID, "attendance", newAttendance.ID, "CREATE", "Create check-in", ip, reqID)

		return "check-in recorded", nil
	}

	if err != nil {
		return "", err // unexpected DB error
	}

	if attendance.CheckOut != nil {
		return "", ErrAlreadyCheckedOut
	}

	// Check-out
	attendance.CheckOut = &now
	attendance.UpdatedAt = now

	if err := config.DB.Save(&attendance).Error; err != nil {
		return "", err
	}

	CreateAuditLog(userID, "attendance", attendance.ID, "UPDATE", "Add check-out", ip, reqID)
	return "check-out recorded", nil
}
