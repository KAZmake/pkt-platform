// Package pdf generates PDF documents for the expertise service.
package pdf

import (
	"fmt"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/model"
	"github.com/jung-kurt/gofpdf"
)

// CommitteeProtocol generates a PDF protocol for a credit committee session (2.2.5).
func CommitteeProtocol(app *model.Application, votes []*model.CommitteeVote) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Протокол кредитного комитета", false)
	pdf.SetAuthor("ТОО «Первое кредитное товарищество»", false)
	pdf.AddPage()

	// Fonts — use built-in helvetica (Cyrillic via Unicode translation)
	pdf.SetFont("Helvetica", "B", 14)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Header
	pdf.CellFormat(0, 10, tr("ТОО «Первое кредитное товарищество»"), "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 8, tr("ПРОТОКОЛ КРЕДИТНОГО КОМИТЕТА"), "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf(tr("Дата: %s"), time.Now().Format("02.01.2006")), "", 1, "C", false, 0, "")
	pdf.Ln(6)

	// Application info
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(0, 8, tr("Сведения о заявке"), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetFillColor(240, 248, 255)

	row := func(label, value string) {
		pdf.CellFormat(60, 7, tr(label), "1", 0, "L", true, 0, "")
		pdf.CellFormat(0, 7, value, "1", 1, "L", false, 0, "")
	}
	row("ID заявки:", app.ID.String())
	row(tr("Сумма:"), fmt.Sprintf("%.2f KZT", app.Amount))
	row(tr("Срок (мес.):"), fmt.Sprintf("%d", app.TermMonths))
	row(tr("Тип платежа:"), app.PaymentType)
	row(tr("Статус:"), app.Status)
	pdf.Ln(6)

	// Votes table
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(0, 8, tr("Голосование членов комитета"), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(30, 77, 189)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(80, 7, tr("Эксперт (ID)"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, tr("Голос"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(0, 7, tr("Комментарий"), "1", 1, "C", true, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(0, 0, 0)
	approved, rejected, abstained := 0, 0, 0
	for _, v := range votes {
		switch v.Vote {
		case model.VoteApproved:
			approved++
			pdf.SetFillColor(240, 253, 244)
		case model.VoteRejected:
			rejected++
			pdf.SetFillColor(254, 242, 242)
		default:
			abstained++
			pdf.SetFillColor(255, 255, 255)
		}
		comment := ""
		if v.Comment != nil {
			comment = *v.Comment
		}
		pdf.CellFormat(80, 7, v.ExpertID.String()[:8]+"...", "1", 0, "L", true, 0, "")
		pdf.CellFormat(40, 7, tr(v.Vote), "1", 0, "C", true, 0, "")
		pdf.CellFormat(0, 7, comment, "1", 1, "L", true, 0, "")
	}
	pdf.Ln(6)

	// Summary
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(0, 8, tr("Итог голосования"), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf(tr("За: %d  |  Против: %d  |  Воздержались: %d"), approved, rejected, abstained), "", 1, "L", false, 0, "")

	decision := tr("ОТКАЗАНО")
	if approved > rejected {
		decision = tr("ОДОБРЕНО")
	}
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Ln(4)
	pdf.CellFormat(0, 10, fmt.Sprintf(tr("Решение комитета: %s"), decision), "1", 1, "C", false, 0, "")

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 6, tr("Документ сформирован автоматически системой PKT Platform"), "", 1, "C", false, 0, "")

	var buf []byte
	buf, err := pdfBytes(pdf)
	return buf, err
}

func pdfBytes(pdf *gofpdf.Fpdf) ([]byte, error) {
	if err := pdf.Error(); err != nil {
		return nil, fmt.Errorf("pdf error: %w", err)
	}
	// Write to in-memory buffer via custom writer
	w := &bytesWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, fmt.Errorf("pdf output: %w", err)
	}
	return w.buf, nil
}

type bytesWriter struct{ buf []byte }

func (bw *bytesWriter) Write(p []byte) (int, error) {
	bw.buf = append(bw.buf, p...)
	return len(p), nil
}
