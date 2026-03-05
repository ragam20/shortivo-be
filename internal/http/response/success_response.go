package response

import (
	"net/http"
	"time"
)

type Meta map[string]any

type SuccessResponse[T any] struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data"`
	Meta      *Meta     `json:"meta,omitempty"`
}

func NewSuccessResponse[T any](message string, data T) SuccessResponse[T] {
	return SuccessResponse[T]{
		Status:    http.StatusOK,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
}

func NewSuccessResponseWithMeta[T any](message string, data T, meta *Meta) SuccessResponse[T] {
	return SuccessResponse[T]{
		Status:    http.StatusOK,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
		Meta:      meta,
	}
}
