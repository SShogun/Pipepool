package app

import (
	"Pipepool/internal/types"
	"time"
)

const (
	defaultWorkerCount   = 2
	defaultQueueSize     = 4
	defaultRunTimeout    = 5 * time.Second
	defaultPerJobTimeout = 250 * time.Millisecond
)

func LoadConfig(w, q int, rt, pjt time.Duration) *types.Config {
	cfg := types.Config{
		WorkerCount:   w,
		QueueSize:     q,
		RunTimeout:    rt,
		PerJobTimeout: pjt,
	}

	normalized := normalizeConfig(cfg)
	return &normalized
}

func normalizeConfig(cfg types.Config) types.Config {
	if cfg.WorkerCount <= 0 {
		cfg.WorkerCount = defaultWorkerCount
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = defaultQueueSize
	}
	if cfg.RunTimeout <= 0 {
		cfg.RunTimeout = defaultRunTimeout
	}
	if cfg.PerJobTimeout <= 0 {
		cfg.PerJobTimeout = defaultPerJobTimeout
	}

	return cfg
}
