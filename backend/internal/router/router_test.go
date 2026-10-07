package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// TestSetupRegistersMergeRoutes ensures the merge endpoints coexist with the
// /plants/:id wildcard route without gin tree conflicts.
func TestSetupRegistersMergeRoutes(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", RateLimitReq: 1000, RateLimitWin: 60}
	var db *gorm.DB // routes do not touch the DB during registration
	r := Setup(cfg, db, util.NewLogger())

	routes := r.Routes()
	want := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/v1/plants/merge/precheck"},
		{http.MethodPost, "/api/v1/plants/merge"},
		{http.MethodGet, "/api/v1/plants/merge/tasks"},
		{http.MethodGet, "/api/v1/plants/merge/tasks/:taskId"},
	}
	for _, w := range want {
		found := false
		for _, ri := range routes {
			if ri.Method == w.method && ri.Path == w.path {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("route %s %s not registered", w.method, w.path)
		}
	}

	// An unauthenticated request to a merge endpoint must hit the merge
	// handler chain (401 from auth middleware), not the /plants/:id handler.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/plants/merge/tasks", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("GET /plants/merge/tasks without token: status = %d, want 401", rec.Code)
	}
}
