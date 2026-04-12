package handler

import (
	"net/http"
	"strconv"

	"shortivo/internal/app/subscription"
	"shortivo/internal/http/response"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subsService subscription.SubscriptionService
}

func NewSubscriptionHandler(subsService subscription.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subsService: subsService,
	}
}

func (h *SubscriptionHandler) GetCurrentSubscription(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	result, err := h.subsService.GetCurrentSubscription(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	res := response.GetCurrentSubscriptionResponse{
		Plan: response.PlanResponse{
			ID:                  result.Plan.PlanID,
			Title:               result.Plan.Title,
			Code:                result.Plan.Code,
			MonthlyLimit:        result.Plan.MonthlyLimit,
			AnalyticsEnabled:    result.Plan.AnalyticsEnabled,
			CustomDomainEnabled: result.Plan.CustomDomainEnabled,
		},
	}

	if result.Subscription != nil {
		res.Subscription = &response.SubscriptionResponse{
			ID:        result.Subscription.ID,
			PeriodEnd: result.Subscription.PeriodEnd,
			Status:    string(result.Subscription.Status),
		}
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("get current subscription successfully", res))
}

func (h *SubscriptionHandler) GetCurrentPlan(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	plan, err := h.subsService.GetCurrentPlan(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("get plan successfully", response.PlanResponse{
		ID:                  plan.ID,
		Code:                plan.Code,
		Title:               plan.Title,
		MonthlyLimit:        plan.MonthlyLimit,
		AnalyticsEnabled:    plan.AnalyticsEnabled,
		CustomDomainEnabled: plan.CustomDomainEnabled,
	}))
}

func (h *SubscriptionHandler) CancelSubsription(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	paramID := c.Param("id")
	if paramID == "" {
		c.Error(response.BadRequestError("subs id is required"))
		return
	}

	subsID, err := strconv.Atoi(paramID)
	if err != nil {
		c.Error(response.BadRequestError("subs id must be a number"))
		return
	}

	if err := h.subsService.CancelSubscription(c.Request.Context(), subscription.CancelSubscriptionCMD{
		UserID: userID,
		SubsID: int64(subsID),
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("subsciption is cancelled", true))
}
