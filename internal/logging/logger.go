package logging

import (
	"log/slog"
	"os"
)

func LoadLogger() *slog.Logger {
	var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}
