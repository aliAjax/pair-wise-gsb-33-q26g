package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/dto"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/service"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// PlantSpeciesHandler exposes plant species endpoints.
type PlantSpeciesHandler struct {
	svc    *service.PlantSpeciesService
	logger *slog.Logger
}

// NewPlantSpeciesHandler creates a PlantSpeciesHandler.
func NewPlantSpeciesHandler(svc *service.PlantSpeciesService, logger *slog.Logger) *PlantSpeciesHandler {
	return &PlantSpeciesHandler{svc: svc, logger: logger}
}

// List handles GET /plants.
func (h *PlantSpeciesHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	speciesType := c.Query("type")
	family := c.Query("family")
	keyword := c.Query("keyword")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}
	items, total, err := h.svc.List(speciesType, family, keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PageData{List: items, Total: total, Page: page, Size: pageSize}))
}

// Get handles GET /plants/:id.
func (h *PlantSpeciesHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid plant id"))
		return
	}
	p, err := h.svc.Get(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(p))
}

// Create handles POST /plants (admin).
func (h *PlantSpeciesHandler) Create(c *gin.Context) {
	var req dto.PlantCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	p := &model.PlantSpecies{
		Family: req.Family, Genus: req.Genus, Name: req.Name, Alias: req.Alias,
		Type: req.Type, Origin: req.Origin, TempMin: req.TempMin, TempMax: req.TempMax,
		LightRequirement: req.LightRequirement, WaterFrequency: req.WaterFrequency,
		Description: req.Description, ImageURLs: req.ImageURLs,
	}
	created, err := h.svc.Create(p)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// Update handles PUT /plants/:id (admin).
func (h *PlantSpeciesHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid plant id"))
		return
	}
	var req dto.PlantCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	p := &model.PlantSpecies{
		Family: req.Family, Genus: req.Genus, Name: req.Name, Alias: req.Alias,
		Type: req.Type, Origin: req.Origin, TempMin: req.TempMin, TempMax: req.TempMax,
		LightRequirement: req.LightRequirement, WaterFrequency: req.WaterFrequency,
		Description: req.Description, ImageURLs: req.ImageURLs,
	}
	updated, err := h.svc.Update(uint(id), p)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(updated))
}

// Delete handles DELETE /plants/:id (admin).
func (h *PlantSpeciesHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid plant id"))
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"deleted": true}))
}
