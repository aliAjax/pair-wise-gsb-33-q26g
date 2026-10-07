package util

import (
	"log/slog"
	"os"
	"time"
)

// NewLogger builds a structured JSON slog logger writing to stdout.
func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// LoggerWithRequest returns a logger pre-bound with request metadata.
func LoggerWithRequest(logger *slog.Logger, requestID, method, path string) *slog.Logger {
	return logger.With(
		"request_id", requestID,
		"method", method,
		"path", path,
	)
}

// RecordLatency measures and logs the latency of a request.
func RecordLatency(start time.Time, logger *slog.Logger, status int) {
	logger.Info("request latency",
		"status", status,
		"latency_ms", time.Since(start).Milliseconds(),
	)
}
