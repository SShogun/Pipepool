package pool

import (
	"Pipepool/internal/pipeline"
	"Pipepool/internal/types"
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type processorFunc func(ctx context.Context, item pipeline.Item, workerID int) types.Result

const (
	defaultWorkerCount   = 1
	defaultQueueSize     = 1
	defaultPerJobTimeout = 250 * time.Millisecond
	defaultProcessDelay  = 25 * time.Millisecond
)

func Run(ctx context.Context, jobs <-chan pipeline.Item, cfg *types.Config, logger *slog.Logger) <-chan types.Result {
	return runWithProcessor(ctx, jobs, cfg, logger, processOne)
}

func runWithProcessor(ctx context.Context, jobs <-chan pipeline.Item, cfg *types.Config, logger *slog.Logger, process processorFunc) <-chan types.Result {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if process == nil {
		process = processOne
	}

	workerCount := defaultWorkerCount
	queueSize := defaultQueueSize
	perJobTimeout := defaultPerJobTimeout
	if cfg != nil {
		if cfg.WorkerCount > 0 {
			workerCount = cfg.WorkerCount
		}
		if cfg.QueueSize > 0 {
			queueSize = cfg.QueueSize
		}
		if cfg.PerJobTimeout > 0 {
			perJobTimeout = cfg.PerJobTimeout
		}
	}

	var wg sync.WaitGroup
	results := make(chan types.Result, queueSize)
	logger.InfoContext(ctx, "pool lifecycle", "component", "pool", "state", "start", "worker_count", workerCount)

	for i := 0; i < workerCount; i++ {
		workerID := i + 1
		wg.Add(1)

		go func() {
			defer wg.Done()
			logger.InfoContext(ctx, "pool lifecycle", "component", "pool", "state", "worker_started", "worker_id", workerID)
			defer logger.InfoContext(ctx, "pool lifecycle", "component", "pool", "state", "worker_stopped", "worker_id", workerID)

			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-jobs:
					if !ok {
						return
					}

					logger.InfoContext(ctx, "pool lifecycle", "component", "pool", "state", "worker_received", "worker_id", workerID, "job_id", item.ID)

					jobCtx, cancel := context.WithTimeout(ctx, perJobTimeout)
					result := process(jobCtx, item, workerID)
					cancel()

					select {
					case <-ctx.Done():
						return
					case results <- result:
						logger.InfoContext(ctx, "pool lifecycle", "component", "pool", "state", "worker_finished", "worker_id", workerID, "job_id", result.ID, "duration", result.Duration)
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
		logger.InfoContext(ctx, "pool lifecycle", "component", "pool", "state", "stop")
	}()

	return results
}

func processOne(ctx context.Context, item pipeline.Item, workerID int) types.Result {
	_ = workerID
	start := time.Now()

	if !item.Valid {
		return types.Result{
			ID:       item.ID,
			Valid:    item.Valid,
			Duration: time.Since(start),
		}
	}

	select {
	case <-ctx.Done():
		return types.Result{
			ID:       item.ID,
			Valid:    item.Valid,
			Duration: time.Since(start),
			Err:      ctx.Err(),
		}
	case <-time.After(defaultProcessDelay):
		return types.Result{
			ID:        item.ID,
			Output:    item.Input + "_processed",
			Valid:     item.Valid,
			WordCount: len(strings.Fields(item.Input)),
			LineCount: countLines(item.Input),
			Duration:  time.Since(start),
		}
	}
}

func countLines(input string) int {
	if input == "" {
		return 0
	}

	return strings.Count(input, "\n") + 1
}
