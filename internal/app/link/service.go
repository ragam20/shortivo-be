package link

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"shortivo/internal/app/subscription"
	"shortivo/internal/http/response"
	"shortivo/internal/idgen"
	"shortivo/internal/logger"
	"shortivo/internal/utils"

	"gorm.io/gorm"
)

type LinkService interface {
	CreateShortLink(ctx context.Context, cmd *CreateLinkCommand) error
	UpdateLink(ctx context.Context, cmd UpdateLinkCommand) (LinkResult, error)
	DeleteLink(ctx context.Context, cmd DeleteLinkCommand) error
	RedirectLink(ctx context.Context, cmd RedirectLinkCommand) (string, error)
	GetLinks(ctx context.Context, query GetLinksQuery) GetLinkResult
	GetMontlyLimit(ctx context.Context, userID int64) (response.GetLinkLimit, error)
}

type linkService struct {
	cfg      *LinkConfig
	linkRepo LinkRepository
	subsRepo subscription.SubscriptionRepository
	planRepo subscription.PlanRepository
	idGen    idgen.IDGenerator
	logger   logger.Logger
	cache    LinkCache
}

func NewLinkService(cfg *LinkConfig,
	linkRepo LinkRepository,
	subsRepo subscription.SubscriptionRepository,
	planRepo subscription.PlanRepository,
	idGen idgen.IDGenerator,
	logger logger.Logger,
	cache LinkCache,
) LinkService {
	return &linkService{
		cfg:      cfg,
		linkRepo: linkRepo,
		subsRepo: subsRepo,
		planRepo: planRepo,
		idGen:    idGen,
		logger:   logger,
		cache:    cache,
	}
}

func (s *linkService) GetLinks(ctx context.Context, query GetLinksQuery) GetLinkResult {
	cacheLink, err := s.cache.GetLinks(ctx, query.UserID, query.Page, query.Limit)
	if err == nil {
		return cacheLink
	}

	links := s.linkRepo.GetLinks(ctx, query.UserID, query.Offset, query.Limit)

	linkItems := make([]LinkResult, len(links))
	for i, v := range links {
		linkItems[i] = LinkResult{
			ID:          v.ID,
			Title:       v.Title,
			OriginalURL: v.OriginalURL,
			Description: v.Description,
			ShortURL:    fmt.Sprintf("%s/%s", s.cfg.ShortenerBaseURL, v.CodeURL),
			ExpiredAt:   v.ExpiredAt,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}

	totalItem := s.linkRepo.GetTotalItemByUser(ctx, query.UserID)

	linkResult := GetLinkResult{
		Links:     linkItems,
		TotalItem: totalItem,
		TotalPage: (int(totalItem) + query.Limit - 1) / query.Limit,
	}

	s.cache.SaveListLinks(ctx, query.UserID, query.Page, query.Limit, linkResult)

	return linkResult
}

func (s *linkService) CreateShortLink(ctx context.Context, cmd *CreateLinkCommand) error {
	if cmd.UserID <= 0 {
		return response.Unauthorized()
	}

	cleanTitle := strings.TrimSpace(cmd.Title)
	if cleanTitle == "" {
		return response.BadRequestError("title is required")
	}

	cleanURL := strings.TrimSpace(cmd.URL)
	if cleanURL == "" {
		return response.BadRequestError("URL is required")
	}

	if !validateURL(cleanURL) {
		return response.BadRequestError("must be valid URL")
	}

	isAllowed, err := s.isAllowedToCreateLink(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	if !isAllowed {
		return response.BadRequestError("you have reached your monthly link limit.")
	}

	linkID := s.idGen.NewID()

	shortCode, err := encodeBase62(linkID)
	if err != nil {
		return response.InternalServerError("error generate shortCode for url")
	}

	link := Link{
		ID:          linkID,
		UserID:      cmd.UserID,
		Title:       cleanTitle,
		Description: cmd.Description,
		OriginalURL: cleanURL,
		CodeURL:     shortCode,
	}
	if err := s.linkRepo.CreateLink(ctx, link); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return response.BadRequestError("already exists")
		}
		return response.InternalServerError("error create link")
	}

	return nil
}

func (s *linkService) isAllowedToCreateLink(ctx context.Context, userID int64) (bool, error) {
	now := time.Now()

	var startTime time.Time
	var endTime time.Time
	var plan subscription.Plan

	sub, err := s.subsRepo.GetActiveSubscriptionByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// FREE plan
			startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endTime = startTime.AddDate(0, 1, 0)

			plan, err = s.planRepo.GetPlanByCode(ctx, "free")
			if err != nil {
				return false, response.InternalServerError("error get free plan")
			}

		} else {
			return false, response.InternalServerError("error get subscription")
		}
	} else {
		startTime = sub.CreatedAt
		endTime = sub.PeriodEnd

		plan, err = s.planRepo.GetPlanByID(ctx, sub.PlanID)
		if err != nil {
			return false, response.InternalServerError("error get plan")
		}
	}

	currLimit, err := s.linkRepo.GetCurrentLimit(ctx, userID, startTime, endTime)
	if err != nil {
		return false, response.InternalServerError("error get current limit")
	}

	return currLimit < int64(plan.MonthlyLimit), nil
}

func (s *linkService) UpdateLink(ctx context.Context, cmd UpdateLinkCommand) (LinkResult, error) {
	s.cache.InvalidateLinks(ctx, cmd.UserID)

	if cmd.UserID <= 0 {
		return LinkResult{}, response.BadRequestError("user id is required")
	}

	if cmd.LinkID <= 0 {
		return LinkResult{}, response.BadRequestError("link id is required")
	}

	fields := make(map[string]any)

	if cmd.Title != "" {
		fields["title"] = cmd.Title
	}

	if cmd.Description != "" {
		fields["description"] = cmd.Description
	}

	if cmd.ExpiredAt != nil {
		fields["expired_at"] = cmd.ExpiredAt
	}

	link, err := s.linkRepo.UpdateLink(ctx, cmd.UserID, cmd.LinkID, fields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LinkResult{}, response.NotFoundError("link not found")
		}

		return LinkResult{}, response.InternalServerError("error update link")
	}

	shortLink := fmt.Sprintf("%s/%s", s.cfg.ShortenerBaseURL, link.CodeURL)

	result := LinkResult{
		ID:          link.ID,
		Title:       link.Title,
		OriginalURL: link.OriginalURL,
		Description: link.Description,
		ShortURL:    shortLink,
		ExpiredAt:   link.ExpiredAt,

		CreatedAt: link.CreatedAt,
		UpdatedAt: link.UpdatedAt,
	}

	return result, nil
}

func (s *linkService) DeleteLink(ctx context.Context, cmd DeleteLinkCommand) error {
	s.cache.InvalidateLinks(ctx, cmd.UserID)

	if cmd.UserID <= 0 {
		return response.BadRequestError("user id is required")
	}

	if cmd.LinkID <= 0 {
		return response.BadRequestError("link id is required")
	}

	err := s.linkRepo.DeleteLink(ctx, cmd.UserID, cmd.LinkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFoundError("link not found")
		}
		return response.InternalServerError("error delete link")
	}
	return nil
}

func (s *linkService) RedirectLink(ctx context.Context, cmd RedirectLinkCommand) (string, error) {
	var linkID int64
	var originalURL string

	cmd.Code = strings.TrimSpace(cmd.Code)
	if !isValidShortCode(cmd.Code) {
		return "", response.BadRequestError("invalid code")
	}

	cache, err := s.cache.GetRedirectLink(ctx, cmd.Code)
	if err == nil {
		originalURL = cache.OriginalURL
		linkID = cache.ID
	}

	if linkID <= 0 || originalURL == "" {
		link, err := s.linkRepo.GetLinkByCode(ctx, cmd.Code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", response.NotFoundError("link not found")
			}
			return "", response.InternalServerError(err.Error())
		}

		s.cache.SaveRedirectLink(ctx, cmd.Code, link)

		linkID = link.ID
		originalURL = link.OriginalURL
	}

	go func(linkID int64, ipAddr, userAgent, referer, device string) {
		trackCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := s.linkRepo.IncrementClickCount(trackCtx, linkID); err != nil {
			s.logger.Error("error increment click count", err, nil)
			return
		}

		if s.cache.IsUniqueClick(trackCtx, linkID, ipAddr, userAgent, referer, device) {
			return
		}

		lc := LinkClicks{
			ID:        s.idGen.NewID(),
			LinkID:    linkID,
			IpAddress: ipAddr,
			UserAgent: userAgent,
			Referer:   referer,
			Device:    device,
		}

		if err := s.linkRepo.TrackClickURL(trackCtx, lc); err != nil {
			s.logger.Error("error track click url", err, map[string]any{
				"link_id": lc.LinkID,
				"code":    cmd.Code,
			})

			return
		}

		if err := s.linkRepo.IncrementUniqueClick(trackCtx, linkID); err != nil {
			s.logger.Error("error increment unique click count", err, nil)
			return
		}

		s.cache.SaveUniqueClick(trackCtx, linkID, ipAddr, userAgent, referer, device)
	}(linkID, cmd.IpAddress, cmd.UserAgent, cmd.Referer, cmd.Device)

	return originalURL, nil
}

func (s *linkService) GetMontlyLimit(ctx context.Context, userID int64) (response.GetLinkLimit, error) {
	linkLimit := 0

	pl, err := s.subsRepo.GetCurrentPlanByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			defaultPlan, _ := s.planRepo.GetPlanByCode(ctx, "free")

			linkLimit = int(defaultPlan.MonthlyLimit)

		} else {
			return response.GetLinkLimit{}, response.InternalServerError(err.Error())
		}
	} else {
		linkLimit = int(pl.MonthlyLimit)
	}

	dateRange := utils.GetDateRange(time.Now())
	currLimit, err := s.linkRepo.GetCurrentLimit(ctx, userID, dateRange.From, dateRange.To)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return response.GetLinkLimit{}, response.InternalServerError(err.Error())
		}
	}

	return response.GetLinkLimit{
		MonthlyLimit: linkLimit,
		CurrLimit:    int(currLimit),
	}, nil
}
