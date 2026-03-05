package link

import "time"

/**
	Domain
**/

type Link struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
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
	IpAddress string
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
	UserID int64
	Title  string
	URL    string
}

type UpdateLinkCommand struct {
	UserID      int64
	LinkID      int64
	Title       string
	Description string
	ExpiredAt   *time.Time
}

type DeleteLinkCommand struct {
	UserID int64
	LinkID int64
}

type RedirectLinkCommand struct {
	Code      string
	IpAddress string
	UserAgent string
	Referer   string
	Device    string
}

// End command

// Query

type GetLinksQuery struct {
	UserID int64
	Offset int
	Limit  int
	Page   int
}

// End Query

/**
	Result
**/

type LinkResult struct {
	ID          int64
	Title       string
	Description string
	OriginalURL string
	ShortURL    string
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
