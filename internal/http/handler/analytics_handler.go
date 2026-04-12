package handler

import (
	"net/http"

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

func (h *AnalyticsHandler) GetTotalCreatedLink(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	dateRange, err := GetDateQueries(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := IsAnalyticsDateRangeValid(dateRange.From, dateRange.To); err != nil {
		c.Error(err)
		return
	}

	filter := analytics.FilterQuery{
		FromDate: dateRange.From,
		ToDate:   dateRange.To,
	}

	resp, err := h.analyticsService.GetTotalCreatedLink(c.Request.Context(), userID, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(
		"get created link successfully",
		map[string]any{
			"totalLink": resp,
		}))
}

func (h *AnalyticsHandler) GetTotalClicks(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	dateRange, err := GetDateQueries(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := IsAnalyticsDateRangeValid(dateRange.From, dateRange.To); err != nil {
		c.Error(err)
		return
	}

	filter := analytics.FilterQuery{
		FromDate: dateRange.From,
		ToDate:   dateRange.To,
	}

	result, err := h.analyticsService.GetTotalClicks(c.Request.Context(), userID, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("get total clicks successfully", response.GetClicksResponse{
		TotalClick:  result.TotalClick,
		UniqueClick: result.UniqueClick,
	}))
}

func (h *AnalyticsHandler) GetTotalClicksByLinkID(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	linkID, isExists := linkIDFromParam(c)
	if !isExists {
		return
	}

	dateRange, err := GetDateQueries(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := IsAnalyticsDateRangeValid(dateRange.From, dateRange.To); err != nil {
		c.Error(err)
		return
	}

	filter := analytics.FilterQuery{
		FromDate: dateRange.From,
		ToDate:   dateRange.To,
	}

	result, err := h.analyticsService.GetTotalClicksByLinkID(c.Request.Context(), analytics.GetTotalClicksQuery{
		UserID: userID,
		LinkID: linkID,
	}, filter)
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

	dateRange, err := GetDateQueries(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := IsAnalyticsDateRangeValid(dateRange.From, dateRange.To); err != nil {
		c.Error(err)
		return
	}

	filter := analytics.FilterQuery{
		FromDate: dateRange.From,
		ToDate:   dateRange.To,
	}

	result, err := h.analyticsService.GetEngagementByDevices(c.Request.Context(), analytics.GetTotalClicksQuery{
		UserID: userID,
		LinkID: linkID,
	}, filter)
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

func (h *AnalyticsHandler) GetEnagementByReferrer(c *gin.Context) {
	userID, isExists := userIDFromContext(c)
	if !isExists {
		c.Error(response.Unauthorized())
		return
	}

	linkID, isExists := linkIDFromParam(c)
	if !isExists {
		return
	}

	dateRange, err := GetDateQueries(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := IsAnalyticsDateRangeValid(dateRange.From, dateRange.To); err != nil {
		c.Error(err)
		return
	}

	filter := analytics.FilterQuery{
		FromDate: dateRange.From,
		ToDate:   dateRange.To,
	}

	result, err := h.analyticsService.GetEngagementByReferer(c.Request.Context(), analytics.GetTotalClicksQuery{
		UserID: userID,
		LinkID: linkID,
	}, filter)
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

	dateRange, err := GetDateQueries(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := IsAnalyticsDateRangeValid(dateRange.From, dateRange.To); err != nil {
		c.Error(err)
		return
	}

	filter := analytics.FilterQuery{
		FromDate: dateRange.From,
		ToDate:   dateRange.To,
	}

	result, err := h.analyticsService.GetTotalEngagementByOvertime(c.Request.Context(), analytics.GetTotalEngagementByOvertimeQuery{
		UserID: userID,
		LinkID: linkID,
	}, filter)
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
