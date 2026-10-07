package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/handler"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
)

func registerQuestionRoutes(v1 *gin.RouterGroup, cfg *config.Config, qh *handler.QuestionHandler, ah *handler.AnswerHandler, limiter *middleware.RateLimiter) {
	questions := v1.Group("/questions")
	questions.GET("", qh.List)
	questions.GET("/:id", qh.Get)
	questions.GET("/:id/answers", ah.List)
	auth := questions.Group("", middleware.AuthRequired(cfg))
	auth.POST("", limiter.Limit(), qh.Create)
	auth.POST("/:id/answers", ah.Create)
	auth.PUT("/:id/adopt", ah.Adopt)
}
