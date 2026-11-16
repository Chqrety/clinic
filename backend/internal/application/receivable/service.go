package receivable

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"clinic-api/internal/application/common"
	journalApp "clinic-api/internal/application/journal"
	"clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/receivable"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrReceivableNotFound     = errors.New("receivable record not found")
	ErrReceivableAlreadyPaid  = errors.New("receivable is already fully settled")
	ErrInvalidPaymentAmount   = errors.New("payment amount must be greater than zero and cannot exceed outstanding amount")
	ErrZeroAmount             = errors.New("receivable total amount must be positive")
)

type CreateReceivableCommand struct {
	shared.BranchScoped
	CustomerType     receivable.CustomerType `json:"customer_type"`
	CustomerID       *uuid.UUID              `json:"customer_id,omitempty"`
	CustomerName     string                  `json:"customer_name"`
	InvoiceReference string                  `json:"invoice_reference"`
	TotalAmount      decimal.Decimal         `json:"total_amount"`
	IssueDate        time.Time               `json:"issue_date"`
	DueDate          time.Time               `json:"due_date"`
	AccountCode      string                  `json:"account_code"` // e.g. 1-1200 / 1-1210 / 1-1220
	Notes            string                  `json:"notes,omitempty"`
}

type PayReceivableCommand struct {
	ReceivableID        uuid.UUID       `json:"receivable_id"`
	Amount              decimal.Decimal `json:"amount"`
	PaymentDate         time.Time       `json:"payment_date"`
	CashBankAccountCode string          `json:"cash_bank_account_code"` // e.g. 1-1111 (Kas Kasir Kidinest), 1-1121 (Bank BCA)
	ReferenceNo         string          `json:"reference_no,omitempty"`
	Notes               string          `json:"notes,omitempty"`
	ReceivedBy          *uuid.UUID      `json:"received_by,omitempty"`
}

type Service struct {
	uow            common.UnitOfWork
	receivableRepo receivable.Repository
	accountRepo    accounting.AccountRepository
	journalService *journalApp.Service
	seqCounter     uint64
}

func NewService(
	uow common.UnitOfWork,
	receivableRepo receivable.Repository,
	accountRepo accounting.AccountRepository,
	journalService *journalApp.Service,
) *Service {
	return &Service{
		uow:            uow,
		receivableRepo: receivableRepo,
		accountRepo:    accountRepo,
		journalService: journalService,
		seqCounter:     uint64(time.Now().UnixNano() % 100000),
	}
}

// CreateReceivable registers a new patient or insurance claim receivable.
func (s *Service) CreateReceivable(ctx context.Context, cmd CreateReceivableCommand) (*receivable.Receivable, error) {
	if !cmd.TotalAmount.IsPositive() {
		return nil, ErrZeroAmount
	}
	if cmd.IssueDate.IsZero() {
		cmd.IssueDate = time.Now()
	}
	if cmd.DueDate.IsZero() {
		cmd.DueDate = cmd.IssueDate.AddDate(0, 0, 30) // Default net 30 days
	}

	accCode := cmd.AccountCode
	if accCode == "" {
		accCode = "1-1200" // Default Piutang Usaha & Klaim
	}

	acc, err := s.accountRepo.GetByCode(ctx, accCode)
	if err != nil {
		return nil, fmt.Errorf("lookup receivable account %s: %w", accCode, err)
	}

	now := time.Now()
	recNumber := s.generateReceivableNumber(cmd.IssueDate)

	rec := &receivable.Receivable{
		BaseEntity: shared.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		BranchScoped: shared.BranchScoped{
			BranchID: cmd.BranchID,
		},
		ReceivableNumber:  recNumber,
		CustomerType:      cmd.CustomerType,
		CustomerID:        cmd.CustomerID,
		CustomerName:      cmd.CustomerName,
		InvoiceReference:  cmd.InvoiceReference,
		TotalAmount:       cmd.TotalAmount,
		OutstandingAmount: cmd.TotalAmount,
		IssueDate:         cmd.IssueDate,
		DueDate:           cmd.DueDate,
		Status:            receivable.StatusOpen,
		AccountID:         acc.ID,
		AccountCode:       acc.Code,
		Notes:             cmd.Notes,
	}

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.receivableRepo.Create(txCtx, rec)
	})
	if err != nil {
		return nil, fmt.Errorf("persist receivable: %w", err)
	}

	return rec, nil
}

// PayReceivable settles all or part of an outstanding receivable and auto-posts a Cash Receipt Journal (JKM).
func (s *Service) PayReceivable(ctx context.Context, cmd PayReceivableCommand) (*receivable.Payment, error) {
	if !cmd.Amount.IsPositive() {
		return nil, ErrInvalidPaymentAmount
	}
	if cmd.PaymentDate.IsZero() {
		cmd.PaymentDate = time.Now()
	}

	rec, err := s.receivableRepo.GetByID(ctx, cmd.ReceivableID)
	if err != nil {
		return nil, ErrReceivableNotFound
	}

	if rec.Status == receivable.StatusPaid || rec.OutstandingAmount.IsZero() {
		return nil, ErrReceivableAlreadyPaid
	}

	if cmd.Amount.GreaterThan(rec.OutstandingAmount) {
		return nil, fmt.Errorf("%w: payment %s exceeds outstanding %s",
			ErrInvalidPaymentAmount, cmd.Amount.StringFixed(2), rec.OutstandingAmount.StringFixed(2))
	}

	newOutstanding := rec.OutstandingAmount.Sub(cmd.Amount)
	newStatus := receivable.StatusPartial
	if newOutstanding.IsZero() {
		newStatus = receivable.StatusPaid
	}

	payNumber := s.generatePaymentNumber(cmd.PaymentDate)

	cashAccountCode := cmd.CashBankAccountCode
	if cashAccountCode == "" {
		cashAccountCode = "1-1111" // Kas Kasir Kidinest
	}

	payment := &receivable.Payment{
		ID:                  uuid.New(),
		ReceivableID:        rec.ID,
		PaymentNumber:       payNumber,
		Amount:              cmd.Amount,
		PaymentDate:         cmd.PaymentDate,
		CashBankAccountCode: cashAccountCode,
		ReferenceNo:         cmd.ReferenceNo,
		Notes:               cmd.Notes,
		CreatedBy:           cmd.ReceivedBy,
		CreatedAt:           time.Now(),
	}

	// Double-Entry Journal:
	// Debit: Kas Kasir / Bank BCA (e.g. 1-1111 or 1-1121)
	// Kredit: Piutang Pasien / Klaim (e.g. 1-1200 or 1-1210)
	journalCmd := journalApp.CreateJournalCommand{
		BranchScoped:    shared.BranchScoped{BranchID: rec.BranchID},
		JournalType:     accounting.JournalTypeCashReceipt,
		TransactionDate: cmd.PaymentDate,
		ReferenceType:   "RECEIVABLE_PAYMENT",
		ReferenceID:     rec.ReceivableNumber,
		Description:     fmt.Sprintf("Pelunasan Piutang %s - Pasien/Klaim: %s", rec.ReceivableNumber, rec.CustomerName),
		CreatedBy:       cmd.ReceivedBy,
		Lines: []journalApp.CreateLineCommand{
			{
				AccountCode: cashAccountCode,
				Debit:       cmd.Amount,
				Credit:      decimal.Zero,
				Description: fmt.Sprintf("Penerimaan Kas/Bank dari %s", rec.CustomerName),
			},
			{
				AccountCode: rec.AccountCode,
				Debit:       decimal.Zero,
				Credit:      cmd.Amount,
				Description: fmt.Sprintf("Pengurangan Piutang Invoice %s", rec.InvoiceReference),
			},
		},
	}

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		// 1. Post JKM
		jEntry, err := s.journalService.CreateJournal(txCtx, journalCmd)
		if err != nil {
			return fmt.Errorf("post receivable cash receipt journal: %w", err)
		}
		payment.JournalID = &jEntry.ID

		// 2. Persist Payment Record
		if err := s.receivableRepo.AddPayment(txCtx, payment); err != nil {
			return fmt.Errorf("record receivable payment: %w", err)
		}

		// 3. Update Outstanding & Status
		if err := s.receivableRepo.UpdateOutstanding(txCtx, rec.ID, newOutstanding, newStatus); err != nil {
			return fmt.Errorf("update receivable status: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("execute payment unit of work: %w", err)
	}

	return payment, nil
}

func (s *Service) GetReceivableByID(ctx context.Context, id uuid.UUID) (*receivable.Receivable, error) {
	return s.receivableRepo.GetByID(ctx, id)
}

func (s *Service) GetReceivableByNumber(ctx context.Context, number string) (*receivable.Receivable, error) {
	return s.receivableRepo.GetByNumber(ctx, number)
}

func (s *Service) ListReceivables(ctx context.Context, filter receivable.Filter) ([]receivable.Receivable, int64, error) {
	return s.receivableRepo.List(ctx, filter)
}

func (s *Service) GetAgingReport(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*receivable.AgingReport, error) {
	if asOfDate.IsZero() {
		asOfDate = time.Now()
	}
	return s.receivableRepo.GetAgingReport(ctx, branchID, asOfDate)
}

func (s *Service) generateReceivableNumber(date time.Time) string {
	seq := atomic.AddUint64(&s.seqCounter, 1) % 10000
	return fmt.Sprintf("AR-%s-%04d", date.Format("20060102"), seq)
}

func (s *Service) generatePaymentNumber(date time.Time) string {
	seq := atomic.AddUint64(&s.seqCounter, 1) % 10000
	return fmt.Sprintf("PAY-AR-%s-%04d", date.Format("20060102"), seq)
}

