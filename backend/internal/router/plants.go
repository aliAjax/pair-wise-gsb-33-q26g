package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/handler"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
)

func registerPlantRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.PlantSpeciesHandler, mergeHandler *handler.PlantMergeHandler, limiter *middleware.RateLimiter) {
	plants := v1.Group("/plants")
	plants.GET("", h.List)
	plants.GET("/:id", h.Get)
	admin := plants.Group("", middleware.AuthRequired(cfg), middleware.RequireRole("admin"))
	admin.POST("", limiter.Limit(), h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
	// Variety merge: static segments take precedence over /:id in gin.
	admin.POST("/merge/precheck", limiter.Limit(), mergeHandler.Precheck)
	admin.POST("/merge", limiter.Limit(), mergeHandler.Merge)
	admin.GET("/merges", mergeHandler.List)
	admin.GET("/merges/:id", mergeHandler.Get)
}
