package entity

import (
	"shortivo/internal/app/account"
	"shortivo/internal/app/auth"
	"time"
)

type UserEntity struct {
	ID       int64   `gorm:"primaryKey;autoIncrement:false"`
	Email    string  `gorm:"size:255;not null;uniqueIndex"`
	Username string  `gorm:"size:50;not null"`
	Password *string `gorm:"size:255"`
	GoogleID *string `gorm:"size:64;uniqueIndex"`

	UpdatedAt time.Time
	CreatedAt time.Time
}

func (UserEntity) TableName() string {
	return "users"
}

func NewAuth(auth auth.Auth) UserEntity {
	return UserEntity{
		ID:        auth.ID,
		Email:     auth.Email,
		Username:  auth.Username,
		GoogleID:  auth.GoogleID,
		CreatedAt: auth.CreatedAt,
		UpdatedAt: auth.UpdatedAt,
	}
}

func (u UserEntity) ToAuthDomain() auth.Auth {
	return auth.Auth{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		GoogleID:  u.GoogleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func NewUser(account account.Account) UserEntity {
	return UserEntity{
		ID:        account.ID,
		Email:     account.Email,
		Username:  account.Username,
		Password:  account.Password,
		GoogleID:  account.GoogleID,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}

func (u UserEntity) ToAccountDomain() account.Account {
	return account.Account{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		GoogleID:  u.GoogleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
