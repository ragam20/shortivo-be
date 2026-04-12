package auth

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	UserAgent string
	IpAddress string
	ExpiredAt time.Time
	Revoked   bool
	RevokedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
