package pipeline

import (
	. "Pipepool/internal/types"
	"context"
)

func Run(ctx context.Context, jobs <-chan Job, queueSize int) <-chan Item {
	ingested := ingest(ctx, jobs)
	normalized := normalize(ctx, ingested)
	validated := validate(ctx, normalized)

	return queue(ctx, validated, queueSize)
}
