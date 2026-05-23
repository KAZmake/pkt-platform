package pdf

import (
	"fmt"
	"math"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/model"
	"github.com/jung-kurt/gofpdf"
)

// LoanAgreement generates a PDF loan agreement with payment schedule (2.2.6).
func LoanAgreement(app *model.Application, annualRate float64) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Договор займа", false)
	pdf.SetAuthor("ТОО «Первое кредитное товарищество»", false)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// ── Header ────────────────────────────────────────────────────────────────
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(0, 10, tr("ДОГОВОР ЗАЙМА"), "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf(tr("г. Уральск, %s"), time.Now().Format("02 января 2006 г.")), "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// ── Parties ───────────────────────────────────────────────────────────────
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(0, 8, tr("СТОРОНЫ ДОГОВОРА"), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.MultiCell(0, 6, tr("Займодавец: ТОО «Первое кредитное товарищество», БИН 123456789012, г. Уральск"), "", "L", false)
	pdf.MultiCell(0, 6, fmt.Sprintf(tr("Заёмщик: (ID %s)"), app.BorrowerID.String()), "", "L", false)
	pdf.Ln(4)

	// ── Terms ─────────────────────────────────────────────────────────────────
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(0, 8, tr("УСЛОВИЯ ЗАЙМА"), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetFillColor(240, 248, 255)
	termRow := func(l, v string) {
		pdf.CellFormat(70, 7, tr(l), "1", 0, "L", true, 0, "")
		pdf.CellFormat(0, 7, v, "1", 1, "L", false, 0, "")
	}
	termRow("ID заявки:", app.ID.String())
	termRow(tr("Сумма займа:"), fmt.Sprintf("%.2f тенге", app.Amount))
	termRow(tr("Срок:"), fmt.Sprintf("%d мес.", app.TermMonths))
	termRow(tr("Процентная ставка (год.):"), fmt.Sprintf("%.2f%%", annualRate))
	termRow(tr("Тип платежа:"), app.PaymentType)
	pdf.Ln(6)

	// ── Payment schedule ──────────────────────────────────────────────────────
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(0, 8, tr("ГРАФИК ПЛАТЕЖЕЙ"), "", 1, "L", false, 0, "")

	schedule := buildSchedule(app.Amount, annualRate, app.TermMonths, app.PaymentType, app.CreatedAt)

	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetFillColor(30, 77, 189)
	pdf.SetTextColor(255, 255, 255)
	for _, h := range []string{"#", tr("Дата"), tr("Основной долг"), tr("Проценты"), tr("Платёж"), tr("Остаток")} {
		w := map[string]float64{"#": 10, tr("Дата"): 28, tr("Основной долг"): 35, tr("Проценты"): 28, tr("Платёж"): 30, tr("Остаток"): 0}[h]
		if h == tr("Остаток") {
			w = 0
		}
		pdf.CellFormat(w, 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(0, 0, 0)
	for i, row := range schedule {
		if i%2 == 0 {
			pdf.SetFillColor(240, 248, 255)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.CellFormat(10, 6, fmt.Sprintf("%d", row.month), "1", 0, "C", true, 0, "")
		pdf.CellFormat(28, 6, row.date, "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 6, fmt.Sprintf("%.2f", row.principal), "1", 0, "R", true, 0, "")
		pdf.CellFormat(28, 6, fmt.Sprintf("%.2f", row.interest), "1", 0, "R", true, 0, "")
		pdf.CellFormat(30, 6, fmt.Sprintf("%.2f", row.payment), "1", 0, "R", true, 0, "")
		pdf.CellFormat(0, 6, fmt.Sprintf("%.2f", row.balance), "1", 1, "R", true, 0, "")

		if pdf.GetY() > 270 {
			pdf.AddPage()
			pdf.SetFont("Helvetica", "", 8)
		}
	}

	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 6, tr("Подписи сторон: _______________________ (Займодавец)   _______________________ (Заёмщик)"), "", 1, "C", false, 0, "")

	return pdfBytes(pdf)
}

type scheduleRow struct {
	month     int
	date      string
	principal float64
	interest  float64
	payment   float64
	balance   float64
}

func buildSchedule(principal, annualRate float64, months int, payType string, start time.Time) []scheduleRow {
	monthlyRate := annualRate / 100 / 12
	balance := principal
	rows := make([]scheduleRow, 0, months)

	switch payType {
	case "annuity":
		var mp float64
		if monthlyRate == 0 {
			mp = principal / float64(months)
		} else {
			f := math.Pow(1+monthlyRate, float64(months))
			mp = principal * monthlyRate * f / (f - 1)
		}
		mp = r2(mp)
		for i := 1; i <= months; i++ {
			interest := r2(balance * monthlyRate)
			pmt := mp
			if i == months {
				pmt = r2(balance + interest)
			}
			pmtPrincipal := r2(pmt - interest)
			balance = r2(balance - pmtPrincipal)
			rows = append(rows, scheduleRow{i, payDate(start, i), pmtPrincipal, interest, pmt, math.Max(balance, 0)})
		}
	default: // differentiated
		pp := r2(principal / float64(months))
		for i := 1; i <= months; i++ {
			interest := r2(balance * monthlyRate)
			p := pp
			if i == months {
				p = balance
			}
			pmt := r2(p + interest)
			balance = r2(balance - p)
			rows = append(rows, scheduleRow{i, payDate(start, i), p, interest, pmt, math.Max(balance, 0)})
		}
	}
	return rows
}

func r2(v float64) float64 { return math.Round(v*100) / 100 }
func payDate(start time.Time, offset int) string {
	d := start.AddDate(0, offset, 0)
	return fmt.Sprintf("%02d.%02d.%04d", 1, d.Month(), d.Year())
}
