package account

import (
	"context"
)

type AccountRepository interface {
	GetProfile(ctx context.Context, userID int64) (Account, error)
}
