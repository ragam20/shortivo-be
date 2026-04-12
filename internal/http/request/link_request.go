package request

import "time"

type CreateLinkRequest struct {
	Title       *string `json:"title,omitempty"`
	URL         string  `json:"url"`
	Description *string `json:"description,omitempty"`
}

type UpdateLinkRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	ExpiredAt   *time.Time `json:"expiredAt"`
}
