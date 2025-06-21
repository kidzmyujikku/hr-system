package models

import "time"

type PayCycle struct {
	ID          uint      `gorm:"primaryKey"`
	StartDate   time.Time `gorm:"not null"`
	EndDate     time.Time `gorm:"not null"`
	MonthLabel  string    `gorm:"type:varchar(20)"`
	IsProcessed bool      `gorm:"default:false"`

	CreatedBy uint `gorm:"not null"`
	Creator   User `gorm:"foreignKey:CreatedBy"`

	UpdatedBy *uint //
	Updater   *User `gorm:"foreignKey:UpdatedBy"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
