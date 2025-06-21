package models

import "time"

type AuditLog struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	Entity      string    `gorm:"not null"`
	EntityID    uint      `gorm:"not null"`
	Action      string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	IPAddress   string    `gorm:"type:varchar(45)"`
	RequestID   string    `gorm:"type:varchar(64)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
