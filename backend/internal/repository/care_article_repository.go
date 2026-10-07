package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
)

// CareArticleRepository handles persistence of care articles.
type CareArticleRepository struct {
	db *gorm.DB
}

// NewCareArticleRepository creates a CareArticleRepository.
func NewCareArticleRepository(db *gorm.DB) *CareArticleRepository {
	return &CareArticleRepository{db: db}
}

// Create inserts an article.
func (r *CareArticleRepository) Create(a *model.CareArticle) error {
	return r.db.Create(a).Error
}

// FindByID locates an article by id.
func (r *CareArticleRepository) FindByID(id uint) (*model.CareArticle, error) {
	var a model.CareArticle
	if err := r.db.First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

// Update persists an article.
func (r *CareArticleRepository) Update(a *model.CareArticle) error {
	return r.db.Save(a).Error
}

// Delete removes an article by id.
func (r *CareArticleRepository) Delete(id uint) error {
	res := r.db.Delete(&model.CareArticle{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// IncrementView bumps the view count of an article.
func (r *CareArticleRepository) IncrementView(id uint) error {
	return r.db.Model(&model.CareArticle{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// List filters articles by topic tag and keyword with pagination.
func (r *CareArticleRepository) List(topicTag, keyword string, page, pageSize int) ([]model.CareArticle, int64, error) {
	var items []model.CareArticle
	var total int64
	q := r.db.Model(&model.CareArticle{}).Where("status = ?", constants.ArticleStatusPublished)
	if topicTag != "" {
		q = q.Where("topic_tag = ?", topicTag)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListLatest returns the latest published articles for the home page.
func (r *CareArticleRepository) ListLatest(limit int) ([]model.CareArticle, error) {
	var items []model.CareArticle
	if err := r.db.Where("status = ?", constants.ArticleStatusPublished).
		Order("id DESC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
