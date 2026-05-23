package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/internal/config"
	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/internal/directus"
	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/internal/handler"
	appMiddleware "github.com/KAZmake/pkt-platform/apps/api/assistant-svc/internal/middleware"
	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/internal/prompt"
	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// ── Anthropic client ──────────────────────────────────────────────────────
	if cfg.AnthropicAPIKey == "" {
		slog.Warn("ANTHROPIC_API_KEY not set — /chat will return 503")
	}
	anthropicClient := anthropic.NewClient(option.WithAPIKey(cfg.AnthropicAPIKey))

	// ── Directus client + system prompt builder ───────────────────────────────
	directusClient := directus.NewClient(cfg.DirectusURL, cfg.DirectusToken)
	promptBuilder := prompt.NewBuilder(directusClient)

	// Warm up the system prompt in the background.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = promptBuilder.Get(ctx)
	}()

	// ── Handlers ──────────────────────────────────────────────────────────────
	chatHandler := handler.NewChatHandler(anthropicClient, promptBuilder, cfg.AnthropicModel, cfg.MaxTokens)

	// ── Rate limiter ──────────────────────────────────────────────────────────
	var rateLimitMiddleware func(http.Handler) http.Handler
	if cfg.RateLimitRPM > 0 {
		rl := appMiddleware.NewRateLimiter(cfg.RateLimitRPM)
		rateLimitMiddleware = rl.Middleware
		slog.Info("rate limiter enabled", "rpm_per_ip", cfg.RateLimitRPM)
	} else {
		rateLimitMiddleware = func(next http.Handler) http.Handler { return next }
	}

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handler.Health)

		// Chat endpoint — rate-limited per IP
		r.With(rateLimitMiddleware).Post("/chat", chatHandler.Chat)
	})

	// ── Server ────────────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second, // Anthropic calls can take ~10–15s
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("assistant-svc starting", "port", cfg.Port, "env", cfg.Environment, "model", cfg.AnthropicModel)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down assistant-svc...")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutCancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}
	slog.Info("assistant-svc stopped")
}
