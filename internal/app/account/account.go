package account

import "time"

type Account struct {
	ID       int64
	Email    string
	Username string
	Password *string
	GoogleID *string

	UpdatedAt time.Time
	CreatedAt time.Time
}

type GetProfile struct {
	ID       int64
	Email    string
	Username string
	GoogleID *string

	UpdatedAt time.Time
	CreatedAt time.Time
}
