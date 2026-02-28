package auth

import "time"

type Auth struct {
	ID       int64
	Email    string
	Username string
	Password *string
	GoogleID *string

	UpdatedAt time.Time
	CreatedAt time.Time
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	RefreshExp   time.Time
}
