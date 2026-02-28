package middleware

import (
	"errors"
	"net/http"
	"shortivo/internal/http/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ErrorMiddleware struct {
	log *zerolog.Logger
}

func NewErrorMiddleware(log *zerolog.Logger) *ErrorMiddleware {
	return &ErrorMiddleware{log: log}
}

func (m *ErrorMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		status := http.StatusInternalServerError
		msg := "internal server error"
		timestamp := time.Now().UTC()

		var apperr *response.AppErrResponse
		if errors.As(err, &apperr) {
			status = apperr.Status
			timestamp = apperr.Timestamp

			if status >= 500 {
				m.log.Error().
					Err(err).
					Int("status", status).
					Str("timestamp", apperr.Timestamp.String()).
					Msg(apperr.Message)
			}

			if status >= 400 && status < 500 {
				msg = apperr.Message
			}
		} else {
			m.log.Error().Err(err).Msg("unexpected error")
		}

		c.AbortWithStatusJSON(status, gin.H{
			"status":    status,
			"message":   msg,
			"timestamp": timestamp,
		})
	}
}
