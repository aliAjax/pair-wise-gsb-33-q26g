package router

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/handler"
	"github.com/gbplantwiki/gbplantwiki/internal/middleware"
)

// TestPlantMergeRoutesCoexistWithPlantID ensures the static merge segments do
// not conflict with /plants/:id and keep their expected methods.
func TestPlantMergeRoutesCoexistWithPlantID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	cfg := &config.Config{}
	limiter := middleware.NewRateLimiter(1000, 60)
	registerPlantRoutes(v1, cfg, handler.NewPlantSpeciesHandler(nil, nil),
		handler.NewPlantMergeHandler(nil, nil), limiter)

	want := map[string]string{
		"GET /api/v1/plants/:id":             "",
		"POST /api/v1/plants/merge":          "",
		"POST /api/v1/plants/merge/precheck": "",
		"GET /api/v1/plants/merges":          "",
		"GET /api/v1/plants/merges/:id":      "",
	}
	found := map[string]bool{}
	for _, ri := range r.Routes() {
		key := ri.Method + " " + ri.Path
		if _, ok := want[key]; ok {
			found[key] = true
		}
	}
	for key := range want {
		if !found[key] {
			t.Errorf("route missing: %s", key)
		}
	}
}
