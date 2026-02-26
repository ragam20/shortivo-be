package entity

import "time"

type LinkClicksEntity struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:false"`
	LinkID    int64      `gorm:"not null;index"`
	Link      LinkEntity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:LinkID;references:ID"`
	IpAddress string     `gorm:"size:64;not null;index"`
	UserAgent string     `gorm:"size:512;not null"`
	Device    string     `gorm:"size:255;not null"`
	Referrer  *string    `gorm:"type:text"`

	CreatedAt time.Time
}
