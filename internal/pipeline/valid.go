package pipeline

import (
	"context"
	"io"
	"log/slog"
)

func validate(ctx context.Context, items <-chan Item, logger *slog.Logger) <-chan Item {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	out := make(chan Item)

	go func() {
		defer close(out)
		logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "validate", "state", "start")

		for {
			select {
			case <-ctx.Done():
				logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "validate", "state", "canceled")
				return
			case item, ok := <-items:
				if !ok {
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "validate", "state", "stop")
					return
				}

				item.Valid = item.Input != ""

				select {
				case <-ctx.Done():
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "validate", "state", "canceled")
					return
				case out <- item:
					logger.InfoContext(ctx, "pipeline lifecycle", "component", "pipeline", "stage", "validate", "state", "done", "job_id", item.ID, "valid", item.Valid)
				}
			}
		}
	}()

	return out
}
