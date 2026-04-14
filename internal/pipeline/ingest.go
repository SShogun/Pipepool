package pipeline

import (
	"Pipepool/internal/types"
	"context"
	"io"
	"log/slog"
)

func ingest(ctx context.Context, inputs []string, logger *slog.Logger) <-chan types.Job {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	out := make(chan types.Job)

	go func() {
		defer close(out)
		logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "ingest", "state", "start")

		for i, input := range inputs {
			job := types.Job{
				ID:    i + 1,
				Input: input,
			}
			select {
			case <-ctx.Done():
				logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "ingest", "state", "canceled")
				return
			case out <- job:
				logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "ingest", "state", "emitted", "job_id", job.ID)
			}
		}

		logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "ingest", "state", "stop", "jobs", len(inputs))
	}()

	return out
}
