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

func (d *Subscription) IsStatusActive() bool {
	return d.Status == SubscriptionActive
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

type CancelSubscriptionCMD struct {
	UserID int64
	SubsID int64
}
