package services

import (
	"errors"
	"hr-system/config"
	"hr-system/internal/models"
	"time"

	"gorm.io/gorm"
)

func getWorkingDays(start, end time.Time) int {
	days := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday >= time.Monday && weekday <= time.Friday {
			days++
		}
	}
	return days
}

func calculateReimbursement(empID uint, start, end time.Time) int64 {
	var reimbursements []models.Reimbursement
	config.DB.Where("user_id = ? AND date BETWEEN ? AND ?", empID, start, end).Find(&reimbursements)

	var total int64
	for _, r := range reimbursements {
		total += r.Amount
	}
	return total
}

func calculateOvertime(empID uint, start, end time.Time) (int64, float32) {
	var employee models.User
	if err := config.DB.First(&employee, empID).Error; err != nil {
		return 0, 0
	}

	var overtimes []models.Overtime
	config.DB.Where("user_id = ? AND date BETWEEN ? AND ?", empID, start, end).Find(&overtimes)

	var totalHours float32
	for _, o := range overtimes {
		hours := o.Hours
		if hours > 3 {
			hours = 3
		}
		totalHours += hours
	}

	workingDays := getWorkingDays(start, end)
	if workingDays == 0 {
		return 0, 0
	}

	// Convert to float64 for safe math
	hourlyRate := float64(employee.Salary) / float64(workingDays*8)
	overtimePay := int64(hourlyRate * float64(totalHours) * 2) // 2x multiplier

	return overtimePay, totalHours
}

func calculateAttendance(empID uint, start, end time.Time) (int64, int64) {
	var employee models.User
	if err := config.DB.First(&employee, empID).Error; err != nil {
		return 0, 0
	}

	var attendanceCount int64
	config.DB.Model(&models.Attendance{}).
		Where("user_id = ? AND date BETWEEN ? AND ?", empID, start, end).
		Count(&attendanceCount)

	totalWorkingDays := getWorkingDays(start, end)
	if totalWorkingDays == 0 {
		return 0, 0
	}

	dailyRate := employee.Salary / int64(totalWorkingDays)
	attendancePay := dailyRate * attendanceCount

	return attendancePay, attendanceCount
}

func RunPayroll(id uint, ip, reqID string, userID uint) (string, error) {
	var cycle models.PayCycle
	err := config.DB.First(&cycle, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrPayCycleNotFound
		}
		return "", err
	}

	if cycle.IsProcessed {
		return "", ErrPayCycleIsProcessed
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var employees []models.User
	if err := tx.Where("role = ?", "employee").Find(&employees).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	for _, emp := range employees {
		attendancePay, totalDays := calculateAttendance(emp.ID, cycle.StartDate, cycle.EndDate)
		overtimePay, totalHours := calculateOvertime(emp.ID, cycle.StartDate, cycle.EndDate)
		reimburse := calculateReimbursement(emp.ID, cycle.StartDate, cycle.EndDate)

		payslip := models.Payslip{
			EmployeeID:     emp.ID,
			PayCycleID:     cycle.ID,
			AttendancePay:  attendancePay,
			OvertimePay:    overtimePay,
			Reimbursements: reimburse,
			TotalTakeHome:  attendancePay + overtimePay + reimburse,
			AttendanceDays: int(totalDays),
			OvertimeHours:  float32(totalHours),
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
		}

		if err := tx.Create(&payslip).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	cycle.IsProcessed = true
	if err := tx.Save(&cycle).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	CreateAuditLog(userID, "payroll", cycle.ID, "RUN", "Run pay payroll", ip, reqID)

	return "payroll successfully run", tx.Commit().Error
}
