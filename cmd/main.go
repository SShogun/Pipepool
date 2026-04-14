package main

import (
	"Pipepool/internal/app"
	"Pipepool/internal/logging"
	"context"
	"time"
)

func main() {
	cfg := app.LoadConfig(4, 2, 5*time.Second, 500*time.Millisecond)
	logger := logging.LoadLogger()
	ctx := context.Background()
	inputs := []string{
		"  alpha beta  ",
		"",
		"line one\r\nline two",
		"go concurrency is fun",
	}

	summary, err := app.Run(ctx, cfg, inputs, logger)
	if err != nil {
		logger.Error("Error running app", "error", err)
		return
	}
	logger.Info("App run completed", "summary", summary)
}
