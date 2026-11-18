package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/aws"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/digitalocean"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/hetzner"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/upcloud"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/vultr"
	"github.com/cloudscanr/cloudscanr/internal/platform"
)

// harvestInterval sets the periodic price-catalog collection cadence.
const harvestInterval = 15 * time.Minute

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := platform.DefaultConfig()
	p, err := platform.Bootstrap(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal worker bootstrap error: %v\n", err)
		os.Exit(1)
	}

	reg := pricing.NewRegistry()
	_ = reg.Register(aws.NewAdapter(nil, "us-east-1"))
	_ = reg.Register(digitalocean.NewAdapter(nil, ""))
	_ = reg.Register(hetzner.NewAdapter(nil, ""))
	_ = reg.Register(upcloud.NewAdapter(nil, "", ""))
	_ = reg.Register(vultr.NewAdapter(nil))

	harvest := func(reason string) {
		hctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		start := time.Now()
		payloads, err := reg.HarvestAll(hctx)
		if err != nil {
			p.Logger.Error("pricing harvest failed", "reason", reason, "err", err)
			return
		}
		p.Logger.Info("pricing harvest completed",
			"reason", reason,
			"providers", len(payloads),
			"elapsed_ms", time.Since(start).Milliseconds())
	}

	p.Logger.Info("cloudscanr background worker started", "env", cfg.Environment,
		"harvest_interval", harvestInterval.String())

	harvest("boot")
	ticker := time.NewTicker(harvestInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.Logger.Info("worker stopping gracefully...")
			return
		case <-ticker.C:
			harvest("interval")
		}
	}
}
