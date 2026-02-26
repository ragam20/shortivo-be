package entity

import "time"

type RefreshTokenEntity struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:false"`
	UserID    int64      `gorm:"not null;index"`
	User      UserEntity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	TokenHash string     `gorm:"size:255;not null;uniqueIndex"`
	UserAgent string     `gorm:"size:512;not null"`
	IpAddress string     `gorm:"size:64;not null"`
	ExpiredAt time.Time  `gorm:"not null;index"`
	Revoked   bool       `gorm:"not null;default:false;index"`
	RevokedAt *time.Time

	UpdatedAt time.Time
	CreatedAt time.Time
}
