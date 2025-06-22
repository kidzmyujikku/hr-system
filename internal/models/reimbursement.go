package models

import "time"

type Reimbursement struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`

	Date        time.Time `gorm:"not null"`
	Amount      int64     `gorm:"not null"`
	Description string    `gorm:"type:varchar(255)"`

	CreatedBy uint
	Creator   User `gorm:"foreignKey:CreatedBy"`
	UpdatedBy *uint
	Updater   *User `gorm:"foreignKey:UpdatedBy"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
