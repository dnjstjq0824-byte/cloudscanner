package main

import (
	"context"
	"fmt"
	"net/http"
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
	"github.com/cloudscanr/cloudscanr/internal/domain/affiliate"
	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/cloudscanr/cloudscanr/internal/domain/identity"
	"github.com/cloudscanr/cloudscanr/internal/domain/management"
	"github.com/cloudscanr/cloudscanr/internal/domain/search"
	"github.com/cloudscanr/cloudscanr/internal/platform"
	"github.com/cloudscanr/cloudscanr/internal/platform/core"
)

// SetupRouter registers all domain routes and middlewares.
func SetupRouter(p *platform.Platform) http.Handler {
	mux := http.NewServeMux()

	// 1. Core & Telemetry
	metricsReg := core.NewMetricsRegistry()
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy","service":"cloudscanr-api","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
	})
	mux.HandleFunc("/api/v1/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if p.IsReady() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ready","service":"cloudscanr-api"}`))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"not_ready"}`))
		}
	})
	mux.HandleFunc("/metrics", metricsReg.MetricsHandler())

	// 2. Identity Domain
	identitySvc := identity.NewService(nil)
	mux.HandleFunc("/api/v1/auth/register", identity.RegisterHandler(identitySvc))
	mux.HandleFunc("/api/v1/auth/login", identity.LoginHandler(identitySvc))

	// 3. Catalog Domain
	catalogRepo := catalog.NewMemoryProviderRepo()
	catalogHandler := catalog.NewHandler(catalogRepo)
	mux.HandleFunc("/api/v1/providers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/providers" || r.URL.Path == "/api/v1/providers/" {
			catalogHandler.ListProvidersHandler(w, r)
		} else {
			catalogHandler.GetProviderBySlugHandler(w, r)
		}
	})
	mux.HandleFunc("/api/v1/providers/", catalogHandler.GetProviderBySlugHandler)
	mux.HandleFunc("/api/v1/catalog/categories", catalogHandler.ListCategoriesHandler)
	mux.HandleFunc("/api/v1/catalog/skus", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/catalog/skus" || r.URL.Path == "/api/v1/catalog/skus/" {
			catalogHandler.ListSKUsHandler(w, r)
		} else {
			catalogHandler.GetSKUHandler(w, r)
		}
	})
	mux.HandleFunc("/api/v1/catalog/skus/", catalogHandler.GetSKUHandler)

	// 4. Search & Ranking Domain
	searchSvc := search.NewService(catalogRepo)
	searchHandler := search.NewHandler(searchSvc)
	mux.HandleFunc("/api/v1/search/skus", searchHandler.SearchSKUsHandler)
	mux.HandleFunc("/api/v1/rankings/best-value", searchHandler.BestValueRankingsHandler)

	// 5. Affiliate Domain
	affiliateSvc := affiliate.NewService(nil, nil)
	mux.HandleFunc("/go/", affiliate.RedirectHandler(affiliateSvc))
	mux.HandleFunc("/api/v1/affiliate/disclosure", affiliate.DisclosureHandler(affiliateSvc))
	mux.HandleFunc("/api/v1/webhooks/affiliate/", affiliate.WebhookHandler(affiliateSvc))

	// 6. BYOC Multi-Cloud Management Domain
	mgmtSvc := management.NewService(nil, nil)
	mgmtHandler := management.NewHandler(mgmtSvc)
	mux.HandleFunc("/api/v1/cloud/connections", mgmtHandler.ConnectionsHandler)
	mux.HandleFunc("/api/v1/cloud/connections/", mgmtHandler.SingleConnectionHandler)
	mux.HandleFunc("/api/v1/cloud/approvals", mgmtHandler.ApprovalsHandler)

	// 7. Pricing Adapter Registry
	pricingReg := pricing.NewRegistry()
	_ = pricingReg.Register(aws.NewAdapter(nil, "us-east-1"))
	_ = pricingReg.Register(digitalocean.NewAdapter(nil, ""))
	_ = pricingReg.Register(hetzner.NewAdapter(nil, ""))
	_ = pricingReg.Register(upcloud.NewAdapter(nil, "", ""))
	_ = pricingReg.Register(vultr.NewAdapter(nil))

	// Apply Middlewares: SecurityHeaders -> PanicRecovery
	var handler http.Handler = mux
	handler = core.SecurityHeadersMiddleware(p.Config.Domain)(handler)
	handler = core.PanicRecoveryMiddleware(p.Logger)(handler)

	return handler
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := platform.DefaultConfig()
	p, err := platform.Bootstrap(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal bootstrap error: %v\n", err)
		os.Exit(1)
	}

	handler := SetupRouter(p)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:      handler,
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
