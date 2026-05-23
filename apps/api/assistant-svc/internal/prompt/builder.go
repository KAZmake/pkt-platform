package prompt

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"text/template"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/internal/directus"
)

const systemTmpl = `Ты — виртуальный ассистент ТОО «Первое кредитное товарищество» (ПКТ), микрофинансовой организации в Западно-Казахстанской области, специализирующейся на кредитовании агросектора.

Отвечай только на русском языке. Будь вежливым, профессиональным и лаконичным. Если вопрос выходит за рамки деятельности ПКТ, вежливо объясни, что можешь помочь только по вопросам кредитования.

## Программы кредитования ПКТ
{{- range .Programs}}
### {{.Name}}
{{if .Description}}{{.Description}}{{end}}
- Ставка: {{if gt .RateMax 0.0}}{{.RateMin}}–{{.RateMax}}%{{else}}{{.RateMin}}%{{end}} годовых
- Срок: {{if gt .TermMax 0}}{{.TermMin}}–{{.TermMax}} мес.{{else}}до {{.TermMin}} мес.{{end}}
{{if gt .AmountMax 0.0}}- Максимальная сумма: {{printf "%.0f" .AmountMax}} {{.Currency}}{{end}}
{{if .Requirements}}- Требования: {{.Requirements}}{{end}}
{{- end}}

## Часто задаваемые вопросы
{{- range .FAQ}}
**В:** {{.Question}}
**О:** {{.Answer}}
{{end}}

Если клиент хочет подать заявку или узнать подробности, направляй его на официальный сайт или предложи создать обращение через личный кабинет.`

var tmpl = template.Must(template.New("system").Parse(systemTmpl))

const refreshInterval = 15 * time.Minute

// Builder builds and caches the system prompt from Directus data.
type Builder struct {
	client    *directus.Client
	mu        sync.RWMutex
	cached    string
	refreshed time.Time
}

func NewBuilder(client *directus.Client) *Builder {
	return &Builder{client: client}
}

// Get returns the (possibly cached) system prompt.
func (b *Builder) Get(ctx context.Context) string {
	b.mu.RLock()
	if b.cached != "" && time.Since(b.refreshed) < refreshInterval {
		p := b.cached
		b.mu.RUnlock()
		return p
	}
	b.mu.RUnlock()

	return b.refresh(ctx)
}

func (b *Builder) refresh(ctx context.Context) string {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Double-checked locking: another goroutine may have refreshed already.
	if b.cached != "" && time.Since(b.refreshed) < refreshInterval {
		return b.cached
	}

	programs, err := b.client.GetPrograms(ctx)
	if err != nil {
		slog.Warn("prompt: failed to fetch programs from Directus", "error", err)
	}
	faqs, err := b.client.GetFAQ(ctx)
	if err != nil {
		slog.Warn("prompt: failed to fetch FAQ from Directus", "error", err)
	}

	var buf bytes.Buffer
	data := struct {
		Programs []directus.Program
		FAQ      []directus.FAQ
	}{programs, faqs}

	if err := tmpl.Execute(&buf, data); err != nil {
		slog.Error("prompt: template execute error", "error", err)
		return fmt.Sprintf("Ты — виртуальный ассистент ТОО «Первое кредитное товарищество». Отвечай только на русском языке.")
	}

	b.cached = buf.String()
	b.refreshed = time.Now()
	slog.Info("prompt: system prompt refreshed",
		"programs", len(programs),
		"faq", len(faqs),
		"len", len(b.cached))
	return b.cached
}
