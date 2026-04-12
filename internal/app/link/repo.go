package link

import (
	"context"
	"time"
)

type LinkRepository interface {
	GetLinkByID(ctx context.Context, linkID int64) (Link, error)
	CreateLink(ctx context.Context, domain Link) error
	DeleteLink(ctx context.Context, userID int64, linkID int64) error
	UpdateLink(ctx context.Context, userID int64, linkID int64, fields map[string]any) (Link, error)
	GetLinkByCode(ctx context.Context, code string) (Link, error)
	TrackClickURL(ctx context.Context, domain LinkClicks) error
	GetTotalItemByUser(ctx context.Context, userID int64, queries LinkQueries) int64
	GetLinks(ctx context.Context, userID int64, queries LinkQueries) []Link
	IncrementClickCount(ctx context.Context, linkID int64) error
	IncrementUniqueClick(ctx context.Context, linkID int64) error
	GetCurrentLimit(ctx context.Context, userID int64, startAt time.Time, endAt time.Time) (int64, error)
}
