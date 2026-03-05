package link

import (
	"context"
	"encoding/json"
	"fmt"
	"shortivo/internal/cache"
	"time"
)

type LinkCache interface {
	GetLinks(ctx context.Context, userID int64, page, limit int) (GetLinkResult, error)
	SaveListLinks(ctx context.Context, userID int64, page, limit int, links GetLinkResult) error
	GetRedirectLink(ctx context.Context, code string) (Link, error)
	SaveRedirectLink(ctx context.Context, code string, link Link) error
	InvalidateLinks(ctx context.Context, userID int64)
	SaveUniqueClick(ctx context.Context, linkID int64, ipAddr, userAgent, referer, device string) error
	IsUniqueClick(ctx context.Context, linkID int64, ipAddr, userAgent, referer, device string) bool
}

type linkCache struct {
	store cache.Store
}

func NewLinkCache(store cache.Store) LinkCache {
	return &linkCache{store: store}
}

func (c *linkCache) GetLinks(ctx context.Context, userID int64, page, limit int) (GetLinkResult, error) {
	value, err := c.store.Get(ctx, fmt.Sprintf("links:%d:%d:%d", userID, page, limit))
	if err != nil {
		return GetLinkResult{}, err
	}

	var result GetLinkResult
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *linkCache) SaveListLinks(ctx context.Context, userID int64, page, limit int, links GetLinkResult) error {
	data, err := json.Marshal(links)
	if err != nil {
		return err
	}

	err = c.store.Set(ctx, fmt.Sprintf("links:%d:%d:%d", userID, page, limit), string(data), 20*time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (c *linkCache) InvalidateLinks(ctx context.Context, userID int64) {
	c.store.Del(ctx, fmt.Sprintf("links:%d:*", userID))
}

func (c *linkCache) GetRedirectLink(ctx context.Context, code string) (Link, error) {
	value, err := c.store.Get(ctx, fmt.Sprintf("links:%s", code))
	if err != nil {
		return Link{}, err
	}

	var link Link
	err = json.Unmarshal([]byte(value), &link)
	if err != nil {
		return Link{}, err
	}

	return link, nil
}

func (c *linkCache) SaveRedirectLink(ctx context.Context, code string, link Link) error {
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	err = c.store.Set(ctx, fmt.Sprintf("links:%s", code), string(data), 168*time.Hour)
	if err != nil {
		return err
	}

	return nil
}

func (c *linkCache) SaveUniqueClick(ctx context.Context, linkID int64, ipAddr, userAgent, referer, device string) error {
	uniqueUser := fmt.Sprintf("link:%d:click:%x", linkID, visitorHash(ipAddr, userAgent, referer, device))
	err := c.store.Set(ctx, uniqueUser, "1", 24*time.Hour)
	if err != nil {
		return err
	}

	return err
}

func (c *linkCache) IsUniqueClick(ctx context.Context, linkID int64, ipAddr, userAgent, referer, device string) bool {
	_, err := c.store.Get(ctx, fmt.Sprintf("link:%d:click:%x", linkID, visitorHash(ipAddr, userAgent, referer, device)))
	return err == nil
}
