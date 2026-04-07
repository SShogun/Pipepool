package pipeline

import (
	. "Pipepool/internal/types"
	"context"
)

func ingest(ctx context.Context, jobs <-chan Job) <-chan Job {
	out := make(chan Job)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobs:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return
				case out <- job:
				}
			}
		}
	}()

	return out
}
