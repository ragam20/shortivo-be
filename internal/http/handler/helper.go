package handler

import (
	"shortivo/internal/constant"
	"shortivo/internal/http/request"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserIDFromContext(c *gin.Context) (int64, bool) {
	val, exists := c.Get(constant.UserIDCtx)

	if exists {
		switch v := val.(type) {
		case int64:
			return v, true
		case int:
			return int64(v), true
		case float64:
			return int64(v), true
		case string:
			parsed, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				return parsed, true
			}
		}
	}

	var req request.SignOutRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.UserID != 0 {
		return req.UserID, true
	}

	return 0, false
}
