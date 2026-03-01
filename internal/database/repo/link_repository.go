package repo

import (
	"shortivo/internal/app/link"

	"gorm.io/gorm"
)

type linkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) link.LinkRepository {
	return &linkRepository{db: db}
}
