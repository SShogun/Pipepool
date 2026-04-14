package pipeline

import (
	"context"
	"io"
	"log/slog"
)

func Run(ctx context.Context, inputs []string, queueSize int, logger *slog.Logger) <-chan Item {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	ingested := ingest(ctx, inputs, logger)
	normalized := normalize(ctx, ingested, logger)
	validated := validate(ctx, normalized, logger)

	return queue(ctx, validated, queueSize, logger)
}
