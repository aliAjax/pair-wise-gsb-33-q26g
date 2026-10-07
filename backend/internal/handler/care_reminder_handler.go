package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/dto"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/service"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// CareReminderHandler exposes care reminder endpoints.
type CareReminderHandler struct {
	svc    *service.CareReminderService
	logger *slog.Logger
}

// NewCareReminderHandler creates a CareReminderHandler.
func NewCareReminderHandler(svc *service.CareReminderService, logger *slog.Logger) *CareReminderHandler {
	return &CareReminderHandler{svc: svc, logger: logger}
}

// List handles GET /reminders.
func (h *CareReminderHandler) List(c *gin.Context) {
	status := c.Query("status")
	items, err := h.svc.ListByUser(middleware.GetUserID(c), status)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// ListByMonth handles GET /reminders/calendar?year=&month=.
func (h *CareReminderHandler) ListByMonth(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))
	items, err := h.svc.ListByMonth(middleware.GetUserID(c), year, month)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// Create handles POST /reminders.
func (h *CareReminderHandler) Create(c *gin.Context) {
	var req dto.ReminderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	m := &model.CareReminder{
		PlantSpeciesID: req.PlantSpeciesID, TaskTitle: req.TaskTitle,
		RemindDate: req.RemindDate, Frequency: req.Frequency,
	}
	created, err := h.svc.Create(middleware.GetUserID(c), m)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// UpdateStatus handles PUT /reminders/:id/status.
func (h *CareReminderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid reminder id"))
		return
	}
	var req dto.ReminderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	m, err := h.svc.UpdateStatus(middleware.GetUserID(c), uint(id), req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(m))
}

// Delete handles DELETE /reminders/:id.
func (h *CareReminderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid reminder id"))
		return
	}
	if err := h.svc.Delete(middleware.GetUserID(c), uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"deleted": true}))
}
