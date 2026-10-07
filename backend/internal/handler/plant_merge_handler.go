package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/dto"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/service"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// PlantMergeHandler exposes plant species merge endpoints (admin only).
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
	pre, err := h.svc.Precheck(req.KeepID, req.SourceIDs)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(pre))
}

// Merge handles POST /plants/merge (admin). Submitting the same batch twice
// returns the existing task, so a second window sees the completed result.
func (h *PlantMergeHandler) Merge(c *gin.Context) {
	var req dto.PlantMergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	task, err := h.svc.Merge(middleware.GetUserID(c), req.KeepID, req.SourceIDs)
	if err != nil {
		c.Error(err)
		return
	}
	view := dto.NewMergeTaskView(task)
	c.JSON(http.StatusOK, dto.OK(gin.H{
		"task":    view,
		"message": mergeStatusMessage(task),
	}))
}

// GetTask handles GET /plants/merge/tasks/:taskId (admin) for status polling.
func (h *PlantMergeHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid merge task id"))
		return
	}
	task, err := h.svc.GetTask(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.NewMergeTaskView(task)))
}

// ListTasks handles GET /plants/merge/tasks (admin).
func (h *PlantMergeHandler) ListTasks(c *gin.Context) {
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
	views := make([]dto.MergeTaskView, 0, len(items))
	for i := range items {
		views = append(views, dto.NewMergeTaskView(&items[i]))
	}
	c.JSON(http.StatusOK, dto.OK(dto.PageData{List: views, Total: total, Page: page, Size: pageSize}))
}

// mergeStatusMessage maps the task status to a user-facing hint.
func mergeStatusMessage(task *model.PlantMergeTask) string {
	switch task.Status {
	case model.MergeTaskSucceeded:
		return constants.MsgPlantMergeOK
	case model.MergeTaskProcessing:
		return constants.MsgPlantMergeProcessing
	case model.MergeTaskFailed:
		return constants.MsgPlantMergeFailed
	}
	return constants.MsgPlantMergeProcessing
}
