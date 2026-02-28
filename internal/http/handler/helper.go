package handler

import (
	"errors"
	"shortivo/internal/http/request"
	"shortivo/internal/http/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserIDFromContext(c *gin.Context) (int64, bool) {
	val, exists := c.Get("userId")
	if !exists {
		val, exists = c.Get("userID")
	}

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

func WriteError(c *gin.Context, err error) {
	var appErr response.AppErrResponse
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, appErr)
		return
	}

	c.Error(err)
}
