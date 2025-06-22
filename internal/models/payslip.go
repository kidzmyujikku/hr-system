package models

import "time"

type Payslip struct {
	ID             uint `gorm:"primaryKey"`
	EmployeeID     uint
	PayCycleID     uint
	AttendanceDays int
	AttendancePay  int64
	OvertimeHours  float32
	OvertimePay    int64
	Reimbursements int64
	TotalTakeHome  int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint
	Creator        User `gorm:"foreignKey:CreatedBy"`
}
