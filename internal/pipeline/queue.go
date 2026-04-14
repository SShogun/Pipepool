package pipeline

import (
	"context"
	"io"
	"log/slog"
	"time"
)

func queue(ctx context.Context, items <-chan Item, queueSize int, logger *slog.Logger) <-chan Item {
	if queueSize <= 0 {
		queueSize = 1
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	out := make(chan Item, queueSize)

	go func() {
		defer close(out)
		logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "enqueue", "state", "start", "queue_size", queueSize)

		for {
			select {
			case <-ctx.Done():
				logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "enqueue", "state", "canceled")
				return
			case item, ok := <-items:
				if !ok {
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "enqueue", "state", "stop")
					return
				}

				waitStart := time.Now()
				logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "enqueue", "state", "waiting", "job_id", item.ID, "queue_size", queueSize)

				select {
				case <-ctx.Done():
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "enqueue", "state", "canceled")
					return
				case out <- item:
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "enqueue", "state", "success", "job_id", item.ID, "duration", time.Since(waitStart), "queue_size", queueSize)
				}
			}
		}
	}()

	return out
}
