package request

import "time"

type CreateLinkRequest struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type UpdateLinkRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ExpiredAt   *time.Time `json:"expiredAt"`
}
