package models

import (
	"time"
)

type Overtime struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	UserID    uint `gorm:"not null"`
	User      User
	Date      time.Time
	Hours     float32
	CreatedBy uint `gorm:"not null"`
	Creator   User `gorm:"foreignKey:CreatedBy"`

	UpdatedBy *uint //
	Updater   *User `gorm:"foreignKey:UpdatedBy"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
