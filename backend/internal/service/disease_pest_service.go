package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/repository"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// DiseasePestService implements disease/pest manual business logic.
type DiseasePestService struct {
	repo   *repository.DiseasePestRepository
	logger *slog.Logger
}

// NewDiseasePestService creates a DiseasePestService.
func NewDiseasePestService(repo *repository.DiseasePestRepository, logger *slog.Logger) *DiseasePestService {
	return &DiseasePestService{repo: repo, logger: logger}
}

// Create adds an entry (admin only).
func (s *DiseasePestService) Create(d *model.DiseasePest) (*model.DiseasePest, error) {
	if d.Images == "" {
		d.Images = "[]"
	}
	if err := s.repo.Create(d); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogPestCreateFailed, d.Name), "error", err)
		return nil, fmt.Errorf("disease pest create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPestCreateSuccess, d.Name), "id", d.ID)
	return d, nil
}

// Get returns an entry by id.
func (s *DiseasePestService) Get(id uint) (*model.DiseasePest, error) {
	d, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("DiseasePest[id=%d] not found", id))
		}
		return nil, fmt.Errorf("disease pest get: %w", err)
	}
	return d, nil
}

// Update edits an entry (admin only).
func (s *DiseasePestService) Update(id uint, d *model.DiseasePest) (*model.DiseasePest, error) {
	exist, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("disease pest update find: %w", err)
	}
	if d.Name != "" {
		exist.Name = d.Name
	}
	if d.Symptoms != "" {
		exist.Symptoms = d.Symptoms
	}
	if d.Cause != "" {
		exist.Cause = d.Cause
	}
	if d.Treatment != "" {
		exist.Treatment = d.Treatment
	}
	if d.RecommendedMedicine != "" {
		exist.RecommendedMedicine = d.RecommendedMedicine
	}
	if d.Keywords != "" {
		exist.Keywords = d.Keywords
	}
	if err := s.repo.Update(exist); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogPestUpdateFailed, id), "error", err)
		return nil, fmt.Errorf("disease pest update: %w", err)
	}
	return exist, nil
}

// Delete removes an entry (admin only).
func (s *DiseasePestService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("disease pest delete: %w", err)
	}
	return nil
}

// List searches entries by plant species and keyword.
func (s *DiseasePestService) List(plantSpeciesID uint, keyword string, page, pageSize int) ([]model.DiseasePest, int64, error) {
	items, total, err := s.repo.List(plantSpeciesID, keyword, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("disease pest list: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPestSearchSuccess, keyword), "total", total)
	return items, total, nil
}

// ListByPlant returns entries for a plant species detail page.
func (s *DiseasePestService) ListByPlant(plantSpeciesID uint) ([]model.DiseasePest, error) {
	items, err := s.repo.ListByPlant(plantSpeciesID)
	if err != nil {
		return nil, fmt.Errorf("disease pest by plant: %w", err)
	}
	return items, nil
}
