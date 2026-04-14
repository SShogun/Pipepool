package pipeline

import (
	"Pipepool/internal/types"
	"context"
	"io"
	"log/slog"
	"strings"
)

func normalize(ctx context.Context, jobs <-chan types.Job, logger *slog.Logger) <-chan Item {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	out := make(chan Item)

	go func() {
		defer close(out)
		logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "normalize", "state", "start")

		for {
			select {
			case <-ctx.Done():
				logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "normalize", "state", "canceled")
				return
			case job, ok := <-jobs:
				if !ok {
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "normalize", "state", "stop")
					return
				}

				clean := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(job.Input, "\r\n", "\n"), "\r", "\n"))
				item := Item{
					ID:    job.ID,
					Input: clean,
				}

				select {
				case <-ctx.Done():
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "normalize", "state", "canceled")
					return
				case out <- item:
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "normalize", "state", "done", "job_id", item.ID)
				}
			}
		}
	}()

	return out
}
