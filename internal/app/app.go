package app

import (
	"Pipepool/internal/pipeline"
	"Pipepool/internal/pool"
	. "Pipepool/internal/types"
	"context"
	"log/slog"
)

func Run(ctx context.Context, cfg *Config, logger *slog.Logger) (Summary, error) {
	jobs := make(chan Job, cfg.QueueSize)
	results := make(chan Result)

	go pipeline.Run(ctx, jobs)
	go pipeline.Run(ctx, jobs)
	go pool.Run(ctx, jobs, results)
	summary := Summary{}.Collect(ctx, results)

	return summary.(Summary), nil
}
