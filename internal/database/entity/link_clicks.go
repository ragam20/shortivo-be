package entity

import "time"

type LinkClicksEntity struct {
	ID        int64 `gorm:"primaryKey;autoIncrement:false"`
	LinkID    int64 `gorm:"index"`
	IpAddress string
	UserAgent string
	Device    string
	Referrer  *string `gorm:"null"`

	CreatedAt time.Time
}
