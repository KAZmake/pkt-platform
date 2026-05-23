package service

import (
	"bytes"
	"fmt"
	"html/template"

	resend "github.com/resend/resend-go/v2"
)

// Mailer wraps the Resend client with HTML email templates.
type Mailer struct {
	client *resend.Client
	from   string
}

func NewMailer(apiKey, from string) *Mailer {
	return &Mailer{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

// Send sends an HTML email.
func (m *Mailer) Send(to, subject, htmlBody string) error {
	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}
	_, err := m.client.Emails.Send(params)
	return err
}

// ── Email templates ───────────────────────────────────────────────────────────

var applicationCreatedTmpl = template.Must(template.New("app_created").Parse(`<!DOCTYPE html>
<html lang="ru"><head><meta charset="UTF-8"><title>Заявка принята</title></head>
<body style="font-family:sans-serif;color:#1a1a1a;max-width:600px;margin:0 auto;padding:24px">
<h2 style="color:#166534">Ваша заявка принята</h2>
<p>Заявка <strong>#{{.ApplicationID}}</strong> успешно зарегистрирована.</p>
<table style="width:100%;border-collapse:collapse;margin:16px 0">
  <tr><td style="padding:8px;background:#f0fdf4;font-weight:bold">Сумма</td><td style="padding:8px">{{printf "%.2f" .Amount}} ₸</td></tr>
  <tr><td style="padding:8px;background:#f0fdf4;font-weight:bold">Срок</td><td style="padding:8px">{{.TermMonths}} мес.</td></tr>
</table>
<p>Следующий шаг — первичный скоринг. Мы уведомим вас о каждом изменении статуса.</p>
<p style="color:#6b7280;font-size:12px">ТОО «Первое кредитное товарищество»</p>
</body></html>`))

var statusChangedTmpl = template.Must(template.New("status_changed").Parse(`<!DOCTYPE html>
<html lang="ru"><head><meta charset="UTF-8"><title>Статус заявки изменён</title></head>
<body style="font-family:sans-serif;color:#1a1a1a;max-width:600px;margin:0 auto;padding:24px">
<h2 style="color:#1d4ed8">Статус заявки обновлён</h2>
<p>Статус заявки <strong>#{{.ApplicationID}}</strong> изменён:</p>
<table style="width:100%;border-collapse:collapse;margin:16px 0">
  <tr><td style="padding:8px;background:#eff6ff;font-weight:bold">Было</td><td style="padding:8px">{{.FromStatus}}</td></tr>
  <tr><td style="padding:8px;background:#eff6ff;font-weight:bold">Стало</td><td style="padding:8px;color:#1d4ed8;font-weight:bold">{{.ToStatus}}</td></tr>
</table>
<p>Войдите в <a href="{{.CabinetURL}}" style="color:#1d4ed8">личный кабинет</a>, чтобы узнать подробности.</p>
<p style="color:#6b7280;font-size:12px">ТОО «Первое кредитное товарищество»</p>
</body></html>`))

type AppCreatedEmailData struct {
	ApplicationID string
	Amount        float64
	TermMonths    int
}

type StatusChangedEmailData struct {
	ApplicationID string
	FromStatus    string
	ToStatus      string
	CabinetURL    string
}

func (m *Mailer) SendApplicationCreated(to string, data AppCreatedEmailData) error {
	var buf bytes.Buffer
	if err := applicationCreatedTmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	return m.Send(to, "Заявка принята — ТОО «ПКТ»", buf.String())
}

func (m *Mailer) SendStatusChanged(to string, data StatusChangedEmailData) error {
	var buf bytes.Buffer
	if err := statusChangedTmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	subject := fmt.Sprintf("Статус заявки #%s изменён на %s", data.ApplicationID, data.ToStatus)
	return m.Send(to, subject, buf.String())
}
