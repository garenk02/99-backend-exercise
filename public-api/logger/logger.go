package logger

import (
	"log/slog"
	"os"
)

// Setup initializes the logger with JSON formatting
func Setup() {
	// Configure slog with JSON handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	
	// Set as default logger
	slog.SetDefault(logger)
}

// Get returns the default logger
func Get() *slog.Logger {
	return slog.Default()
}