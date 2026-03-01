package auth

import (
	"context"
)

type AuthRepository interface {
	SignUp(ctx context.Context, auth Auth) error
	SignIn(ctx context.Context, email string) (Auth, error)
	SaveRefreshToken(ctx context.Context, token RefreshToken) error
	RevokeTokenByUser(ctx context.Context, userID int64) error
	FindByTokenHash(ctx context.Context, token string) (RefreshToken, error)
	DeleteToken(ctx context.Context, tokenHash string) error
	WithTransaction(ctx context.Context, transaction func(repo AuthRepository) error) error
}
