package types

import (
	"time"
)

type Config struct {
	WorkerCount   int
	QueueSize     int
	RunTimeout    time.Duration
	PerJobTimeout time.Duration
}

type Job struct {
	ID    int
	Input string
}

type Result struct {
	ID        int
	Output    string
	Valid     bool
	WordCount int
	LineCount int
	Duration  time.Duration
	Err       error
}

type Summary struct {
	TotalJobs     int
	SuccessCount  int
	FailureCount  int
	TotalWords    int
	TotalLines    int
	TotalDuration time.Duration
	SlowestJobID  int
	MaxDuration   time.Duration
	Errors        []error
}
