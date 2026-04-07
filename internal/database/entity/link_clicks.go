package entity

import (
	"time"

	"shortivo/internal/app/link"
)

type LinkClicksEntity struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:false"`
	LinkID    int64      `gorm:"not null;index"`
	Link      LinkEntity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:LinkID;references:ID"`
	IPAddress string     `gorm:"size:64;not null;index"`
	UserAgent string     `gorm:"size:512;not null"`
	Device    string     `gorm:"size:255;not null"`
	Referer   string     `gorm:"type:text"`

	CreatedAt time.Time
}

func (LinkClicksEntity) TableName() string {
	return "link_clicks"
}

func (e *LinkClicksEntity) ToDomain() link.LinkClicks {
	return link.LinkClicks{
		ID:        e.ID,
		LinkID:    e.ID,
		IPAddress: e.IPAddress,
		UserAgent: e.UserAgent,
		Device:    e.Device,
		Referer:   e.Referer,
		CreatedAt: e.CreatedAt,
	}
}

func NewLinkClicks(domain link.LinkClicks) LinkClicksEntity {
	return LinkClicksEntity{
		ID:        domain.ID,
		LinkID:    domain.LinkID,
		IPAddress: domain.IPAddress,
		UserAgent: domain.UserAgent,
		Device:    domain.Device,
		Referer:   domain.Referer,
		CreatedAt: domain.CreatedAt,
	}
}
