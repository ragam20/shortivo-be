package subscription

import "time"

type SubscriptionStatus string

// Domain

const (
	SubscriptionActive    SubscriptionStatus = "active"
	SubscriptionExpired   SubscriptionStatus = "expired"
	SubscriptionCancelled SubscriptionStatus = "cancelled"
)

type Subscription struct {
	ID        int64
	UserID    int64
	PlanID    int64
	PeriodEnd time.Time
	Status    SubscriptionStatus

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

	UpdatedAt time.Time
	CreatedAt time.Time
}

// End Domain
