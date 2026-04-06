package app

import (
	"Pipepool/internal/pipeline"
	. "Pipepool/internal/types"
	"context"
	"log/slog"
)

func Run(ctx context.Context, cfg *Config, logger *slog.Logger) {
	jobs := make(chan Job, cfg.QueueSize)
	results := make(chan Result)

	go pipeline.Run(ctx, jobs, results)
	
}
