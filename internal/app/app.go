package app

import (
	"Pipepool/internal/pipeline"
	"Pipepool/internal/pool"
	summarypkg "Pipepool/internal/summary"
	. "Pipepool/internal/types"
	"context"
	"log/slog"
)

func Run(ctx context.Context, cfg *Config, logger *slog.Logger) (Summary, error) {
	jobs := make(chan Job)
	queuedJobs := pipeline.Run(ctx, jobs, cfg.QueueSize)
	results := pool.Run(ctx, queuedJobs, cfg)
	summary := summarypkg.Collect(ctx, results)

	return summary, nil
}
