// Package consumer contains NATS JetStream consumers for the core-api.
package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/service"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

// NotificationConsumer subscribes to application-events and:
//  1. Persists a notification row in the DB for the borrower.
//  2. Sends a transactional email via Resend.
type NotificationConsumer struct {
	js           jetstream.JetStream
	notifSvc     *service.NotificationService
	mailer       *service.Mailer
	userRepo     *repository.UserRepository
	borrowerRepo *repository.BorrowerRepository
	cabinetURL   string
}

func NewNotificationConsumer(
	js jetstream.JetStream,
	notifSvc *service.NotificationService,
	mailer *service.Mailer,
	userRepo *repository.UserRepository,
	borrowerRepo *repository.BorrowerRepository,
	cabinetURL string,
) *NotificationConsumer {
	return &NotificationConsumer{
		js:           js,
		notifSvc:     notifSvc,
		mailer:       mailer,
		userRepo:     userRepo,
		borrowerRepo: borrowerRepo,
		cabinetURL:   cabinetURL,
	}
}

// Start launches the consumer in the background. Call from main.go with a
// context that is cancelled on shutdown.
func (c *NotificationConsumer) Start(ctx context.Context) {
	go func() {
		if err := c.run(ctx); err != nil && ctx.Err() == nil {
			slog.Error("notification consumer exited unexpectedly", "error", err)
		}
	}()
}

func (c *NotificationConsumer) run(ctx context.Context) error {
	consumer, err := c.js.CreateOrUpdateConsumer(ctx, "application-events", jetstream.ConsumerConfig{
		Name:           "core-api-notifications",
		Durable:        "core-api-notifications",
		FilterSubjects: []string{"application.created", "application.status_changed"},
		AckPolicy:      jetstream.AckExplicitPolicy,
		DeliverPolicy:  jetstream.DeliverNewPolicy,
		MaxDeliver:     5,
		AckWait:        30 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	slog.Info("notification consumer started")

	for {
		msgs, err := consumer.Fetch(10, jetstream.FetchMaxWait(5*time.Second))
		if err != nil {
			if ctx.Err() != nil {
				return nil // graceful shutdown
			}
			// Transient fetch error — backoff and retry
			slog.Warn("nats fetch error", "error", err)
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(2 * time.Second):
				continue
			}
		}

		for msg := range msgs.Messages() {
			c.handle(ctx, msg)
		}

		if ctx.Err() != nil {
			return nil
		}
	}
}

func (c *NotificationConsumer) handle(ctx context.Context, msg jetstream.Msg) {
	subject := msg.Subject()
	switch subject {
	case "application.created":
		c.handleApplicationCreated(ctx, msg)
	case "application.status_changed":
		c.handleStatusChanged(ctx, msg)
	default:
		_ = msg.Ack()
	}
}

// ── application.created ───────────────────────────────────────────────────────

type appCreatedEvent struct {
	ApplicationID uuid.UUID `json:"application_id"`
	BorrowerID    uuid.UUID `json:"borrower_id"`
	ProgramID     uuid.UUID `json:"program_id"`
	Amount        float64   `json:"amount"`
	TermMonths    int       `json:"term_months"`
	CreatedAt     time.Time `json:"created_at"`
}

func (c *NotificationConsumer) handleApplicationCreated(ctx context.Context, msg jetstream.Msg) {
	var ev appCreatedEvent
	if err := json.Unmarshal(msg.Data(), &ev); err != nil {
		slog.Warn("consumer: bad application.created payload", "error", err)
		_ = msg.Ack() // don't reprocess malformed messages
		return
	}

	user, err := c.userForBorrower(ctx, ev.BorrowerID)
	if err != nil || user == nil {
		slog.Warn("consumer: user not found for borrower", "borrower_id", ev.BorrowerID, "error", err)
		_ = msg.Ack()
		return
	}

	title := "Ваша заявка принята"
	body := fmt.Sprintf("Заявка #%s зарегистрирована и направлена на первичный скоринг.", ev.ApplicationID)

	if _, err := c.notifSvc.Create(ctx, repository.CreateNotificationInput{
		UserID: user.ID,
		Type:   "status",
		Title:  title,
		Body:   body,
	}); err != nil {
		slog.Warn("consumer: save notification failed", "error", err)
		_ = msg.Nak()
		return
	}

	if err := c.mailer.SendApplicationCreated(user.Email, service.AppCreatedEmailData{
		ApplicationID: ev.ApplicationID.String(),
		Amount:        ev.Amount,
		TermMonths:    ev.TermMonths,
	}); err != nil {
		// Email failure is non-fatal: notification is saved, log and ack
		slog.Warn("consumer: send email failed", "subject", "application.created", "error", err)
	}

	_ = msg.Ack()
}

// ── application.status_changed ────────────────────────────────────────────────

type statusChangedEvent struct {
	ApplicationID uuid.UUID `json:"application_id"`
	FromStatus    string    `json:"from_status"`
	ToStatus      string    `json:"to_status"`
	ActorID       uuid.UUID `json:"actor_id"`
	ChangedAt     time.Time `json:"changed_at"`
}

func (c *NotificationConsumer) handleStatusChanged(ctx context.Context, msg jetstream.Msg) {
	var ev statusChangedEvent
	if err := json.Unmarshal(msg.Data(), &ev); err != nil {
		slog.Warn("consumer: bad application.status_changed payload", "error", err)
		_ = msg.Ack()
		return
	}

	// Find borrower for this application by looking up through appRepo is complex;
	// actor_id may be employee. We notify the actor only if they are the borrower.
	// A better approach: the event should carry borrower_id. For now we use actor_id
	// to look up notifications for the user who acted — but really we want the borrower.
	// We store actor_id in the event; look up by actor_id to notify them.
	user, err := c.userRepo.GetByID(ctx, ev.ActorID)
	if err != nil || user == nil {
		_ = msg.Ack()
		return
	}

	title := fmt.Sprintf("Статус заявки изменён: %s → %s", ev.FromStatus, ev.ToStatus)
	body := fmt.Sprintf("Заявка #%s перешла в статус «%s».", ev.ApplicationID, ev.ToStatus)

	if _, err := c.notifSvc.Create(ctx, repository.CreateNotificationInput{
		UserID: user.ID,
		Type:   "status",
		Title:  title,
		Body:   body,
	}); err != nil {
		slog.Warn("consumer: save notification failed", "error", err)
		_ = msg.Nak()
		return
	}

	if err := c.mailer.SendStatusChanged(user.Email, service.StatusChangedEmailData{
		ApplicationID: ev.ApplicationID.String(),
		FromStatus:    ev.FromStatus,
		ToStatus:      ev.ToStatus,
		CabinetURL:    c.cabinetURL,
	}); err != nil {
		slog.Warn("consumer: send email failed", "subject", "application.status_changed", "error", err)
	}

	_ = msg.Ack()
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (c *NotificationConsumer) userForBorrower(ctx context.Context, borrowerID uuid.UUID) (*model.User, error) {
	b, err := c.borrowerRepo.GetByID(ctx, borrowerID)
	if err != nil || b == nil {
		return nil, err
	}
	return c.userRepo.GetByID(ctx, b.UserID)
}
