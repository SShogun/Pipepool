package pipeline

import (
	. "Pipepool/internal/types"
	"context"
	"strings"
)

func normalize(ctx context.Context, jobs <-chan Job) <-chan Item {
	out := make(chan Item)

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

				clean := strings.TrimSpace(job.Input)
				item := Item{
					ID:     job.ID,
					Input:  clean,
					Valid:  false,
					Result: Result{},
				}

				select {
				case <-ctx.Done():
					return
				case out <- item:
				}
			}
		}
	}()

	return out
}
