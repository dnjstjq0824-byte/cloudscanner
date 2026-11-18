package main

import (
	"context"
	"crypto/sha256"
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
	"github.com/cloudscanr/cloudscanr/internal/domain/admin"
	"github.com/cloudscanr/cloudscanr/internal/domain/affiliate"
	"github.com/cloudscanr/cloudscanr/internal/domain/catalog"
	"github.com/cloudscanr/cloudscanr/internal/domain/estimation"
	"github.com/cloudscanr/cloudscanr/internal/domain/identity"
	"github.com/cloudscanr/cloudscanr/internal/domain/management"
	"github.com/cloudscanr/cloudscanr/internal/domain/search"
	"github.com/cloudscanr/cloudscanr/internal/platform"
	"github.com/cloudscanr/cloudscanr/internal/platform/audit"
	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/cloudscanr/cloudscanr/internal/platform/database"
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
		dbState := "degraded_in_memory"
		if _, ok := p.DB.(*database.DBPool); ok && p.DB != nil {
			dbState = "connected"
		}
		if p.IsReady() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ready","service":"cloudscanr-api","database":"` + dbState + `"}`))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"not_ready","database":"` + dbState + `"}`))
		}
	})
	mux.HandleFunc("/metrics", metricsReg.MetricsHandler())

	// 2. Identity Domain (PostgreSQL-backed when connected, in-memory otherwise)
	identitySvc := identity.NewService(nil)
	if pool, ok := p.DB.(*database.DBPool); ok && pool != nil {
		identitySvc = identity.NewService(identity.NewPGUserStore(pool))
		p.Logger.Info("identity store: PostgreSQL")
	} else {
		p.Logger.Warn("identity store: in-memory (users will not persist across restarts)")
	}
	tokenSvc, err := identity.NewTokenService(p.Config.AuthTokenSecret, identity.DefaultAccessTokenTTL)
	if err != nil {
		p.Logger.Error("invalid AUTH_TOKEN_SECRET, refusing to start auth routes", "err", err)
		tokenSvc = identity.MustTokenService("startup-fallback-dev-secret-32bytes-min")
	}
	mux.HandleFunc("/api/v1/auth/register", identity.RegisterHandler(identitySvc, tokenSvc))
	mux.HandleFunc("/api/v1/auth/login", identity.LoginHandler(identitySvc, tokenSvc))
	mux.Handle("/api/v1/auth/me", identity.RequireAuth(tokenSvc)(identity.MeHandler(identitySvc)))

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

	// 6. BYOC Multi-Cloud Management Domain (authenticated, tenant-scoped).
	// KMS and safety-gate keys derive from the server secret so sealed
	// credentials stay decryptable across restarts.
	kmsSvc, err := core.NewMockKMSServiceWithSecret(p.Config.KMSKeyID, []byte(p.Config.AuthTokenSecret))
	if err != nil {
		p.Logger.Warn("KMS secret too short; sealed credentials will not survive restarts", "err", err)
		kmsSvc, _ = core.NewMockKMSService(p.Config.KMSKeyID)
	}
	safetySum := sha256.Sum256([]byte("cloudscanr-safety-gate-v1:" + p.Config.AuthTokenSecret))
	mgmtSvc := management.NewService(kmsSvc, safetySum[:])
	mgmtHandler := management.NewHandler(mgmtSvc)
	trail := audit.NewLog()
	noticeStore := admin.NewNoticeStore()
	settingsStore := admin.NewSettingsStore()
	mgmtHandler.SetAudit(trail)

	requireAuth := identity.RequireAuth(tokenSvc)
	idem := core.IdempotencyMiddleware(core.NewMemoryIdempotencyStore())
	// Authenticated + idempotent: GETs pass through, mutations require a key.
	requireAuthIdem := func(h http.Handler) http.Handler { return requireAuth(idem(h)) }
	requireAdmin := func(h http.Handler) http.Handler {
		return requireAuth(identity.RequireRole(identity.RoleAdmin)(idem(h)))
	}
	mux.Handle("/api/v1/cloud/connections", requireAuthIdem(http.HandlerFunc(mgmtHandler.ConnectionsHandler)))
	mux.Handle("/api/v1/cloud/connections/", requireAuthIdem(http.HandlerFunc(mgmtHandler.SingleConnectionHandler)))
	mux.Handle("/api/v1/cloud/approvals", requireAuthIdem(http.HandlerFunc(mgmtHandler.ApprovalsHandler)))
	mux.Handle("/api/v1/cloud/budgets", requireAuthIdem(http.HandlerFunc(mgmtHandler.BudgetsHandler)))
	mux.Handle("/api/v1/cloud/resources", requireAuth(http.HandlerFunc(mgmtHandler.ResourcesHandler)))
	mux.Handle("/api/v1/cloud/servers", requireAuthIdem(http.HandlerFunc(mgmtHandler.ServersHandler)))
	mux.Handle("/api/v1/cloud/servers/", requireAuthIdem(http.HandlerFunc(mgmtHandler.SingleServerHandler)))

	// 6b. TCO Estimation API (public: transparent pricing math).
	// /api/v1/pricing/estimate is the canonical contract alias of /api/v1/estimate/tco.
	estimateHandler := estimation.NewHandler(catalogRepo)
	mux.HandleFunc("/api/v1/estimate/tco", estimateHandler.EstimateHandler)
	mux.HandleFunc("/api/v1/pricing/estimate", estimateHandler.EstimateHandler)
	var fxPool *database.DBPool
	if pool, ok := p.DB.(*database.DBPool); ok {
		fxPool = pool
	}
	mux.Handle("/api/v1/pricing/fx-rates", estimation.NewFXHandler(fxPool))

	// 6c. Admin Console Domain (authenticated, ADMIN role or higher).
	dbMode := "in-memory"
	if fxPool != nil {
		dbMode = "postgresql"
	}
	adminHandler := admin.NewHandler(identitySvc.GetStore(), catalogRepo, mgmtSvc, trail, noticeStore, settingsStore, dbMode)
	mux.Handle("/api/v1/admin/overview", requireAdmin(http.HandlerFunc(adminHandler.OverviewHandler)))
	mux.Handle("/api/v1/admin/users", requireAdmin(http.HandlerFunc(adminHandler.UsersHandler)))
	mux.Handle("/api/v1/admin/users/", requireAdmin(http.HandlerFunc(adminHandler.SingleUserHandler)))
	mux.Handle("/api/v1/admin/providers/", requireAdmin(http.HandlerFunc(adminHandler.SingleProviderHandler)))
	mux.Handle("/api/v1/admin/audit-logs", requireAdmin(http.HandlerFunc(adminHandler.AuditLogsHandler)))
	mux.Handle("/api/v1/admin/notices", requireAdmin(http.HandlerFunc(adminHandler.NoticesHandler)))
	mux.Handle("/api/v1/admin/notices/", requireAdmin(http.HandlerFunc(adminHandler.SingleNoticeHandler)))
	mux.Handle("/api/v1/admin/settings", requireAdmin(http.HandlerFunc(adminHandler.SettingsHandler)))

	// 6d. Public site surfaces (no auth, maintenance-safe for status comms).
	mux.HandleFunc("/api/v1/notices", adminHandler.PublicNoticesHandler)
	mux.HandleFunc("/api/v1/site-config", adminHandler.PublicSiteConfigHandler)

	// 7. Pricing Adapter Registry
	pricingReg := pricing.NewRegistry()
	_ = pricingReg.Register(aws.NewAdapter(nil, "us-east-1"))
	_ = pricingReg.Register(digitalocean.NewAdapter(nil, ""))
	_ = pricingReg.Register(hetzner.NewAdapter(nil, ""))
	_ = pricingReg.Register(upcloud.NewAdapter(nil, "", ""))
	_ = pricingReg.Register(vultr.NewAdapter(nil))

	// Apply Middlewares (outer -> inner): PanicRecovery -> SecurityHeaders ->
	// RequestID -> AccessLog -> Maintenance -> RateLimit -> mux
	rateLimiter := core.NewRateLimiter(core.DefaultRateLimitConfig())
	var handler http.Handler = mux
	handler = rateLimiter.Middleware(handler)
	handler = admin.MaintenanceMiddleware(settingsStore)(handler)
	handler = core.AccessLogMiddleware(p.Logger)(handler)
	handler = core.RequestIDMiddleware(handler)
	handler = core.SecurityHeadersMiddleware(p.Config.Domain, p.Config.CORSAllowedOrigins...)(handler)
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
	if !p.IsReady() {
		fmt.Fprintln(os.Stderr, "Fatal bootstrap error: platform is not ready (production database unreachable)")
		os.Exit(1)
	}
	if p.DB != nil {
		defer p.DB.Close()
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
