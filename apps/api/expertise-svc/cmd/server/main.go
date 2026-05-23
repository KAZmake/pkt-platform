package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/config"
	dbutil "github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/db"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/handler"
	apimw "github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/middleware"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/repository"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/service"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := db.Ping(pingCtx); err != nil {
		slog.Error("database ping failed", "error", err)
		os.Exit(1)
	}
	slog.Info("database connected", "url", maskDSN(cfg.DatabaseURL))

	// ── Migrations ────────────────────────────────────────────────────────────
	if err := dbutil.RunMigrations(cfg.DatabaseURL, cfg.MigrationsDir, "schema_migrations_expertise"); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	// ── JWKS (Keycloak) ───────────────────────────────────────────────────────
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs",
		cfg.KeycloakURL, cfg.KeycloakRealm)
	jwks, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		slog.Error("failed to fetch JWKS", "url", jwksURL, "error", err)
		os.Exit(1)
	}
	slog.Info("JWKS loaded", "url", jwksURL)

	// ── NATS JetStream (optional) ─────────────────────────────────────────────
	var js jetstream.JetStream
	nc, err := natsgo.Connect(cfg.NatsURL)
	if err != nil {
		slog.Warn("NATS unavailable — events disabled", "url", cfg.NatsURL, "error", err)
	} else {
		defer nc.Drain()
		if js, err = jetstream.New(nc); err != nil {
			slog.Warn("JetStream init failed", "error", err)
			js = nil
		} else {
			slog.Info("NATS JetStream connected", "url", cfg.NatsURL)
		}
	}

	// ── Dependencies ──────────────────────────────────────────────────────────
	appRepo := repository.NewApplicationRepository(db)
	colRepo := repository.NewCollateralRepository(db)
	concRepo := repository.NewConclusionRepository(db)

	appSvc := service.NewApplicationService(appRepo, colRepo, concRepo, js)
	colSvc := service.NewCollateralService(colRepo)
	concSvc := service.NewConclusionService(concRepo)

	healthHandler := handler.NewHealthHandler(db)
	appHandler := handler.NewApplicationHandler(appSvc)
	colHandler := handler.NewCollateralHandler(colSvc)
	concHandler := handler.NewConclusionHandler(concSvc)
	pdfHandler := handler.NewPDFHandler(appSvc, concSvc)

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(apimw.CORS([]string{"http://localhost:3000"}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler.Check)

		// All expertise endpoints require authentication
		r.Group(func(r chi.Router) {
			r.Use(apimw.Authenticate(jwks))
			r.Use(handler.StoreActor)

			// ── Applications (expert workstation, 2.2.3) ──────────────────────
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Get("/applications", appHandler.ListApplications)
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Get("/applications/{id}", appHandler.GetApplication)

			// Status change: role checked inside service (2.2.2)
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Patch("/applications/{id}/status", appHandler.ChangeStatus)

			// Assign
			r.With(apimw.RequireRole("employee", "admin")).
				Patch("/applications/{id}/assign", appHandler.Assign)

			// ── Collaterals (2.2.4) ───────────────────────────────────────────
			r.With(apimw.RequireRole("expert", "admin")).
				Post("/collaterals", colHandler.CreateCollateral)
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Get("/collaterals/{id}", colHandler.GetCollateral)
			r.With(apimw.RequireRole("expert", "admin")).
				Put("/collaterals/{id}", colHandler.UpdateCollateral)
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Get("/applications/{id}/collaterals", colHandler.ListApplicationCollaterals)
			r.With(apimw.RequireRole("expert", "admin")).
				Post("/applications/{id}/collaterals", colHandler.AttachCollateral)
			r.With(apimw.RequireRole("expert", "admin")).
				Delete("/applications/{id}/collaterals/{col_id}", colHandler.ReleaseCollateral)

			// ── Expert conclusions (2.2.3) ────────────────────────────────────
			r.With(apimw.RequireRole("expert", "admin")).
				Post("/applications/{id}/conclusions", concHandler.SubmitConclusion)
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Get("/applications/{id}/conclusions", concHandler.ListConclusions)

			// ── Committee votes (2.2.5) ───────────────────────────────────────
			r.With(apimw.RequireRole("expert", "admin")).
				Post("/applications/{id}/votes", concHandler.AddVote)
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Get("/applications/{id}/votes", concHandler.ListVotes)

			// ── PDF generation (2.2.5, 2.2.6) ────────────────────────────────
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Get("/applications/{id}/pdf/protocol", pdfHandler.CommitteeProtocol)
			r.With(apimw.RequireRole("employee", "expert", "admin")).
				Get("/applications/{id}/pdf/agreement", pdfHandler.LoanAgreement)
		})
	})

	// ── Server ────────────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("expertise-svc starting", "port", cfg.Port, "env", cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down expertise-svc...")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutCancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}
	slog.Info("expertise-svc stopped")
}

func maskDSN(dsn string) string {
	if i := strings.Index(dsn, "@"); i != -1 {
		if j := strings.LastIndex(dsn[:i], ":"); j != -1 {
			return dsn[:j+1] + "***" + dsn[i:]
		}
	}
	return dsn
}
