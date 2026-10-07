package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/gbplantwiki/gbplantwiki/internal/model"
)

// FavoriteRepository handles persistence of favorites.
type FavoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository creates a FavoriteRepository.
func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

// Create inserts a favorite.
func (r *FavoriteRepository) Create(f *model.Favorite) error {
	if err := r.db.Create(f).Error; err != nil {
		if isDuplicate(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

// Find locates a favorite by user and target.
func (r *FavoriteRepository) Find(userID uint, targetType string, targetID uint) (*model.Favorite, error) {
	var f model.Favorite
	if err := r.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		First(&f).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &f, nil
}

// Delete removes a favorite by id.
func (r *FavoriteRepository) Delete(id uint) error {
	return r.db.Delete(&model.Favorite{}, id).Error
}

// ListByUser returns favorites of a user optionally filtered by target type.
func (r *FavoriteRepository) ListByUser(userID uint, targetType string) ([]model.Favorite, error) {
	var items []model.Favorite
	q := r.db.Where("user_id = ?", userID)
	if targetType != "" {
		q = q.Where("target_type = ?", targetType)
	}
	if err := q.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
