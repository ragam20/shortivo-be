package link

import "time"

/**
	Domain
**/

type Link struct {
	ID          int64
	UserID      int64
	Title       *string
	Description *string
	OriginalURL string
	CodeURL     string
	ClickCount  int64
	UniqueClick int64
	ExpiredAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type LinkClicks struct {
	ID        int64
	LinkID    int64
	IPAddress string
	UserAgent string
	Device    string
	Referer   string

	CreatedAt time.Time
}

// End domain

/**
	Command
**/

type CreateLinkCommand struct {
	UserID      int64
	Title       *string
	URL         string
	Description *string
}

type UpdateLinkCommand struct {
	UserID      int64
	LinkID      int64
	Title       *string
	Description *string
	ExpiredAt   *time.Time
}

type DeleteLinkCommand struct {
	UserID int64
	LinkID int64
}

type RedirectLinkCommand struct {
	Code      string
	IPAddress string
	UserAgent string
	Referer   string
	Device    string
}

// End command

// Query

type LinkQueries struct {
	UserID        int64
	Offset        int
	Limit         int
	Page          int
	SortBy        string
	Order         string
	Search        *string
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}

// End Query
/**

	Result
**/

type LinkResult struct {
	ID          int64
	Title       *string
	Description *string
	OriginalURL string
	ShortURL    string
	ClickCount  int64
	UniqueClick int64
	ExpiredAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetLinkResult struct {
	Links     []LinkResult
	TotalItem int64
	TotalPage int
}

type RedirectLinkResult struct {
	OriginalURL string
}
