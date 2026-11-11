package journal

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"clinic-api/internal/application/common"
	"clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrUnbalancedJournal   = errors.New("journal is unbalanced: total debit must equal total credit")
	ErrEmptyJournalLines   = errors.New("journal entry must have at least two valid lines")
	ErrZeroJournalAmount   = errors.New("journal total amount must be greater than zero")
	ErrAccountNotFound     = errors.New("account not found")
	ErrAccountIsHeader     = errors.New("cannot post transactions to a header account")
	ErrAccountInactive     = errors.New("cannot post transactions to an inactive account")
	ErrFiscalPeriodLocked  = errors.New("cannot post transactions to a closed or locked fiscal period")
	ErrJournalAlreadyVoid  = errors.New("cannot reverse a void or reversed journal")
	ErrBranchMismatch      = errors.New("branch ID mismatch between context and transaction payload")
)

type CreateLineCommand struct {
	AccountID   *uuid.UUID      `json:"account_id,omitempty"`
	AccountCode string          `json:"account_code,omitempty"`
	Debit       decimal.Decimal `json:"debit"`
	Credit      decimal.Decimal `json:"credit"`
	Description string          `json:"description,omitempty"`
}

type CreateJournalCommand struct {
	shared.BranchScoped
	JournalType     accounting.JournalType `json:"journal_type"`
	TransactionDate time.Time              `json:"transaction_date"`
	ReferenceType   string                 `json:"reference_type,omitempty"`
	ReferenceID     string                 `json:"reference_id,omitempty"`
	Description     string                 `json:"description"`
	Lines           []CreateLineCommand    `json:"lines"`
	CreatedBy       *uuid.UUID             `json:"created_by,omitempty"`
}

type Service struct {
	uow         common.UnitOfWork
	accountRepo accounting.AccountRepository
	journalRepo accounting.JournalRepository
	periodRepo  accounting.FiscalPeriodRepository
	seqCounter  uint64
}

func NewService(
	uow common.UnitOfWork,
	accountRepo accounting.AccountRepository,
	journalRepo accounting.JournalRepository,
	periodRepo accounting.FiscalPeriodRepository,
) *Service {
	return &Service{
		uow:         uow,
		accountRepo: accountRepo,
		journalRepo: journalRepo,
		periodRepo:  periodRepo,
		seqCounter:  uint64(time.Now().UnixNano() % 100000),
	}
}

// CreateJournal validates double-entry balance, verifies period lock status, resolves accounts, and persists.
func (s *Service) CreateJournal(ctx context.Context, cmd CreateJournalCommand) (*accounting.JournalEntry, error) {
	if len(cmd.Lines) < 2 {
		return nil, ErrEmptyJournalLines
	}

	if cmd.TransactionDate.IsZero() {
		cmd.TransactionDate = time.Now()
	}

	// 1. Verify Fiscal Period lock status
	isLocked, err := s.periodRepo.IsPeriodLocked(ctx, cmd.BranchID, cmd.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("check fiscal period lock: %w", err)
	}
	if isLocked {
		return nil, ErrFiscalPeriodLocked
	}

	// 2. Validate double-entry balance: sum(debit) == sum(credit) > 0
	totalDebit := decimal.Zero
	totalCredit := decimal.Zero

	var resolvedLines []accounting.JournalLine
	for idx, lineCmd := range cmd.Lines {
		if lineCmd.Debit.IsNegative() || lineCmd.Credit.IsNegative() {
			return nil, fmt.Errorf("line [%d]: negative debit or credit amount is prohibited", idx)
		}
		if lineCmd.Debit.IsZero() && lineCmd.Credit.IsZero() {
			return nil, fmt.Errorf("line [%d]: both debit and credit are zero", idx)
		}
		if lineCmd.Debit.IsPositive() && lineCmd.Credit.IsPositive() {
			return nil, fmt.Errorf("line [%d]: line cannot have both positive debit and credit", idx)
		}

		// Resolve Account
		var acc *accounting.Account
		if lineCmd.AccountID != nil && *lineCmd.AccountID != uuid.Nil {
			acc, err = s.accountRepo.GetByID(ctx, *lineCmd.AccountID)
		} else if lineCmd.AccountCode != "" {
			acc, err = s.accountRepo.GetByCode(ctx, lineCmd.AccountCode)
		} else {
			return nil, fmt.Errorf("line [%d]: neither account_id nor account_code was provided", idx)
		}

		if err != nil {
			return nil, fmt.Errorf("line [%d] resolve account: %w", idx, ErrAccountNotFound)
		}
		if acc.IsHeader {
			return nil, fmt.Errorf("line [%d] account %s: %w", idx, acc.Code, ErrAccountIsHeader)
		}
		if !acc.IsActive {
			return nil, fmt.Errorf("line [%d] account %s: %w", idx, acc.Code, ErrAccountInactive)
		}

		totalDebit = totalDebit.Add(lineCmd.Debit)
		totalCredit = totalCredit.Add(lineCmd.Credit)

		resolvedLines = append(resolvedLines, accounting.JournalLine{
			ID:          uuid.New(),
			AccountID:   acc.ID,
			AccountCode: acc.Code,
			AccountName: acc.Name,
			Debit:       lineCmd.Debit,
			Credit:      lineCmd.Credit,
			Description: lineCmd.Description,
			LineOrder:   idx + 1,
		})
	}

	if totalDebit.IsZero() || totalCredit.IsZero() {
		return nil, ErrZeroJournalAmount
	}

	if !totalDebit.Equal(totalCredit) {
		return nil, fmt.Errorf("%w (Total Debit: %s, Total Credit: %s)", ErrUnbalancedJournal, totalDebit.StringFixed(2), totalCredit.StringFixed(2))
	}

	journalNumber := s.GenerateJournalNumber(cmd.JournalType, cmd.TransactionDate)
	now := time.Now()

	entry := &accounting.JournalEntry{
		BaseEntity: shared.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: cmd.CreatedBy,
		},
		BranchScoped: shared.BranchScoped{
			BranchID: cmd.BranchID,
		},
		JournalNumber:   journalNumber,
		JournalType:     cmd.JournalType,
		TransactionDate: cmd.TransactionDate,
		PostingDate:     &now,
		Status:          accounting.JournalStatusPosted,
		ReferenceType:   cmd.ReferenceType,
		ReferenceID:     cmd.ReferenceID,
		TotalDebit:      totalDebit,
		TotalCredit:     totalCredit,
		Description:     cmd.Description,
		PostedBy:        cmd.CreatedBy,
		Lines:           resolvedLines,
	}

	// 3. Atomically persist journal
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.journalRepo.Create(txCtx, entry)
	})
	if err != nil {
		return nil, fmt.Errorf("persist journal: %w", err)
	}

	return entry, nil
}

// ReverseJournal creates an exact inverted journal entry to void a posted journal with an audit trail.
func (s *Service) ReverseJournal(ctx context.Context, id uuid.UUID, reason string, createdBy *uuid.UUID) (*accounting.JournalEntry, error) {
	orig, err := s.journalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find original journal: %w", err)
	}

	if orig.Status == accounting.JournalStatusReversed || orig.Status == accounting.JournalStatusVoid {
		return nil, ErrJournalAlreadyVoid
	}

	isLocked, err := s.periodRepo.IsPeriodLocked(ctx, orig.BranchID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("check period lock: %w", err)
	}
	if isLocked {
		return nil, ErrFiscalPeriodLocked
	}

	var revLines []CreateLineCommand
	for _, l := range orig.Lines {
		revLines = append(revLines, CreateLineCommand{
			AccountID:   &l.AccountID,
			AccountCode: l.AccountCode,
			Debit:       l.Credit,
			Credit:      l.Debit,
			Description: fmt.Sprintf("[Pembalikan %s] %s", orig.JournalNumber, l.Description),
		})
	}

	revCmd := CreateJournalCommand{
		BranchScoped:    shared.BranchScoped{BranchID: orig.BranchID},
		JournalType:     accounting.JournalTypeReversal,
		TransactionDate: time.Now(),
		ReferenceType:   "REVERSAL",
		ReferenceID:     orig.JournalNumber,
		Description:     fmt.Sprintf("Pembalikan Jurnal %s - Alasan: %s", orig.JournalNumber, reason),
		Lines:           revLines,
		CreatedBy:       createdBy,
	}

	var revEntry *accounting.JournalEntry
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		re, err := s.CreateJournal(txCtx, revCmd)
		if err != nil {
			return err
		}
		revEntry = re

		return s.journalRepo.UpdateStatus(txCtx, orig.ID, accounting.JournalStatusReversed, createdBy)
	})

	if err != nil {
		return nil, fmt.Errorf("execute reversal unit of work: %w", err)
	}

	return revEntry, nil
}

func (s *Service) GetJournalByID(ctx context.Context, id uuid.UUID) (*accounting.JournalEntry, error) {
	return s.journalRepo.GetByID(ctx, id)
}

func (s *Service) GetJournalByNumber(ctx context.Context, number string) (*accounting.JournalEntry, error) {
	return s.journalRepo.GetByJournalNumber(ctx, number)
}

func (s *Service) ListJournals(ctx context.Context, filter accounting.JournalFilter) ([]accounting.JournalEntry, int64, error) {
	return s.journalRepo.List(ctx, filter)
}

func (s *Service) GenerateJournalNumber(jType accounting.JournalType, date time.Time) string {
	prefix := "JU"
	switch jType {
	case accounting.JournalTypeCashReceipt:
		prefix = "JKM"
	case accounting.JournalTypeCashDisbursement:
		prefix = "JKK"
	case accounting.JournalTypeBank:
		prefix = "JB"
	case accounting.JournalTypeSales:
		prefix = "JJ"
	case accounting.JournalTypePurchase:
		prefix = "JB"
	case accounting.JournalTypeAdjusting:
		prefix = "JC"
	case accounting.JournalTypeReversal:
		prefix = "JR"
	case accounting.JournalTypeClosing:
		prefix = "JT"
	}

	seq := atomic.AddUint64(&s.seqCounter, 1) % 10000
	return fmt.Sprintf("%s-%s-%04d", prefix, date.Format("20060102"), seq)
}
