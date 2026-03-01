package entity

import "time"

type LinkEntity struct {
	ID          int64      `gorm:"primaryKey;autoIncrement:false"`
	UserID      int64      `gorm:"not null;index"`
	User        UserEntity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	Title       string     `gorm:"size:255;not null"`
	OriginalURL string     `gorm:"column:original_url;type:text;not null"`
	CodeURL     string     `gorm:"column:code_url;size:32;not null;uniqueIndex"`
	ClickCount  int64      `gorm:"not null;default:0"`
	UniqueClick int64      `gorm:"not null;default:0"`
	ExpiredAt   *time.Time `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (LinkEntity) TableName() string {
	return "links"
}
