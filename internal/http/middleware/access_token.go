package middleware

import (
	"shortivo/internal/app/auth"
	"shortivo/internal/constant"
	"shortivo/internal/http/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccessTokenMiddleware struct {
	jwtManager auth.JwtManager
}

func NewAccessTokenMiddleware(jwtManager auth.JwtManager) *AccessTokenMiddleware {
	return &AccessTokenMiddleware{jwtManager: jwtManager}
}

func (m *AccessTokenMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader(constant.Authorization)
		if bearer == "" {
			c.Error(response.Unauthorized())
			c.Abort()
			return
		}

		s := strings.Split(bearer, " ")
		if len(s) < 2 || strings.ToLower(s[0]) != "bearer" {
			c.Error(response.Unauthorized())
			c.Abort()
			return
		}

		jwtUser, err := m.jwtManager.ValidateAccessToken(s[1])
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Set(constant.UserIDCtx, jwtUser.UserId)
		c.Next()
	}
}
