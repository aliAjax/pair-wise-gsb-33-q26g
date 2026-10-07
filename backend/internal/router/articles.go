package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/handler"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
)

func registerArticleRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.CareArticleHandler, limiter *middleware.RateLimiter) {
	articles := v1.Group("/articles")
	articles.GET("", h.List)
	articles.GET("/:id", h.Get)
	auth := articles.Group("", middleware.AuthRequired(cfg))
	auth.POST("", limiter.Limit(), h.Create)
	auth.PUT("/:id", h.Update)
	auth.DELETE("/:id", h.Delete)
}
