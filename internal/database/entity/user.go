package entity

import "time"

type UserEntity struct {
	ID       int64   `gorm:"primaryKey;autoIncrement:false"`
	Email    string  `gorm:"size:255;not null;uniqueIndex"`
	Username string  `gorm:"size:50;not null"`
	Password *string `gorm:"size:255"`
	GoogleID *string `gorm:"size:64;uniqueIndex"`

	UpdatedAt time.Time
	CreatedAt time.Time
}
