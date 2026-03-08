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
