package pool

import (
	"Pipepool/internal/pipeline"
	. "Pipepool/internal/types"
	"context"
	"strings"
	"sync"
	"time"
)

func Run(ctx context.Context, jobs <-chan pipeline.Item, cfg *Config) <-chan Result {
	var wg sync.WaitGroup
	results := make(chan Result, cfg.QueueSize)

	for i := 0; i < cfg.WorkerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-jobs:
					if !ok {
						return
					}

					jobCtx, cancel := context.WithTimeout(ctx, cfg.PerJobTimeout)
					result := process(jobCtx, item)
					cancel()

					select {
					case <-ctx.Done():
						return
					case results <- result:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func process(ctx context.Context, item pipeline.Item) Result {
	start := time.Now()

	if !item.Valid {
		return Result{
			ID:       item.ID,
			Valid:    item.Valid,
			Duration: time.Since(start),
		}
	}

	select {
	case <-ctx.Done():
		return Result{
			ID:       item.ID,
			Valid:    item.Valid,
			Duration: time.Since(start),
			Err:      ctx.Err(),
		}
	case <-time.After(10 * time.Millisecond):
		return Result{
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
