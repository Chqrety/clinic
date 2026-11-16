package journal

import (
	"context"
	"fmt"
	"log"
	"time"

	cashierDomain "clinic-api/internal/cashier/domain"
	"clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	BranchKidinestID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	BranchSSCID      = uuid.MustParse("a0000000-0000-0000-0000-000000000002")
)

// CashierJournalListener handles automatic asynchronous posting from cashier settlements to the journal.
type CashierJournalListener struct {
	journalSvc *Service
}

func NewCashierJournalListener(journalSvc *Service) *CashierJournalListener {
	return &CashierJournalListener{journalSvc: journalSvc}
}

// OnInvoiceSettled records double-entry: (Debit) Kas Kasir / Bank EDC vs (Credit) Pendapatan Jasa Medis & Obat.
func (l *CashierJournalListener) OnInvoiceSettled(ctx context.Context, inv *cashierDomain.Invoice) error {
	if inv.TotalAmount <= 0 {
		return nil
	}

	branchID := BranchKidinestID
	cashAccount := "1-1111" // Kas Kasir Kidinest
	if inv.ClinicUnit == "SSC" {
		branchID = BranchSSCID
		cashAccount = "1-1112" // Kas Kasir SSC
	}

	// Tentukan Akun Debit (Kas Kasir atau Bank BCA Operasional / EDC)
	debitAccount := cashAccount
	if inv.PaymentMethod != cashierDomain.PaymentMethodCash {
		debitAccount = "1-1113" // Bank BCA Operasional (EDC / Transfer)
	}

	totalNet := decimal.FromInt(int64(inv.TotalAmount))

	lines := []CreateLineCommand{
		{
			AccountCode: debitAccount,
			Debit:       totalNet,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Penerimaan %s: Inv %s - %s", inv.PaymentMethod, inv.InvoiceNumber, inv.PatientName),
		},
	}

	// Alokasikan Kredit ke Pendapatan Pendaftaran, Obat, dan Jasa Medis
	allocated := 0

	// 1. Pendapatan Pendaftaran
	if inv.SubtotalRegistration > 0 {
		regAmt := inv.SubtotalRegistration
		if regAmt > inv.TotalAmount {
			regAmt = inv.TotalAmount
		}
		allocated += regAmt
		lines = append(lines, CreateLineCommand{
			AccountCode: "4-1100", // Pendapatan Pendaftaran & Skrining Anak
			Debit:       decimal.Zero,
			Credit:      decimal.FromInt(int64(regAmt)),
			Description: fmt.Sprintf("Registrasi: %s", inv.InvoiceNumber),
		})
	}

	// 2. Pendapatan Obat & Skincare
	medCode := "4-1500" // Pendapatan Obat Anak (Kidinest)
	if inv.ClinicUnit == "SSC" {
		medCode = "4-2400" // Pendapatan Skincare & Produk SSC
	}
	if inv.SubtotalMedication > 0 && allocated < inv.TotalAmount {
		medAmt := inv.SubtotalMedication
		if allocated+medAmt > inv.TotalAmount {
			medAmt = inv.TotalAmount - allocated
		}
		allocated += medAmt
		lines = append(lines, CreateLineCommand{
			AccountCode: medCode,
			Debit:       decimal.Zero,
			Credit:      decimal.FromInt(int64(medAmt)),
			Description: fmt.Sprintf("Pendapatan Obat/Skincare: %s", inv.InvoiceNumber),
		})
	}

	// 3. Pendapatan Jasa Medis / Tindakan (sisanya)
	actCode := "4-1200" // Pendapatan Jasa Dokter Spesialis Anak (Sp.A)
	if inv.ClinicUnit == "SSC" {
		actCode = "4-2200" // Pendapatan Treatment Facial & Estetika SSC
	}
	remaining := inv.TotalAmount - allocated
	if remaining > 0 {
		lines = append(lines, CreateLineCommand{
			AccountCode: actCode,
			Debit:       decimal.Zero,
			Credit:      decimal.FromInt(int64(remaining)),
			Description: fmt.Sprintf("Pendapatan Tindakan/Jasa Medis: %s", inv.InvoiceNumber),
		})
	} else if len(lines) == 1 {
		// Fallback seimbang jika subtotal per item 0
		lines = append(lines, CreateLineCommand{
			AccountCode: actCode,
			Debit:       decimal.Zero,
			Credit:      totalNet,
			Description: fmt.Sprintf("Pendapatan Layanan Medis: %s", inv.InvoiceNumber),
		})
	}

	cmd := CreateJournalCommand{
		BranchScoped: shared.BranchScoped{
			BranchID: branchID,
		},
		JournalType:     accounting.JournalTypeSales,
		TransactionDate: time.Now(),
		ReferenceType:   "CASHIER_INVOICE",
		ReferenceID:     inv.InvoiceNumber,
		Description:     fmt.Sprintf("Pendapatan Kasir [%s] %s - Pasien: %s", inv.ClinicUnit, inv.InvoiceNumber, inv.PatientName),
		Lines:           lines,
	}

	entry, err := l.journalSvc.CreateJournal(ctx, cmd)
	if err != nil {
		return fmt.Errorf("create sales journal for inv %s: %w", inv.InvoiceNumber, err)
	}

	log.Printf("[CASHIER-JOURNAL-SYNC] Jurnal %s berhasil dicatat untuk Invoice %s senilai Rp %d", entry.JournalNumber, inv.InvoiceNumber, inv.TotalAmount)
	return nil
}

// OnShiftClosed mencatat ringkasan tutup shift ke log akuntansi.
func (l *CashierJournalListener) OnShiftClosed(ctx context.Context, report *cashierDomain.CashierShiftReport) error {
	log.Printf("[CASHIER-JOURNAL-SYNC] Tutup shift dicatat ke akuntansi: Kasir %s, Tgl %s, Net: Rp %d", report.StaffName, report.ShiftDate, report.NetTotal)
	return nil
}
