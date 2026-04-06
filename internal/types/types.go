package types

import (
	"context"
	"time"
)

type Job struct {
	ID    int
	Input string
}

type Result struct {
	ID     int
	Output string
}

type Summary struct {
	TotalJobs    int
	SuccessCount int
	FailureCount int

	TotalDuration time.Duration

	SlowestJobID int
	MaxDuration  time.Duration

	Errors []error
}

func (s Summary) Collect(ctx context.Context, results chan Result) any {
	panic("unimplemented")
}
