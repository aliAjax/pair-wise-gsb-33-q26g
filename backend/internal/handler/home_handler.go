package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/dto"
	"github.com/gbplantwiki/gbplantwiki/internal/service"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// HomeHandler aggregates home page data.
type HomeHandler struct {
	plantService   *service.PlantSpeciesService
	articleService *service.CareArticleService
}

// NewHomeHandler creates a HomeHandler.
func NewHomeHandler(plantService *service.PlantSpeciesService, articleService *service.CareArticleService) *HomeHandler {
	return &HomeHandler{plantService: plantService, articleService: articleService}
}

// Overview handles GET /home/overview.
func (h *HomeHandler) Overview(c *gin.Context) {
	hotPlants, err := h.plantService.ListHot(8)
	if err != nil {
		c.Error(err)
		return
	}
	latestArticles, err := h.articleService.ListLatest(6)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{
		"hot_plants":     hotPlants,
		"latest_articles": latestArticles,
		"season_task":    util.CurrentSeasonTask(),
	}))
}
