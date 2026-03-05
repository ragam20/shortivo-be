package middleware

import (
	"errors"
	"net/http"
	"shortivo/internal/http/response"
	"shortivo/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorMiddleware struct {
	log logger.Logger
}

func NewErrorMiddleware(log logger.Logger) *ErrorMiddleware {
	return &ErrorMiddleware{log: log}
}

func (m *ErrorMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

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
				m.log.Error(
					apperr.Message,
					err,
					map[string]any{
						"status":    apperr.Status,
						"timestamp": apperr.Timestamp,
						"path":      c.Request.URL.Path,
						"ipAddress": c.ClientIP(),
					},
				)
			}

			if status >= 400 && status < 500 {
				msg = apperr.Message
			}
		}

		c.AbortWithStatusJSON(status, gin.H{
			"status":    status,
			"message":   msg,
			"timestamp": timestamp,
		})
	}
}
