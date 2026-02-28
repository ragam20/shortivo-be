package middleware

import (
	"context"
	"net/http"
	"shortivo/internal/config"
	"shortivo/internal/http/response"

	"github.com/gin-gonic/gin"
)

type RequestTimeoutMiddleware struct {
	cfg *config.AppConfig
}

func NewRequestTimeoutMiddleware(cfg *config.AppConfig) *RequestTimeoutMiddleware {
	return &RequestTimeoutMiddleware{cfg: cfg}
}

func (m *RequestTimeoutMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), m.cfg.RequestTimeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, response.NewAppErr(
				http.StatusGatewayTimeout,
				"Request timeout",
			))
		case <-done:
		}
	}
}
