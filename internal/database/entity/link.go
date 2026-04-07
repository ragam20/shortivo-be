package entity

import (
	"time"

	"shortivo/internal/app/link"

	"gorm.io/gorm"
)

type LinkEntity struct {
	ID          int64      `gorm:"primaryKey;autoIncrement:false"`
	UserID      int64      `gorm:"not null;index"`
	User        UserEntity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	Title       *string    `gorm:"size:255"`
	Description *string    `gorm:"type:text"`
	OriginalURL string     `gorm:"column:original_url;type:text;not null"`
	CodeURL     string     `gorm:"column:code_url;size:32;not null;uniqueIndex"`
	ClickCount  int64      `gorm:"not null;default:0"`
	UniqueClick int64      `gorm:"not null;default:0"`
	ExpiredAt   *time.Time `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

func (LinkEntity) TableName() string {
	return "links"
}

func (e *LinkEntity) ToDomain() link.Link {
	return link.Link{
		ID:          e.ID,
		UserID:      e.UserID,
		Title:       e.Title,
		Description: e.Description,
		OriginalURL: e.OriginalURL,
		CodeURL:     e.CodeURL,
		ClickCount:  e.ClickCount,
		UniqueClick: e.ClickCount,
		ExpiredAt:   e.ExpiredAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func NewLink(domain link.Link) LinkEntity {
	return LinkEntity{
		ID:          domain.ID,
		UserID:      domain.UserID,
		Title:       domain.Title,
		Description: domain.Description,
		OriginalURL: domain.OriginalURL,
		CodeURL:     domain.CodeURL,
		ClickCount:  domain.ClickCount,
		UniqueClick: domain.UniqueClick,
		ExpiredAt:   domain.ExpiredAt,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}
