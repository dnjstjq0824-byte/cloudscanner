package main

import (
	"context"
	"fmt"
	"net/http"
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
		fmt.Fprintf(os.Stderr, "Fatal bootstrap error: %v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"cloudscanr-api","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
	})
	mux.HandleFunc("/api/v1/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if p.IsReady() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ready","service":"cloudscanr-api"}`))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not_ready"}`))
		}
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		p.Logger.Info("starting CloudScanr API server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		p.Logger.Info("shutting down CloudScanr API server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			p.Logger.Error("server forced shutdown", "err", err)
		}
	case err := <-serverErr:
		p.Logger.Error("server unexpected error", "err", err)
	}
}
