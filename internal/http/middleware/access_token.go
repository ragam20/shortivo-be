package middleware

import (
	"net/http"
	"shortivo/internal/app/auth"
	"shortivo/internal/constant"
	"shortivo/internal/http/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccessTokenMiddleware struct {
	jwtManager auth.JwtManager
}

func NewAccesTokenMiddleware(jwtManager auth.JwtManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader(constant.Authorization)
		if bearer == "" {
			c.AbortWithError(http.StatusUnauthorized, response.Unauthorized())
			return
		}

		s := strings.Split(bearer, " ")
		if len(s) < 2 || strings.ToLower(s[0]) != "bearer" {
			c.AbortWithError(http.StatusUnauthorized, response.Unauthorized())
			return
		}

		jwtUser, err := jwtManager.ValidateAccessToken(s[1])
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, response.Unauthorized())
			return
		}

		c.Set(constant.UserIDCtx, jwtUser.UserId)

		c.Next()
	}
}
