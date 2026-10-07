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

// PlantSpeciesService implements plant species business logic.
type PlantSpeciesService struct {
	repo   *repository.PlantSpeciesRepository
	logger *slog.Logger
}

// NewPlantSpeciesService creates a PlantSpeciesService.
func NewPlantSpeciesService(repo *repository.PlantSpeciesRepository, logger *slog.Logger) *PlantSpeciesService {
	return &PlantSpeciesService{repo: repo, logger: logger}
}

// Create adds a new plant species (admin only).
func (s *PlantSpeciesService) Create(p *model.PlantSpecies) (*model.PlantSpecies, error) {
	if !constants.IsValidPlantType(p.Type) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("PlantSpecies[type=%s] create failed: invalid plant type", p.Type))
	}
	if p.ImageURLs == "" {
		p.ImageURLs = "[]"
	}
	if err := s.repo.Create(p); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("PlantSpecies[name=%s] create failed: name already exists", p.Name))
		}
		s.logger.Error(fmt.Sprintf(constants.LogPlantCreateFailed, p.Name), "error", err)
		return nil, fmt.Errorf("plant species create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantCreateSuccess, p.Name, p.Type), "id", p.ID)
	return p, nil
}

// Get returns a plant species by id.
func (s *PlantSpeciesService) Get(id uint) (*model.PlantSpecies, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("PlantSpecies[id=%d] not found", id))
		}
		return nil, fmt.Errorf("plant species get: %w", err)
	}
	return p, nil
}

// Update modifies an existing plant species (admin only).
func (s *PlantSpeciesService) Update(id uint, p *model.PlantSpecies) (*model.PlantSpecies, error) {
	exist, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("plant species update find: %w", err)
	}
	if p.Type != "" && !constants.IsValidPlantType(p.Type) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("PlantSpecies[id=%d] type=%s invalid", id, p.Type))
	}
	if p.Name != "" {
		exist.Name = p.Name
	}
	if p.Type != "" {
		exist.Type = p.Type
	}
	if p.Family != "" {
		exist.Family = p.Family
	}
	if p.Genus != "" {
		exist.Genus = p.Genus
	}
	if p.Alias != "" {
		exist.Alias = p.Alias
	}
	if p.Origin != "" {
		exist.Origin = p.Origin
	}
	if p.Description != "" {
		exist.Description = p.Description
	}
	exist.TempMin = p.TempMin
	exist.TempMax = p.TempMax
	if p.LightRequirement != "" {
		exist.LightRequirement = p.LightRequirement
	}
	if p.WaterFrequency != "" {
		exist.WaterFrequency = p.WaterFrequency
	}
	if p.ImageURLs != "" {
		exist.ImageURLs = p.ImageURLs
	}
	if err := s.repo.Update(exist); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogPlantUpdateFailed, id), "error", err)
		return nil, fmt.Errorf("plant species update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantUpdateSuccess, id), "id", id)
	return exist, nil
}

// Delete removes a plant species (admin only).
func (s *PlantSpeciesService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("PlantSpecies[id=%d] not found", id))
		}
		return fmt.Errorf("plant species delete: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantDeleteSuccess, id), "id", id)
	return nil
}

// List filters plant species by type/family/keyword with pagination.
func (s *PlantSpeciesService) List(speciesType, family, keyword string, page, pageSize int) ([]model.PlantSpecies, int64, error) {
	items, total, err := s.repo.List(speciesType, family, keyword, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("plant species list: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantListSuccess, page, pageSize), "total", total)
	return items, total, nil
}

// ListHot returns recent plants for the home page.
func (s *PlantSpeciesService) ListHot(limit int) ([]model.PlantSpecies, error) {
	items, err := s.repo.ListHot(limit)
	if err != nil {
		return nil, fmt.Errorf("plant species hot: %w", err)
	}
	return items, nil
}
