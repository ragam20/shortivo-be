package response

import "time"

type SuccessResponse[T any] struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data"`
}
