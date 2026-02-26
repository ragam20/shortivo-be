package entity

import "time"

type UserEntity struct {
	ID       int64  `gorm:"primaryKey;autoIncrement:false"`
	Email    string `gorm:"unique"`
	Username string
	Password string
	GoogleID *int64 `gorm:"null"`

	UpdatedAt time.Time
	CreatedAt time.Time
}
