package handler

import (
	"net/http"
	"shortivo/internal/app/account"
	"shortivo/internal/constant"
	"shortivo/internal/http/response"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accService account.AccountService
}

func NewAccountHandler(accService account.AccountService) *AccountHandler {
	return &AccountHandler{accService: accService}
}

func (h *AccountHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64(constant.UserIDCtx)
	if userID <= 0 {
		c.Error(response.Unauthorized())
		return
	}

	profile, err := h.accService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse[response.GetProfileResponse]{
		Status:  http.StatusOK,
		Message: "get profile success",
		Data: response.GetProfileResponse{
			ID:        profile.ID,
			Username:  profile.Username,
			Email:     profile.Email,
			GoogleID:  profile.GoogleID,
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
		},
	})
}
