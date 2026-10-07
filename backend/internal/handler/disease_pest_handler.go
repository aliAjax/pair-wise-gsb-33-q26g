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

// DiseasePestHandler exposes disease/pest manual endpoints.
type DiseasePestHandler struct {
	svc    *service.DiseasePestService
	logger *slog.Logger
}

// NewDiseasePestHandler creates a DiseasePestHandler.
func NewDiseasePestHandler(svc *service.DiseasePestService, logger *slog.Logger) *DiseasePestHandler {
	return &DiseasePestHandler{svc: svc, logger: logger}
}

// List handles GET /pests.
func (h *DiseasePestHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	plantID, _ := strconv.ParseUint(c.Query("plant_species_id"), 10, 64)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := h.svc.List(uint(plantID), keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PageData{List: items, Total: total, Page: page, Size: pageSize}))
}

// Get handles GET /pests/:id.
func (h *DiseasePestHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid pest id"))
		return
	}
	d, err := h.svc.Get(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(d))
}

// Create handles POST /pests (admin).
func (h *DiseasePestHandler) Create(c *gin.Context) {
	var req dto.PestCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	d := &model.DiseasePest{
		PlantSpeciesID: req.PlantSpeciesID, Name: req.Name, Symptoms: req.Symptoms,
		Cause: req.Cause, Treatment: req.Treatment, RecommendedMedicine: req.RecommendedMedicine,
		Images: req.Images, Keywords: req.Keywords,
	}
	created, err := h.svc.Create(d)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// Update handles PUT /pests/:id (admin).
func (h *DiseasePestHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid pest id"))
		return
	}
	var req dto.PestCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	d := &model.DiseasePest{
		PlantSpeciesID: req.PlantSpeciesID, Name: req.Name, Symptoms: req.Symptoms,
		Cause: req.Cause, Treatment: req.Treatment, RecommendedMedicine: req.RecommendedMedicine,
		Images: req.Images, Keywords: req.Keywords,
	}
	updated, err := h.svc.Update(uint(id), d)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(updated))
}

// Delete handles DELETE /pests/:id (admin).
func (h *DiseasePestHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid pest id"))
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"deleted": true}))
}
