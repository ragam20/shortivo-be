package link

import (
	"context"
)

type LinkRepository interface {
	CreateLink(ctx context.Context, domain Link) error
	DeleteLink(ctx context.Context, userID int64, linkID int64) error
	UpdateLink(ctx context.Context, userID int64, linkID int64, fields map[string]any) (Link, error)
	GetLinkByCode(ctx context.Context, code string) (Link, error)
	TrackClickURL(ctx context.Context, domain LinkClicks) error
	GetLinks(ctx context.Context, userID int64, offset int, limit int) []Link
	GetTotalItemByUser(ctx context.Context, userID int64) int64
	IncrementClickCount(ctx context.Context, linkID int64) error
	IncrementUniqueClick(ctx context.Context, linkID int64) error
}
