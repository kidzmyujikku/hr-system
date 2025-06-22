package models

import "time"

type Attendance struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"not null"`
	User     User
	Date     time.Time `gorm:"not null;index"` // YYYY-MM-DD
	CheckIn  time.Time `gorm:"not null"`
	CheckOut *time.Time

	CreatedBy uint `gorm:"not null"`
	Creator   User `gorm:"foreignKey:CreatedBy"`

	UpdatedBy *uint //
	Updater   *User `gorm:"foreignKey:UpdatedBy"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
