package entity

import "time"

type RefreshTokenEntity struct {
	ID        int64 `gorm:"primaryKey;autoIncrement:false"`
	UserID    int64
	TokenHash string
	UserAgent string
	IpAddress string
	ExpiredAt time.Time
	Revoked   bool
	RevokedAt time.Time

	CreatedAt time.Time
}
