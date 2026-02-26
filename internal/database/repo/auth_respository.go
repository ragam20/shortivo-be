package repo

import (
	"shortivo/internal/app/auth"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) auth.AuthRepository {
	return authRepository{db: db}
}
