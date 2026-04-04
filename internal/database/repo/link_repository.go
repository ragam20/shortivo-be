package repo

import (
	"context"
	"time"

	"shortivo/internal/app/link"
	"shortivo/internal/database/entity"

	"gorm.io/gorm"
)

type linkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) link.LinkRepository {
	return &linkRepository{db: db}
}

func (r *linkRepository) CreateLink(ctx context.Context, domain link.Link) error {
	link := entity.NewLink(domain)
	if err := r.db.WithContext(ctx).Create(&link).Error; err != nil {
		return err
	}

	return nil
}

func (r *linkRepository) GetTotalItemByUser(ctx context.Context, userID int64) int64 {
	var totalItem int64
	r.db.WithContext(ctx).Model(&entity.LinkEntity{}).Where("user_id = ?", userID).Count(&totalItem)

	return totalItem
}

func (r *linkRepository) GetLinks(ctx context.Context, userID int64, offset int, limit int) []link.Link {
	var entities []entity.LinkEntity
	r.db.WithContext(ctx).Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).Order("created_at DESC").Find(&entities)

	result := make([]link.Link, len(entities))
	for i, v := range entities {
		result[i] = v.ToDomain()
	}

	return result
}

func (r *linkRepository) UpdateLink(ctx context.Context, userID int64, linkID int64, fields map[string]any) (link.Link, error) {
	tx := r.db.WithContext(ctx).Begin()

	result := tx.Model(&entity.LinkEntity{}).
		Where("id = ? AND user_id = ?", linkID, userID).
		Updates(fields)

	if result.Error != nil {
		tx.Rollback()
		return link.Link{}, result.Error
	}

	var e entity.LinkEntity
	if err := tx.
		Where("id = ? AND user_id = ?", linkID, userID).
		First(&e).Error; err != nil {

		tx.Rollback()
		return link.Link{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return link.Link{}, err
	}

	return e.ToDomain(), nil
}

func (r *linkRepository) DeleteLink(ctx context.Context, userID int64, linkID int64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", linkID, userID).
		Delete(&entity.LinkEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *linkRepository) GetLinkByCode(ctx context.Context, code string) (link.Link, error) {
	var entity entity.LinkEntity
	if err := r.db.WithContext(ctx).Where("code_url = ?", code).First(&entity).Error; err != nil {
		return link.Link{}, err
	}

	return entity.ToDomain(), nil
}

func (r *linkRepository) TrackClickURL(ctx context.Context, domain link.LinkClicks) error {
	linkClicks := entity.NewLinkClicks(domain)
	if err := r.db.WithContext(ctx).Create(&linkClicks).Error; err != nil {
		return err
	}
	return nil
}

func (r *linkRepository) IncrementClickCount(ctx context.Context, linkID int64) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.LinkEntity{}).
		Where("id = ?", linkID).
		Update("click_count", gorm.Expr("click_count + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

func (r *linkRepository) IncrementUniqueClick(ctx context.Context, linkID int64) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.LinkEntity{}).
		Where("id = ?", linkID).
		Update("unique_click", gorm.Expr("unique_click + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

func (r *linkRepository) GetCurrentLimit(ctx context.Context, userID int64, startAt time.Time, endAt time.Time) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.LinkEntity{}).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startAt, endAt).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
