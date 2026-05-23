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

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/config"
	dbutil "github.com/KAZmake/pkt-platform/apps/api/core-api/internal/db"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/handler"
	apimw "github.com/KAZmake/pkt-platform/apps/api/core-api/internal/middleware"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/pkg/response"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// ── Database ─────────────────────────────────────────────────────────────
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

	// ── Migrations ───────────────────────────────────────────────────────────
	if err := dbutil.RunMigrations(cfg.DatabaseURL, cfg.MigrationsDir, "schema_migrations_core"); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	// ── JWKS (Keycloak) ───────────────────────────────────────────────────────
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs",
		cfg.KeycloakURL, cfg.KeycloakRealm)

	jwks, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		slog.Error("failed to fetch JWKS from Keycloak", "url", jwksURL, "error", err)
		os.Exit(1)
	}
	slog.Info("JWKS loaded", "url", jwksURL)

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(apimw.CORS([]string{"http://localhost:3000", "http://localhost:8055"}))

	healthHandler := handler.NewHealthHandler(db)

	r.Route("/api/v1", func(r chi.Router) {
		// Public
		r.Get("/health", healthHandler.Check)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(apimw.Authenticate(jwks))

			r.Get("/me", func(w http.ResponseWriter, req *http.Request) {
				claims, _ := apimw.ClaimsFromCtx(req.Context())
				response.OK(w, map[string]any{
					"email":    claims.Email,
					"username": claims.PreferredUsername,
					"roles":    claims.RealmAccess.Roles,
				})
			})

			// Employee+ only example
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("employee", "expert", "admin"))
				r.Get("/admin/ping", func(w http.ResponseWriter, _ *http.Request) {
					response.OK(w, map[string]string{"message": "access granted"})
				})
			})
		})
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
		slog.Info("core-api starting", "port", cfg.Port, "env", cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutCancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}
	slog.Info("server stopped")
}

func maskDSN(dsn string) string {
	if i := strings.Index(dsn, "@"); i != -1 {
		if j := strings.LastIndex(dsn[:i], ":"); j != -1 {
			return dsn[:j+1] + "***" + dsn[i:]
		}
	}
	return dsn
}
