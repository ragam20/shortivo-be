package response

import "time"

type LinkResponse struct {
	ID          string     `json:"id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	OriginalURL string     `json:"originalURL"`
	ShortURL    string     `json:"shortURL"`
	ExpiredAt   *time.Time `json:"expiredAt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
