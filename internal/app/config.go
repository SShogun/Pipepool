package app

import "time"

type Config struct {
	WorkerCount   int
	QueueSize     int
	RunTimeout    time.Duration
	PerJobTimeout time.Duration
}

func LoadConfig(w, q int, rt, pjt time.Duration) *Config {
	return &Config{
		WorkerCount:   w,
		QueueSize:     q,
		RunTimeout:    rt,
		PerJobTimeout: pjt,
	}
}
