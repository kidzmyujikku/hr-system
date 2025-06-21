package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"type:varchar(10);not null;check:role IN ('admin','employee')";default:'employee'"`
	Salary   int64  `gorm:"not null" json:"-"`

	CreatedAt time.Time
}
