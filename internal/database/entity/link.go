package entity

import "time"

type LinkEntity struct {
	ID          int64 `gorm:"primaryKey;autoIncrement:false"`
	UserID      int64 `gorm:"index"`
	Title       string
	OriginalUrl string
	CodeUrl     string `gorm:"unique;index"`
	ClickCount  int32
	UniqueClick int32
	ExpiredAt   *time.Time `gorm:"null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
