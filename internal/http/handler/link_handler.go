package handler

import (
	"net/http"
	"shortivo/internal/app/link"
	"shortivo/internal/http/request"
	"shortivo/internal/http/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	linkService link.LinkService
}

func NewLinkHandler(linkService link.LinkService) *LinkHandler {
	return &LinkHandler{
		linkService: linkService,
	}
}

func (h *LinkHandler) GetLinks(c *gin.Context) {
	userID, exists := userIDFromContext(c)
	if !exists {
		c.Error(response.Unauthorized())
		return
	}

	pagination, err := pagination(c)
	if err != nil {
		c.Error(err)
		return
	}

	result := h.linkService.GetLinks(c.Request.Context(), link.GetLinksQuery{
		UserID: userID,
		Offset: pagination.Offset(),
		Limit:  pagination.Limit,
		Page:   pagination.Page,
	})

	resp := make([]response.LinkResponse, len(result.Links))
	for i, v := range result.Links {
		resp[i] = response.LinkResponse{
			ID:          v.ID,
			Title:       v.Title,
			ShortURL:    v.ShortURL,
			OriginalURL: v.OriginalURL,
			Description: v.Description,
			ExpiredAt:   v.ExpiredAt,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response.NewSuccessResponseWithMeta(
		"success get links",
		resp,
		&response.Meta{
			"page":      pagination.Page,
			"limit":     pagination.Limit,
			"totalItem": result.TotalItem,
			"totalPage": result.TotalPage,
		},
	))
}

func (h *LinkHandler) RedirectShortLink(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.Error(response.BadRequestError("invalid code"))
		return
	}

	userAgent := strings.TrimSpace(c.Request.UserAgent())
	if userAgent == "" {
		userAgent = "unknown"
	}

	ipAddress := strings.TrimSpace(c.ClientIP())
	if ipAddress == "" {
		ipAddress = "unknown"
	}

	referer := strings.TrimSpace(c.GetHeader("Referer"))
	device := detectDevice(userAgent)

	originalURL, err := h.linkService.RedirectLink(c.Request.Context(), link.RedirectLinkCommand{
		Code:      code,
		UserAgent: userAgent,
		IpAddress: ipAddress,
		Referer:   referer,
		Device:    device,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

func (h *LinkHandler) CreateShortLink(c *gin.Context) {
	userID, exists := userIDFromContext(c)
	if !exists {
		c.Error(response.Unauthorized())
		return
	}

	var req request.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(response.BadRequestError("invalid request: " + err.Error()))
		return
	}

	if err := h.linkService.CreateShortLink(c.Request.Context(), &link.CreateLinkCommand{
		UserID:      userID,
		Title:       req.Title,
		URL:         req.URL,
		Description: req.Description,
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("success create a new shortlink", true))
}

func (h *LinkHandler) UpdateLink(c *gin.Context) {
	userID, exists := userIDFromContext(c)
	if !exists {
		c.Error(response.Unauthorized())
		return
	}

	linkID, ok := linkIDFromParam(c)
	if !ok {
		return
	}

	var req request.UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(response.BadRequestError("invalid request: " + err.Error()))
		return
	}

	result, err := h.linkService.UpdateLink(c.Request.Context(), link.UpdateLinkCommand{
		UserID:      userID,
		LinkID:      linkID,
		Title:       req.Title,
		Description: req.Description,
		ExpiredAt:   req.ExpiredAt,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("success update link", response.LinkResponse{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		OriginalURL: result.OriginalURL,
		ShortURL:    result.ShortURL,
		ExpiredAt:   result.ExpiredAt,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}))
}

func (h *LinkHandler) DeleteLink(c *gin.Context) {
	userID, exists := userIDFromContext(c)
	if !exists {
		c.Error(response.Unauthorized())
		return
	}

	linkID, ok := linkIDFromParam(c)
	if !ok {
		return
	}

	if err := h.linkService.DeleteLink(c.Request.Context(), link.DeleteLinkCommand{
		UserID: userID,
		LinkID: linkID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("success delete link", true))
}

func linkIDFromParam(c *gin.Context) (int64, bool) {
	paramLinkID := c.Param("linkId")
	if paramLinkID == "" {
		c.Error(response.BadRequestError("link id is required"))
		return 0, false
	}

	linkID, err := strconv.ParseInt(paramLinkID, 10, 64)
	if err != nil || linkID <= 0 {
		c.Error(response.BadRequestError("invalid link id"))
		return 0, false
	}

	return linkID, true
}
