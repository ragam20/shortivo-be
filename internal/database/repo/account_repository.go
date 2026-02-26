package repo

import (
	"shortivo/internal/app/account"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) account.AccountRepository {
	return accountRepository{}
}
