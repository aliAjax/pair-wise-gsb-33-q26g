package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DB_HOST", "")
	cfg := Load()
	if cfg.ServerPort != "8080" {
		t.Errorf("ServerPort = %s, want 8080", cfg.ServerPort)
	}
	if cfg.DBName != "gbplantwiki_db" {
		t.Errorf("DBName = %s, want gbplantwiki_db", cfg.DBName)
	}
	if cfg.RateLimitReq <= 0 {
		t.Errorf("RateLimitReq should be positive, got %d", cfg.RateLimitReq)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("DB_HOST", "mysql-test")
	t.Setenv("SERVER_PORT", "9090")
	cfg := Load()
	if cfg.DBHost != "mysql-test" {
		t.Errorf("DBHost = %s, want mysql-test", cfg.DBHost)
	}
	if cfg.ServerPort != "9090" {
		t.Errorf("ServerPort = %s, want 9090", cfg.ServerPort)
	}
	if got := cfg.DSN(); got == "" {
		t.Error("DSN should not be empty")
	}
}
