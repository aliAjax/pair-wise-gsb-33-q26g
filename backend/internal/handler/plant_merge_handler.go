package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/dto"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
	"github.com/gbplantwiki/gbplantwiki/internal/service"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// PlantMergeHandler exposes variety merge endpoints (admin only).
type PlantMergeHandler struct {
	svc    *service.PlantMergeService
	logger *slog.Logger
}

// NewPlantMergeHandler creates a PlantMergeHandler.
func NewPlantMergeHandler(svc *service.PlantMergeService, logger *slog.Logger) *PlantMergeHandler {
	return &PlantMergeHandler{svc: svc, logger: logger}
}

// Precheck handles POST /plants/merge/precheck (admin).
func (h *PlantMergeHandler) Precheck(c *gin.Context) {
	var req dto.PlantMergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	result, err := h.svc.Precheck(req.KeepID, req.SourceIDs)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(result))
}

// Merge handles POST /plants/merge (admin). Resubmitting the same batch is
// idempotent: finished pairs replay their result, failed pairs resume.
func (h *PlantMergeHandler) Merge(c *gin.Context) {
	var req dto.PlantMergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	results, err := h.svc.Merge(middleware.GetUserID(c), req.KeepID, req.SourceIDs)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(results))
}

// List handles GET /plants/merges (admin).
func (h *PlantMergeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := h.svc.ListTasks(page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PageData{List: items, Total: total, Page: page, Size: pageSize}))
}

// Get handles GET /plants/merges/:id (admin).
func (h *PlantMergeHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid merge task id"))
		return
	}
	result, err := h.svc.GetTask(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(result))
}
