package middleware

import (
	"net/http"
	"shortivo/internal/http/response"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type RecoveryMiddleware struct {
	log *zerolog.Logger
}

func NewRecoveryMiddleware(log *zerolog.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{log: log}
}

func (m *RecoveryMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {

				log.Error().
					Interface("panic", r).
					Str("path", c.Request.URL.Path).
					Msg("panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerError(""))
			}
		}()

		c.Next()
	}
}
