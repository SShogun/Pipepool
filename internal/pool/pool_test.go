package pool

import (
	"Pipepool/internal/pipeline"
	"Pipepool/internal/testutil"
	"Pipepool/internal/types"
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunWorkerCountNeverExceedsConfig(t *testing.T) {
	cfg := &types.Config{
		WorkerCount:   3,
		QueueSize:     8,
		PerJobTimeout: time.Second,
	}

	const totalJobs = 24
	in := make(chan pipeline.Item, totalJobs)
	for i := 0; i < totalJobs; i++ {
		in <- pipeline.Item{ID: i + 1, Input: "alpha beta", Valid: true}
	}
	close(in)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var active int32
	var maxActive int32

	processor := func(ctx context.Context, item pipeline.Item, workerID int) types.Result {
		current := atomic.AddInt32(&active, 1)
		for {
			recorded := atomic.LoadInt32(&maxActive)
			if current <= recorded {
				break
			}
			if atomic.CompareAndSwapInt32(&maxActive, recorded, current) {
				break
			}
		}

		defer atomic.AddInt32(&active, -1)

		select {
		case <-ctx.Done():
			return types.Result{
				ID:       item.ID,
				Valid:    item.Valid,
				Duration: 0,
				Err:      ctx.Err(),
			}
		case <-time.After(35 * time.Millisecond):
			return types.Result{
				ID:        item.ID,
				Valid:     item.Valid,
				WordCount: 2,
				LineCount: 1,
				Duration:  35 * time.Millisecond,
			}
		}
	}

	results := runWithProcessor(ctx, in, cfg, testutil.NewDiscardLogger(), processor)

	count := 0
	for range results {
		count++
	}

	if count != totalJobs {
		t.Fatalf("got %d results, want %d", count, totalJobs)
	}

	if got := int(atomic.LoadInt32(&maxActive)); got > cfg.WorkerCount {
		t.Fatalf("active workers exceeded cap: got %d, cap %d", got, cfg.WorkerCount)
	}
}
