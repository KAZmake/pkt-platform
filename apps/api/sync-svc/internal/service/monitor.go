package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/repository"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	insuranceWarningDays  = 30  // publish warning if expiry < 30 days
	insuranceCriticalDays = 7   // publish critical if expiry < 7 days
	inventoryOverdueDays  = 365 // publish alert if no inventory for > 1 year
)

// MonitorService checks collateral health and publishes NATS alerts (2.3.3).
type MonitorService struct {
	colRepo *repository.CollateralRepository
	js      jetstream.JetStream
}

func NewMonitorService(colRepo *repository.CollateralRepository, js jetstream.JetStream) *MonitorService {
	return &MonitorService{colRepo: colRepo, js: js}
}

// Run checks all active collaterals and publishes alerts for violations.
// Called by the cron scheduler alongside the sync cycle.
func (s *MonitorService) Run(ctx context.Context) {
	cols, err := s.colRepo.ListActive(ctx)
	if err != nil {
		slog.Error("monitor: failed to list collaterals", "error", err)
		return
	}

	now := time.Now()
	alerts := 0

	for _, c := range cols {
		// Insurance expiry check
		if c.InsuranceExpiry != nil {
			daysUntil := int(c.InsuranceExpiry.Sub(now).Hours() / 24)
			if daysUntil <= insuranceCriticalDays {
				s.publish(ctx, "sync.collateral.insurance_expiring", model.CollateralAlert{
					CollateralID: c.ID,
					Type:         "insurance_expiring",
					Severity:     "critical",
					Message:      fmt.Sprintf("Страховка истекает через %d дн. (%s)", daysUntil, c.InsuranceExpiry.Format("02.01.2006")),
					DaysUntil:    daysUntil,
				})
				alerts++
			} else if daysUntil <= insuranceWarningDays {
				s.publish(ctx, "sync.collateral.insurance_expiring", model.CollateralAlert{
					CollateralID: c.ID,
					Type:         "insurance_expiring",
					Severity:     "warning",
					Message:      fmt.Sprintf("Страховка истекает через %d дн. (%s)", daysUntil, c.InsuranceExpiry.Format("02.01.2006")),
					DaysUntil:    daysUntil,
				})
				alerts++
			}
		}

		// Inventory overdue check
		if c.LastInventoryDate != nil {
			daysSince := int(now.Sub(*c.LastInventoryDate).Hours() / 24)
			if daysSince >= inventoryOverdueDays {
				s.publish(ctx, "sync.collateral.inventory_overdue", model.CollateralAlert{
					CollateralID: c.ID,
					Type:         "inventory_overdue",
					Severity:     "warning",
					Message:      fmt.Sprintf("Последняя инвентаризация %d дн. назад (%s)", daysSince, c.LastInventoryDate.Format("02.01.2006")),
					DaysSince:    daysSince,
				})
				alerts++
			}
		} else {
			// No inventory on record at all
			s.publish(ctx, "sync.collateral.inventory_overdue", model.CollateralAlert{
				CollateralID: c.ID,
				Type:         "inventory_overdue",
				Severity:     "warning",
				Message:      "Инвентаризация залога не проводилась",
			})
			alerts++
		}
	}

	if alerts > 0 {
		slog.Warn("monitor: collateral alerts published", "count", alerts)
	} else {
		slog.Info("monitor: all collaterals OK", "checked", len(cols))
	}
}

func (s *MonitorService) publish(ctx context.Context, subject string, alert model.CollateralAlert) {
	if s.js == nil {
		return
	}
	data, err := json.Marshal(alert)
	if err != nil {
		slog.Warn("monitor: marshal error", "error", err)
		return
	}
	pubCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if _, err := s.js.Publish(pubCtx, subject, data); err != nil {
		slog.Warn("monitor: publish error", "subject", subject, "error", err)
	}
}
