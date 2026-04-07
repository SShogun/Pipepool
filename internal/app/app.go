package app

import (
	"Pipepool/internal/pipeline"
	"Pipepool/internal/pool"
	. "Pipepool/internal/types"
	"context"
	"log/slog"
)

func Run(ctx context.Context, cfg *Config, logger *slog.Logger) (Summary, error) {
	jobs := make(chan Job)
	results := make(chan Result)
	queuedJobs := pipeline.Run(ctx, jobs, cfg.QueueSize)

	go pool.Run(ctx, queuedJobs, results)
	summary := Summary{}.Collect(ctx, results)

	return summary.(Summary), nil
}
