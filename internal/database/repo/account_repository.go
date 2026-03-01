package repo

import (
	"context"
	"shortivo/internal/app/account"
	"shortivo/internal/database/entity"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) account.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetProfile(ctx context.Context, userID int64) (account.Account, error) {
	var acc entity.UserEntity
	if err := r.db.WithContext(ctx).First(&acc).Error; err != nil {
		return account.Account{}, err
	}

	return acc.ToAccountDomain(), nil
}
