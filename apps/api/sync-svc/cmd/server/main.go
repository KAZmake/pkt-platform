package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/cache"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/config"
	dbutil "github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/db"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/handler"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/onec"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/repository"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/robfig/cron/v3"
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
	if err := dbutil.RunMigrations(cfg.DatabaseURL, cfg.MigrationsDir, "schema_migrations_sync"); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	// ── Redis / Valkey ────────────────────────────────────────────────────────
	redisClient, err := cache.New(cfg.ValkeyURL)
	if err != nil {
		slog.Error("failed to init Redis client", "error", err)
		os.Exit(1)
	}
	pingCtx2, pingCancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel2()
	if err := redisClient.Ping(pingCtx2); err != nil {
		slog.Warn("Redis unavailable — cache degraded, DB fallback active", "error", err)
	} else {
		slog.Info("Redis connected", "url", cfg.ValkeyURL)
	}

	// ── NATS JetStream (optional) ─────────────────────────────────────────────
	var js jetstream.JetStream
	nc, err := natsgo.Connect(cfg.NatsURL)
	if err != nil {
		slog.Warn("NATS unavailable — collateral alerts disabled", "url", cfg.NatsURL, "error", err)
	} else {
		defer nc.Drain()
		if js, err = jetstream.New(nc); err != nil {
			slog.Warn("JetStream init failed", "error", err)
			js = nil
		} else {
			slog.Info("NATS JetStream connected", "url", cfg.NatsURL)
		}
	}

	// ── 1С client ─────────────────────────────────────────────────────────────
	var oneCClient onec.Client
	if cfg.OneCBaseURL != "" {
		oneCClient = onec.NewHTTPClient(cfg.OneCBaseURL, cfg.OneCUser, cfg.OneCPassword)
		slog.Info("1С HTTP client ready", "url", cfg.OneCBaseURL)
	} else {
		oneCClient = onec.NewMockClient()
	}

	// ── Dependencies ──────────────────────────────────────────────────────────
	loanRepo := repository.NewLoanRepository(db)
	colRepo := repository.NewCollateralRepository(db)

	syncSvc := service.NewSyncService(oneCClient, loanRepo, redisClient)
	monitorSvc := service.NewMonitorService(colRepo, js)
	loanSvc := service.NewLoanService(loanRepo, redisClient)

	healthHandler := handler.NewHealthHandler(db)
	loanHandler := handler.NewLoanHandler(loanSvc)

	// ── Cron scheduler ────────────────────────────────────────────────────────
	cr := cron.New()
	if _, err := cr.AddFunc(cfg.CronSchedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		syncSvc.Run(ctx)
		monitorSvc.Run(ctx)
	}); err != nil {
		slog.Error("failed to register cron job", "schedule", cfg.CronSchedule, "error", err)
		os.Exit(1)
	}
	cr.Start()
	slog.Info("cron scheduler started", "schedule", cfg.CronSchedule)

	// Run sync once immediately on startup (warm cache before first request)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		syncSvc.Run(ctx)
		monitorSvc.Run(ctx)
	}()

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:8081"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler.Check)

		// Loan data — internal endpoints (served to core-api for ЛК, expertise-svc for monitoring)
		r.Get("/loans", loanHandler.ListLoans) // ?status=active|overdue|closed
		r.Get("/loans/{id}", loanHandler.GetLoan)
		r.Get("/loans/{id}/schedule", loanHandler.GetSchedule)
		r.Get("/loans/{id}/debts", loanHandler.GetDebts)
		r.Get("/borrowers/{borrower_id}/loans", loanHandler.GetBorrowerLoans)
	})

	// ── Server ────────────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("sync-svc starting", "port", cfg.Port, "env", cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down sync-svc...")
	cr.Stop()
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutCancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}
	slog.Info("sync-svc stopped")
}

func maskDSN(dsn string) string {
	if i := strings.Index(dsn, "@"); i != -1 {
		if j := strings.LastIndex(dsn[:i], ":"); j != -1 {
			return dsn[:j+1] + "***" + dsn[i:]
		}
	}
	return dsn
}
