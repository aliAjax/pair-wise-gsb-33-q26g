package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/gbplantwiki/gbplantwiki/internal/model"
)

// DiseasePestRepository handles persistence of disease/pest entries.
type DiseasePestRepository struct {
	db *gorm.DB
}

// NewDiseasePestRepository creates a DiseasePestRepository.
func NewDiseasePestRepository(db *gorm.DB) *DiseasePestRepository {
	return &DiseasePestRepository{db: db}
}

// Create inserts a disease/pest entry.
func (r *DiseasePestRepository) Create(d *model.DiseasePest) error {
	return r.db.Create(d).Error
}

// FindByID locates an entry by id.
func (r *DiseasePestRepository) FindByID(id uint) (*model.DiseasePest, error) {
	var d model.DiseasePest
	if err := r.db.First(&d, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &d, nil
}

// Update persists an entry.
func (r *DiseasePestRepository) Update(d *model.DiseasePest) error {
	return r.db.Save(d).Error
}

// Delete removes an entry.
func (r *DiseasePestRepository) Delete(id uint) error {
	res := r.db.Delete(&model.DiseasePest{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List filters entries by plant species and keyword with pagination.
func (r *DiseasePestRepository) List(plantSpeciesID uint, keyword string, page, pageSize int) ([]model.DiseasePest, int64, error) {
	var items []model.DiseasePest
	var total int64
	q := r.db.Model(&model.DiseasePest{})
	if plantSpeciesID > 0 {
		q = q.Where("plant_species_id = ?", plantSpeciesID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR keywords LIKE ? OR symptoms LIKE ?", like, like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListByPlant returns all entries related to a plant species.
func (r *DiseasePestRepository) ListByPlant(plantSpeciesID uint) ([]model.DiseasePest, error) {
	var items []model.DiseasePest
	if err := r.db.Where("plant_species_id = ?", plantSpeciesID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
