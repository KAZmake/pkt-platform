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
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/consumer"
	dbutil "github.com/KAZmake/pkt-platform/apps/api/core-api/internal/db"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/handler"
	apimw "github.com/KAZmake/pkt-platform/apps/api/core-api/internal/middleware"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/service"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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

	// ── MinIO ─────────────────────────────────────────────────────────────────
	mc, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		slog.Error("failed to init MinIO client", "error", err)
		os.Exit(1)
	}
	slog.Info("MinIO client ready", "endpoint", cfg.MinioEndpoint)

	// ── NATS JetStream (optional — non-fatal if unavailable) ──────────────────
	var js jetstream.JetStream
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		slog.Warn("NATS unavailable — events disabled", "url", cfg.NatsURL, "error", err)
	} else {
		defer nc.Drain()
		js, err = jetstream.New(nc)
		if err != nil {
			slog.Warn("JetStream init failed — events disabled", "error", err)
			js = nil
		} else {
			slog.Info("NATS JetStream connected", "url", cfg.NatsURL)
		}
	}

	// ── Dependencies ─────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	borrowerRepo := repository.NewBorrowerRepository(db)
	programRepo := repository.NewProgramRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	docRepo := repository.NewDocumentRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	userSvc := service.NewUserService(userRepo, borrowerRepo)
	programSvc := service.NewProgramService(programRepo)
	appSvc := service.NewApplicationService(appRepo, js)
	docSvc := service.NewDocumentService(docRepo, mc)
	scheduleSvc := service.NewScheduleService(appRepo, programRepo)
	ticketSvc := service.NewTicketService(ticketRepo)
	notifSvc := service.NewNotificationService(notifRepo)
	mailer := service.NewMailer(cfg.ResendAPIKey, cfg.EmailFrom)

	healthHandler := handler.NewHealthHandler(db)
	userHandler := handler.NewUserHandler(userSvc)
	programHandler := handler.NewProgramHandler(programSvc)
	appHandler := handler.NewApplicationHandler(appSvc)
	docHandler := handler.NewDocumentHandler(docSvc, scheduleSvc)
	ticketHandler := handler.NewTicketHandler(ticketSvc)
	notifHandler := handler.NewNotificationHandler(notifSvc)

	// ── NATS consumer ─────────────────────────────────────────────────────────
	if js != nil {
		notifConsumer := consumer.NewNotificationConsumer(
			js, notifSvc, mailer, userRepo, borrowerRepo, cfg.CabinetURL,
		)
		notifConsumer.Start(context.Background())
	}

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(apimw.CORS([]string{"http://localhost:3000", "http://localhost:8055"}))

	r.Route("/api/v1", func(r chi.Router) {
		// ── Public ────────────────────────────────────────────────────────────
		r.Get("/health", healthHandler.Check)
		r.Get("/programs", programHandler.ListPrograms)
		r.Get("/programs/{id}", programHandler.GetProgram)

		// ── Authenticated ─────────────────────────────────────────────────────
		r.Group(func(r chi.Router) {
			r.Use(apimw.Authenticate(jwks))
			r.Use(handler.SyncUser(userSvc))

			// Current user profile
			r.Get("/me", userHandler.GetMe)
			r.Put("/me", userHandler.UpdateMe)

			// Borrower profile (borrower role only)
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("borrower", "admin"))
				r.Get("/me/borrower", userHandler.GetMyBorrower)
				r.Put("/me/borrower", userHandler.UpdateMyBorrower)
			})

			// Employee+ endpoints
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("employee", "expert", "admin"))
				r.Get("/users", userHandler.ListUsers)
				r.Get("/users/{id}", userHandler.GetUser)
			})

			// Admin-only program management
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("admin"))
				r.Post("/programs", programHandler.CreateProgram)
				r.Put("/programs/{id}", programHandler.UpdateProgram)
				r.Delete("/programs/{id}", programHandler.DeactivateProgram)
			})

			// Applications
			r.Get("/applications", appHandler.ListApplications)
			r.Get("/applications/{id}", appHandler.GetApplication)
			r.Get("/applications/{id}/schedule", docHandler.GetSchedule)
			r.Get("/applications/{id}/documents", docHandler.ListApplicationDocuments)
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("borrower"))
				r.Post("/applications", appHandler.CreateApplication)
				r.Post("/applications/{id}/documents/upload-url", docHandler.InitiateUpload)
			})
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("employee", "expert", "admin"))
				r.Patch("/applications/{id}/status", appHandler.ChangeStatus)
			})

			// Documents
			r.Get("/documents/{id}/download-url", docHandler.GetDownloadURL)
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("borrower", "admin"))
				r.Delete("/documents/{id}", docHandler.DeleteDocument)
			})

			// Tickets (обращения)
			r.Get("/tickets", ticketHandler.ListTickets)
			r.Get("/tickets/{id}", ticketHandler.GetTicket)
			r.Post("/tickets/{id}/messages", ticketHandler.AddMessage)
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("borrower"))
				r.Post("/tickets", ticketHandler.CreateTicket)
			})
			r.Group(func(r chi.Router) {
				r.Use(apimw.RequireRole("employee", "expert", "admin"))
				r.Patch("/tickets/{id}/status", ticketHandler.ChangeTicketStatus)
			})

			// Notifications
			r.Get("/notifications", notifHandler.ListNotifications)
			r.Get("/notifications/unread-count", notifHandler.UnreadCount)
			r.Patch("/notifications/read-all", notifHandler.MarkAllRead)
			r.Patch("/notifications/{id}/read", notifHandler.MarkRead)
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
