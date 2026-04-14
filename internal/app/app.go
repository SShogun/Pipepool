package app

import (
	"Pipepool/internal/pipeline"
	"Pipepool/internal/pool"
	summarypkg "Pipepool/internal/summary"
	"Pipepool/internal/types"
	"context"
	"io"
	"log/slog"
)

func Run(ctx context.Context, cfg *types.Config, inputs []string, logger *slog.Logger) (types.Summary, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	effectiveConfig := normalizeConfig(types.Config{})
	if cfg != nil {
		effectiveConfig = normalizeConfig(*cfg)
	}

	runCtx, cancel := context.WithTimeout(ctx, effectiveConfig.RunTimeout)
	defer cancel()

	logger.InfoContext(runCtx, "pipepool lifecycle", "component", "app", "state", "start", "worker_count", effectiveConfig.WorkerCount, "queue_size", effectiveConfig.QueueSize, "inputs", len(inputs))

	queuedJobs := pipeline.Run(runCtx, inputs, effectiveConfig.QueueSize, logger)
	results := pool.Run(runCtx, queuedJobs, &effectiveConfig, logger)
	summary := summarypkg.Collect(runCtx, results, logger)

	if err := runCtx.Err(); err != nil {
		logger.WarnContext(ctx, "pipepool lifecycle", "component", "app", "state", "canceled", "error", err)
		return summary, err
	}

	logger.InfoContext(runCtx, "pipepool lifecycle", "component", "app", "state", "stop", "total_jobs", summary.TotalJobs, "successes", summary.SuccessCount, "failures", summary.FailureCount)
	return summary, nil
}
