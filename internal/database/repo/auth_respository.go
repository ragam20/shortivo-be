package repo

import (
	"context"
	"shortivo/internal/app/auth"
	"shortivo/internal/database/entity"
	"time"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) auth.AuthRepository {
	return &authRepository{db: db}
}

func (a *authRepository) SignIn(ctx context.Context, email string) (auth.Auth, error) {
	var u entity.UserEntity
	if err := a.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return auth.Auth{}, err
	}
	return u.ToAuthDomain(), nil
}

func (a *authRepository) SignUp(ctx context.Context, auth auth.Auth) error {
	u := entity.NewAuth(auth)
	if err := a.db.WithContext(ctx).Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (a *authRepository) SaveRefreshToken(ctx context.Context, token auth.RefreshToken) error {
	if err := a.db.WithContext(ctx).Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (a *authRepository) FindByTokenHash(ctx context.Context, token string) (auth.RefreshToken, error) {
	var refreshToken entity.RefreshTokenEntity
	if err := a.db.WithContext(ctx).
		Model(&entity.RefreshTokenEntity{}).
		Where("token_hash = ?", token).
		First(&refreshToken).Error; err != nil {
		return auth.RefreshToken{}, err
	}
	return refreshToken.ToDomain(), nil
}

func (a *authRepository) RevokeTokenByUser(ctx context.Context, userID int64) error {
	if err := a.db.WithContext(ctx).
		Model(&entity.RefreshTokenEntity{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"revoked_at": time.Now().UTC(),
			"revoked":    true,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (a *authRepository) DeleteToken(ctx context.Context, tokenHash string) error {
	if err := a.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		Delete(&entity.RefreshTokenEntity{}).Error; err != nil {
		return err
	}
	return nil
}
