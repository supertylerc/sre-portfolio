package internal

import (
	"os"

	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/viper"
)

func LogConfig() {
	loggerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, loggerOpts)

	if viper.Get("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
