package handler

import (
	"net/http"
	"time"

	"shortivo/internal/app/analytics"
	"shortivo/internal/http/response"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService analytics.AnalyticsService
}

func NewAnalyticsHandler(analyticsService analytics.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetTotalClicks(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	linkID, isExists := linkIDFromParam(c)
	if !isExists {
		return
	}

	result, err := h.analyticsService.GetTotalClicks(c.Request.Context(), analytics.GetTotalClicksQuery{
		UserID: userID,
		LinkID: linkID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("get total clicks successfully", response.GetClicksResponse{
		TotalClick:  result.TotalClick,
		UniqueClick: result.UniqueClick,
	}))
}

func (h *AnalyticsHandler) GetEngagementByDevices(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	linkID, isExists := linkIDFromParam(c)
	if !isExists {
		return
	}

	result, err := h.analyticsService.GetEngagementByDevices(c.Request.Context(), analytics.GetTotalClicksQuery{
		UserID: userID,
		LinkID: linkID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	resp := make([]response.GetTotalEngagementByDeviceResponse, len(result))
	for i, r := range result {
		resp[i] = response.GetTotalEngagementByDeviceResponse{
			Device: r.Device,
			Clicks: r.Clicks,
		}
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("get engagement by device successfully", resp))
}

func (h *AnalyticsHandler) GetEnagementByReferer(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	linkID, isExists := linkIDFromParam(c)
	if !isExists {
		return
	}

	result, err := h.analyticsService.GetEngagementByReferer(c.Request.Context(), analytics.GetTotalClicksQuery{
		UserID: userID,
		LinkID: linkID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	resp := make([]response.GetTotalEngagementByRefererResponse, len(result))
	for i, v := range result {
		resp[i] = response.GetTotalEngagementByRefererResponse{
			Referer: v.Referer,
			Clicks:  v.Clicks,
		}
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("get engagement by referer successfully", resp))
}

func (h *AnalyticsHandler) GetTotalEngagementByOvertime(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	linkID, isExists := linkIDFromParam(c)
	if !isExists {
		return
	}

	rangeQ := c.Query("range")

	var startTime time.Time
	var endTime time.Time

	if rangeQ != "" {
		switch rangeQ {
		case "30d":
			startTime = time.Now().AddDate(0, 0, -30)
		case "90d":
			startTime = time.Now().AddDate(0, 0, -90)
		default:
			startTime = time.Now().AddDate(0, 0, -7)
		}

		endTime = time.Now()
	} else {
		startTimeQ := c.Query("startTime")
		endTimeQ := c.Query("endTime")

		var err error

		startTime, err = time.Parse("2006-01-02", startTimeQ)
		if err != nil {
			c.Error(response.BadRequestError("invalid parse startTime"))
			return
		}

		endTime, err = time.Parse("2006-01-02", endTimeQ)
		if err != nil {
			c.Error(response.BadRequestError("invalid parse endTime"))
			return
		}

		if startTime.After(endTime) {
			c.Error(response.BadRequestError("endtTime must be after startTime"))
			return
		}

		dur := endTime.Sub(startTime)
		if dur > 90*24*time.Hour {
			c.Error(response.BadRequestError("range cannot exceed 90 days"))
			return
		}
	}

	result, err := h.analyticsService.GetTotalEngagementByOvertime(c.Request.Context(), analytics.GetTotalEngagementByOvertimeQuery{
		UserID:    userID,
		LinkID:    linkID,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		c.Error(err)
		return
	}

	resp := make([]response.GetTotalEngagementByOvertimeResponse, len(result))
	for i, v := range result {
		resp[i] = response.GetTotalEngagementByOvertimeResponse{
			Date:   v.Date,
			Clicks: v.Clicks,
		}
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("Get total engagement successfully", resp))
}
