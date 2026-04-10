package response

import "time"

type LinkResponse struct {
	ID          string     `json:"id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	OriginalURL string     `json:"originalURL"`
	ClickCount  int64      `json:"clickCount"`
	UniqueClick int64      `json:"uniqueClick"`
	ShortURL    string     `json:"shortURL"`
	ExpiredAt   *time.Time `json:"expiredAt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
