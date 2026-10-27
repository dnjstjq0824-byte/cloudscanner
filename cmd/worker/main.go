package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/platform"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := platform.DefaultConfig()
	p, err := platform.Bootstrap(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal worker bootstrap error: %v\n", err)
		os.Exit(1)
	}

	p.Logger.Info("cloudscanr background worker started", "env", cfg.Environment)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.Logger.Info("worker stopping gracefully...")
			return
		case <-ticker.C:
			p.Logger.Debug("worker heartbeat check")
		}
	}
}
