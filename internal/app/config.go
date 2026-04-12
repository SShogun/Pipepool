package app

import (
	. "Pipepool/internal/types"
	"time"
)

func LoadConfig(w, q int, rt, pjt time.Duration) *Config {
	return &Config{
		WorkerCount:   w,
		QueueSize:     q,
		RunTimeout:    rt,
		PerJobTimeout: pjt,
	}
}
