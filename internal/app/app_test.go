package app

import (
	"Pipepool/internal/testutil"
	"bytes"
	"context"
	"errors"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunHappyPathReturnsExpectedSummary(t *testing.T) {
	cfg := LoadConfig(2, 1, 2*time.Second, 500*time.Millisecond)
	inputs := testutil.BuildInputs(
		"  hello world  ",
		"",
		"line one\r\nline two",
	)

	summary, err := Run(context.Background(), cfg, inputs, testutil.NewDiscardLogger())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summary.TotalJobs != 3 {
		t.Fatalf("TotalJobs = %d, want 3", summary.TotalJobs)
	}
	if summary.SuccessCount != 2 {
		t.Fatalf("SuccessCount = %d, want 2", summary.SuccessCount)
	}
	if summary.FailureCount != 1 {
		t.Fatalf("FailureCount = %d, want 1", summary.FailureCount)
	}
	if summary.TotalWords != 6 {
		t.Fatalf("TotalWords = %d, want 6", summary.TotalWords)
	}
	if summary.TotalLines != 3 {
		t.Fatalf("TotalLines = %d, want 3", summary.TotalLines)
	}
}

func TestRunCancellationStopsPipelineAndWorkers(t *testing.T) {
	cfg := LoadConfig(2, 1, 40*time.Millisecond, 500*time.Millisecond)
	inputs := testutil.BuildManyInputs(60, "alpha beta gamma")

	start := time.Now()
	summary, err := Run(context.Background(), cfg, inputs, testutil.NewDiscardLogger())
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected cancellation error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
	if elapsed > time.Second {
		t.Fatalf("canceled run returned too slowly: %v", elapsed)
	}
	if summary.TotalJobs >= len(inputs) {
		t.Fatalf("expected partial processing before cancel, got %d of %d", summary.TotalJobs, len(inputs))
	}
}

func TestRunGoroutinesReturnNearBaselineAfterShutdown(t *testing.T) {
	cfg := LoadConfig(3, 2, 2*time.Second, 300*time.Millisecond)
	inputs := testutil.BuildManyInputs(80, "alpha beta")

	baseline := runtime.NumGoroutine()

	_, err := Run(context.Background(), cfg, inputs, testutil.NewDiscardLogger())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	maxAllowed := baseline + 12
	deadline := time.Now().Add(1500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= maxAllowed {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("goroutines did not return near baseline: baseline=%d current=%d allowed=%d", baseline, runtime.NumGoroutine(), maxAllowed)
}

func TestRunLogsLifecycleStates(t *testing.T) {
	cfg := LoadConfig(2, 1, time.Second, 500*time.Millisecond)
	inputs := testutil.BuildInputs("hello", "world")

	var buf bytes.Buffer
	summary, err := Run(context.Background(), cfg, inputs, testutil.NewBufferLogger(&buf))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if summary.TotalJobs != 2 {
		t.Fatalf("TotalJobs = %d, want 2", summary.TotalJobs)
	}

	logs := buf.String()
	required := []string{
		"component=app",
		"state=start",
		"state=stop",
		"component=pipeline",
		"stage=enqueue",
		"state=waiting",
		"component=pool",
		"state=worker_started",
		"state=worker_finished",
	}

	for _, token := range required {
		if !strings.Contains(logs, token) {
			t.Fatalf("expected logs to contain %q, got:\n%s", token, logs)
		}
	}
}
