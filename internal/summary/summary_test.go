package summary

import (
	. "Pipepool/internal/types"
	"context"
	"errors"
	"testing"
	"time"
)

func TestCollectAggregatesResults(t *testing.T) {
	results := make(chan Result, 3)
	workerErr := errors.New("worker failed")

	results <- Result{ID: 1, Valid: true, WordCount: 2, LineCount: 1, Duration: 10 * time.Millisecond}
	results <- Result{ID: 2, Valid: false, Duration: 20 * time.Millisecond}
	results <- Result{ID: 3, Valid: true, WordCount: 4, LineCount: 2, Duration: 5 * time.Millisecond, Err: workerErr}
	close(results)

	summary := Collect(context.Background(), results)

	if summary.TotalJobs != 3 {
		t.Fatalf("TotalJobs = %d, want 3", summary.TotalJobs)
	}
	if summary.SuccessCount != 1 {
		t.Fatalf("SuccessCount = %d, want 1", summary.SuccessCount)
	}
	if summary.FailureCount != 2 {
		t.Fatalf("FailureCount = %d, want 2", summary.FailureCount)
	}
	if summary.TotalWords != 6 {
		t.Fatalf("TotalWords = %d, want 6", summary.TotalWords)
	}
	if summary.TotalLines != 3 {
		t.Fatalf("TotalLines = %d, want 3", summary.TotalLines)
	}
	if summary.TotalDuration != 35*time.Millisecond {
		t.Fatalf("TotalDuration = %v, want 35ms", summary.TotalDuration)
	}
	if summary.SlowestJobID != 2 {
		t.Fatalf("SlowestJobID = %d, want 2", summary.SlowestJobID)
	}
	if summary.MaxDuration != 20*time.Millisecond {
		t.Fatalf("MaxDuration = %v, want 20ms", summary.MaxDuration)
	}
	if len(summary.Errors) != 1 {
		t.Fatalf("len(Errors) = %d, want 1", len(summary.Errors))
	}
	if !errors.Is(summary.Errors[0], workerErr) {
		t.Fatalf("Errors[0] = %v, want workerErr", summary.Errors[0])
	}
}

func TestCollectReturnsOnContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results := make(chan Result)
	summary := Collect(ctx, results)

	if len(summary.Errors) != 1 {
		t.Fatalf("len(Errors) = %d, want 1", len(summary.Errors))
	}
	if !errors.Is(summary.Errors[0], context.Canceled) {
		t.Fatalf("Errors[0] = %v, want context.Canceled", summary.Errors[0])
	}
}
