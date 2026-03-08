package subscription

import (
	"context"
)

type SubscriptionRepository interface {
	GetSubscriptionByUser(ctx context.Context, userID int64) (Subscription, error)
	ChangeStatus(ctx context.Context, userID, subsID int64, sub SubscriptionStatus) error
}
