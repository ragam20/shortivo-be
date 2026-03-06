package request

import "time"

type CreateLinkRequest struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type UpdateLinkRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ExpiredAt   *time.Time `json:"expiredAt"`
}
