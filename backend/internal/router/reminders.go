package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/handler"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
)

func registerReminderRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.CareReminderHandler, limiter *middleware.RateLimiter) {
	reminders := v1.Group("/reminders", middleware.AuthRequired(cfg))
	reminders.GET("", h.List)
	reminders.GET("/calendar", h.ListByMonth)
	reminders.POST("", limiter.Limit(), h.Create)
	reminders.PUT("/:id/status", h.UpdateStatus)
	reminders.DELETE("/:id", h.Delete)
}
