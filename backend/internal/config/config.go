package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config aggregates all runtime configuration loaded from environment variables.
type Config struct {
	ServerPort   string
	CORSOrigins  string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	JWTExpire    time.Duration
	UploadDir    string
	RateLimitReq int
	RateLimitWin time.Duration
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		CORSOrigins:  getEnv("APP_CORS_ORIGINS", "http://localhost:8102"),
		DBHost:       getEnv("DB_HOST", "db"),
		DBPort:       getEnv("DB_PORT", "3502"),
		DBUser:       getEnv("DB_USER", "gbplantwiki_user"),
		DBPassword:   getEnv("DB_PASSWORD", "gbplantwiki_pwd"),
		DBName:       getEnv("DB_NAME", "gbplantwiki_db"),
		JWTSecret:    getEnv("JWT_SECRET", "change_me_to_a_long_random_string"),
		JWTExpire:    time.Duration(getEnvInt("JWT_EXPIRE_HOURS", 72)) * time.Hour,
		UploadDir:    getEnv("UPLOAD_DIR", "/app/uploads"),
		RateLimitReq: getEnvInt("RATE_LIMIT_REQUESTS", 60),
		RateLimitWin: time.Duration(getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
	}
}

// DSN returns the MySQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
