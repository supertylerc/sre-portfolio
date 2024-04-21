package internal

import (
	"log/slog"
	"os"
)

func LogConfig(c *Config) {
	loggerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, loggerOpts)

	if c.Environment == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
