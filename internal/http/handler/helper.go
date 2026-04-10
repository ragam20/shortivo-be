package handler

import (
	"net/url"
	"strconv"
	"time"

	"shortivo/internal/constant"
	"shortivo/internal/http/request"
	"shortivo/internal/http/response"
	"shortivo/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

func IsAnalyticsDateRangeValid(from *time.Time, to *time.Time) error {
	if from == nil || to == nil {
		return nil
	}

	f := from.UTC()
	t := to.UTC()

	if f.After(t) {
		return response.BadRequestError("invalid date range: 'from' must be before 'to'")
	}

	maxRange := 90 * 24 * time.Hour
	minRange := 1 * time.Hour

	diff := t.Sub(f)

	if diff > maxRange {
		return response.BadRequestError("date range too large (max 90 days)")
	}

	if diff < minRange {
		return response.BadRequestError("date range too small (min 1 hour)")
	}

	return nil
}

func GetDateQueries(c *gin.Context) (*utils.DateRange, error) {
	qCreatedFrom := c.Query("created_from")
	qCreatedTo := c.Query("created_to")
	if qCreatedFrom == "" && qCreatedTo == "" {
		return &utils.DateRange{}, nil
	}

	if qCreatedFrom == "" || qCreatedTo == "" {
		return nil, response.BadRequestError("created_from and created_to must both be provided")
	}

	from, err := utils.ParseRFC3339Time(qCreatedFrom)
	if err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	to, err := utils.ParseRFC3339Time(qCreatedTo)
	if err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	return &utils.DateRange{
		From: from,
		To:   to,
	}, nil
}

type Pagination struct {
	Page  int
	Limit int
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func pagination(c *gin.Context) (Pagination, error) {
	const (
		defaultPage  = 1
		defaultLimit = 10
		maxLimit     = 100
	)

	page := defaultPage
	limit := defaultLimit

	if pageQ := c.Query("page"); pageQ != "" {
		p, err := strconv.Atoi(pageQ)
		if err != nil {
			return Pagination{}, response.BadRequestError("page must be a number")
		}
		page = p
	}

	if limitQ := c.Query("limit"); limitQ != "" {
		l, err := strconv.Atoi(limitQ)
		if err != nil {
			return Pagination{}, response.BadRequestError("limit must be a number")
		}
		limit = l
	}

	if page < 1 {
		page = defaultPage
	}

	if limit < 1 {
		limit = defaultLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}, nil
}

func userIDFromContext(c *gin.Context) (int64, bool) {
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

func validateURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}

func detectDevice(userAgent string) string {
	ua := useragent.New(userAgent)

	if ua.Bot() {
		return "bot"
	}

	if ua.Mobile() {
		return "mobile"
	}

	return "desktop"
}
