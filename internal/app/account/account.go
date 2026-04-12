package account

import (
	"time"
)

type Account struct {
	ID       int64
	Email    string
	Username string
	Password *string
	GoogleID *string

	UpdatedAt time.Time
	CreatedAt time.Time
}

type GetAccount struct {
	ID       int64
	Email    string
	Username string
	GoogleID *string

	UpdatedAt time.Time
	CreatedAt time.Time
}

type Plan struct {
	ID                  int64
	Code                string
	Title               string
	MonthlyLimit        int32
	AnalyticsEnabled    bool
	CustomDomainEnabled bool
}

type UserPlan struct {
	User GetAccount
	Plan *Plan
}
