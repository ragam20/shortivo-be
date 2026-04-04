package handler

import (
	"net/http"
	"time"

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

func (h *AccountHandler) GetAccount(c *gin.Context) {
	userID := c.GetInt64(constant.UserIDCtx)
	if userID <= 0 {
		c.Error(response.Unauthorized())
		return
	}

	userPlan, err := h.accService.GetAccount(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	var planResp *response.PlanResponse
	if userPlan.Plan != nil {
		planResp = &response.PlanResponse{
			ID:                  userPlan.Plan.ID,
			Title:               userPlan.Plan.Title,
			Code:                userPlan.Plan.Code,
			MonthlyLimit:        userPlan.Plan.MonthlyLimit,
			AnalyticsEnabled:    userPlan.Plan.AnalyticsEnabled,
			CustomDomainEnabled: userPlan.Plan.CustomDomainEnabled,
		}
	}

	c.JSON(http.StatusOK, response.SuccessResponse[response.GetAccountResponse]{
		Status:    http.StatusOK,
		Message:   "get user successfully",
		Timestamp: time.Now().Truncate(time.Second),
		Data: response.GetAccountResponse{
			User: response.GetUserResponse{
				ID:        userPlan.User.ID,
				Email:     userPlan.User.Email,
				Username:  userPlan.User.Username,
				GoogleID:  userPlan.User.GoogleID,
				CreatedAt: userPlan.User.CreatedAt.Truncate(time.Second),
				UpdatedAt: userPlan.User.UpdatedAt.Truncate(time.Second),
			},
			Plan: planResp,
		},
	})
}
