package entity

import (
	"shortivo/internal/app/auth"
	"time"
)

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

func (RefreshTokenEntity) TableName() string {
	return "refresh_tokens"
}

func ToRefreshTokenEntity(refToken auth.RefreshToken) RefreshTokenEntity {
	return RefreshTokenEntity{
		ID:        refToken.ID,
		UserID:    refToken.UserID,
		TokenHash: refToken.TokenHash,
		IpAddress: refToken.IpAddress,
		UserAgent: refToken.UserAgent,
		ExpiredAt: refToken.ExpiredAt,
		Revoked:   refToken.Revoked,
		RevokedAt: refToken.RevokedAt,
	}
}

func (e RefreshTokenEntity) ToDomain() auth.RefreshToken {
	return auth.RefreshToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		IpAddress: e.IpAddress,
		UserAgent: e.UserAgent,
		ExpiredAt: e.ExpiredAt,
		Revoked:   e.Revoked,
		RevokedAt: e.RevokedAt,
	}
}
