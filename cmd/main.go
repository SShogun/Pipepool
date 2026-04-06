package main

import (
	. "Pipepool/internal/app"
	. "Pipepool/internal/logging"
	"context"
	"time"
)

func main() {
	cfg := LoadConfig(5, 100, 4*time.Second, 16*time.Second) // or however you initialize your Config struct
	logger := LoadLogger()                                   // or however you initialize your Logger
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RunTimeout)
	defer cancel()

}
