package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/handler"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
)

func registerFavoriteRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.FavoriteHandler, limiter *middleware.RateLimiter) {
	favorites := v1.Group("/favorites", middleware.AuthRequired(cfg))
	favorites.GET("", h.List)
	favorites.POST("", limiter.Limit(), h.Add)
	favorites.DELETE("/:targetType/:targetId", h.Remove)
}
