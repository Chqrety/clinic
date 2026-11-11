package accounting

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrUnbalancedJournal    = errors.New("total debit must equal total credit")
	ErrFiscalPeriodLocked   = errors.New("accounting transaction rejected: fiscal period is closed or locked")
	ErrInvalidJournalStatus = errors.New("invalid journal status for this operation")
	ErrZeroLines            = errors.New("journal entry must have at least two balanced lines")
	ErrAccountNotFound      = errors.New("specified account does not exist or is inactive")
	ErrHeaderAccountCannotPost = errors.New("cannot post journal directly to a header/parent account")
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	uow         UnitOfWork
	accountRepo domain.AccountRepository
	journalRepo domain.JournalRepository
	periodRepo  domain.FiscalPeriodRepository
}

func NewService(
	uow UnitOfWork,
	accountRepo domain.AccountRepository,
	journalRepo domain.JournalRepository,
	periodRepo domain.FiscalPeriodRepository,
) *Service {
	return &Service{
		uow:         uow,
		accountRepo: accountRepo,
		journalRepo: journalRepo,
		periodRepo:  periodRepo,
	}
}

type CreateLineCommand struct {
	AccountID   uuid.UUID       `json:"account_id"`
	AccountCode string          `json:"account_code,omitempty"`
	Debit       decimal.Decimal `json:"debit"`
	Credit      decimal.Decimal `json:"credit"`
	Description string          `json:"description"`
}

type CreateJournalCommand struct {
	BranchID        uuid.UUID           `json:"branch_id"`
	JournalType     domain.JournalType  `json:"journal_type"`
	TransactionDate time.Time           `json:"transaction_date"`
	ReferenceType   string              `json:"reference_type"`
	ReferenceID     string              `json:"reference_id"`
	Description     string              `json:"description"`
	CreatedBy       *uuid.UUID          `json:"created_by,omitempty"`
	AutoPost        bool                `json:"auto_post"`
	Lines           []CreateLineCommand `json:"lines"`
}

// CreateJournal creates a new balanced double-entry journal entry.
func (s *Service) CreateJournal(ctx context.Context, cmd CreateJournalCommand) (*domain.JournalEntry, error) {
	if len(cmd.Lines) < 2 {
		return nil, ErrZeroLines
	}

	// 1. Validate Fiscal Period Lock
	locked, err := s.periodRepo.IsPeriodLocked(ctx, cmd.BranchID, cmd.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("check period lock: %w", err)
	}
	if locked {
		return nil, ErrFiscalPeriodLocked
	}

	// 2. Validate Lines, Accounts, and Debit/Credit Balance
	totalDebit := decimal.Zero
	totalCredit := decimal.Zero
	lines := make([]domain.JournalLine, 0, len(cmd.Lines))

	for idx, l := range cmd.Lines {
		if l.Debit.IsPositive() && l.Credit.IsPositive() {
			return nil, fmt.Errorf("line %d cannot have both debit and credit positive", idx+1)
		}
		if l.Debit.IsZero() && l.Credit.IsZero() {
			return nil, fmt.Errorf("line %d must have either debit or credit greater than zero", idx+1)
		}

		accID := l.AccountID
		if accID == uuid.Nil && l.AccountCode != "" {
			acc, err := s.accountRepo.GetByCode(ctx, l.AccountCode)
			if err != nil {
				return nil, fmt.Errorf("lookup account %s: %w", l.AccountCode, err)
			}
			accID = acc.ID
			if acc.IsHeader {
				return nil, fmt.Errorf("%w: %s", ErrHeaderAccountCannotPost, acc.Code)
			}
			if !acc.IsActive {
				return nil, fmt.Errorf("%w: %s", ErrAccountNotFound, acc.Code)
			}
		} else {
			acc, err := s.accountRepo.GetByID(ctx, accID)
			if err != nil {
				return nil, fmt.Errorf("lookup account id %s: %w", accID, err)
			}
			if acc.IsHeader {
				return nil, fmt.Errorf("%w: %s", ErrHeaderAccountCannotPost, acc.Code)
			}
			if !acc.IsActive {
				return nil, fmt.Errorf("%w: %s", ErrAccountNotFound, acc.Code)
			}
		}

		totalDebit = totalDebit.Add(l.Debit)
		totalCredit = totalCredit.Add(l.Credit)

		lines = append(lines, domain.JournalLine{
			AccountID:   accID,
			Debit:       l.Debit,
			Credit:      l.Credit,
			Description: l.Description,
			LineOrder:   idx + 1,
		})
	}

	if !decimal.IsBalanced(totalDebit, totalCredit) {
		return nil, fmt.Errorf("%w (Debit: %s, Credit: %s)", ErrUnbalancedJournal, totalDebit, totalCredit)
	}

	// 3. Generate sequential Journal Number
	journalNumber := s.generateJournalNumber(cmd.JournalType, cmd.TransactionDate)

	status := domain.JournalStatusDraft
	var postingDate *time.Time
	var postedBy *uuid.UUID

	if cmd.AutoPost {
		status = domain.JournalStatusPosted
		now := time.Now()
		postingDate = &now
		postedBy = cmd.CreatedBy
	}

	entry := &domain.JournalEntry{
		BranchScoped:    shared.BranchScoped{BranchID: cmd.BranchID},
		JournalNumber:   journalNumber,
		JournalType:     cmd.JournalType,
		TransactionDate: cmd.TransactionDate,
		PostingDate:     postingDate,
		Status:          status,
		ReferenceType:   cmd.ReferenceType,
		ReferenceID:     cmd.ReferenceID,
		Description:     cmd.Description,
		TotalDebit:      totalDebit,
		TotalCredit:     totalCredit,
		PostedBy:        postedBy,
		CreatedBy:       cmd.CreatedBy,
		Lines:           lines,
	}

	// 4. Save Header & Lines in an Atomic Unit of Work
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.journalRepo.Create(txCtx, entry)
	})
	if err != nil {
		return nil, fmt.Errorf("execute create journal uow: %w", err)
	}

	return entry, nil
}

// PostJournal posts a DRAFT journal to the General Ledger.
func (s *Service) PostJournal(ctx context.Context, journalID uuid.UUID, postedBy *uuid.UUID) (*domain.JournalEntry, error) {
	entry, err := s.journalRepo.GetByID(ctx, journalID)
	if err != nil {
		return nil, fmt.Errorf("find journal: %w", err)
	}

	if entry.Status != domain.JournalStatusDraft {
		return nil, fmt.Errorf("%w: journal status is %s, expected DRAFT", ErrInvalidJournalStatus, entry.Status)
	}

	locked, err := s.periodRepo.IsPeriodLocked(ctx, entry.BranchID, entry.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("check period lock: %w", err)
	}
	if locked {
		return nil, ErrFiscalPeriodLocked
	}

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.journalRepo.UpdateStatus(txCtx, journalID, domain.JournalStatusPosted, postedBy)
	})
	if err != nil {
		return nil, fmt.Errorf("execute post journal uow: %w", err)
	}

	entry.Status = domain.JournalStatusPosted
	now := time.Now()
	entry.PostingDate = &now
	entry.PostedBy = postedBy

	return entry, nil
}

// ReverseJournal creates an offsetting reversal journal and updates original status to REVERSED.
func (s *Service) ReverseJournal(ctx context.Context, journalID uuid.UUID, reason string, reversedBy *uuid.UUID) (*domain.JournalEntry, error) {
	orig, err := s.journalRepo.GetByID(ctx, journalID)
	if err != nil {
		return nil, fmt.Errorf("find original journal: %w", err)
	}

	if orig.Status != domain.JournalStatusPosted {
		return nil, fmt.Errorf("%w: only POSTED journals can be reversed (status is %s)", ErrInvalidJournalStatus, orig.Status)
	}

	today := time.Now()
	locked, err := s.periodRepo.IsPeriodLocked(ctx, orig.BranchID, today)
	if err != nil {
		return nil, fmt.Errorf("check period lock: %w", err)
	}
	if locked {
		return nil, ErrFiscalPeriodLocked
	}

	reversalLines := make([]CreateLineCommand, 0, len(orig.Lines))
	for _, l := range orig.Lines {
		// Invert: Debit becomes Credit, Credit becomes Debit
		reversalLines = append(reversalLines, CreateLineCommand{
			AccountID:   l.AccountID,
			Debit:       l.Credit,
			Credit:      l.Debit,
			Description: fmt.Sprintf("Reversal: %s", l.Description),
		})
	}

	reversalCmd := CreateJournalCommand{
		BranchID:        orig.BranchID,
		JournalType:     domain.JournalTypeReversal,
		TransactionDate: today,
		ReferenceType:   "JOURNAL_REVERSAL",
		ReferenceID:     orig.JournalNumber,
		Description:     fmt.Sprintf("Reversal of %s. Reason: %s", orig.JournalNumber, reason),
		CreatedBy:       reversedBy,
		AutoPost:        true,
		Lines:           reversalLines,
	}

	var reversalEntry *domain.JournalEntry
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		rev, err := s.CreateJournal(txCtx, reversalCmd)
		if err != nil {
			return err
		}
		reversalEntry = rev

		// Update original journal status to REVERSED
		return s.journalRepo.UpdateStatus(txCtx, orig.ID, domain.JournalStatusReversed, reversedBy)
	})
	if err != nil {
		return nil, fmt.Errorf("execute reversal uow: %w", err)
	}

	return reversalEntry, nil
}

// GetAccounts returns filtered list of accounts with total count.
func (s *Service) GetAccounts(ctx context.Context, filter domain.AccountFilter) ([]domain.Account, int64, error) {
	return s.accountRepo.List(ctx, filter)
}

// GetAccountByID returns an account by UUID.
func (s *Service) GetAccountByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return s.accountRepo.GetByID(ctx, id)
}

// CreateAccount creates a new account in the Chart of Accounts.
func (s *Service) CreateAccount(ctx context.Context, acc *domain.Account) error {
	return s.accountRepo.Create(ctx, acc)
}

// ListJournals returns paginated journal entries with filters.
func (s *Service) ListJournals(ctx context.Context, filter domain.JournalFilter) ([]domain.JournalEntry, int64, error) {
	return s.journalRepo.List(ctx, filter)
}

// GetJournalByID returns full journal entry with line items.
func (s *Service) GetJournalByID(ctx context.Context, id uuid.UUID) (*domain.JournalEntry, error) {
	return s.journalRepo.GetByID(ctx, id)
}

// ClosePeriod closes an accounting period, locks further modifications, and records the closing audit.
func (s *Service) ClosePeriod(ctx context.Context, branchID uuid.UUID, year, month int, closedBy *uuid.UUID, notes string) (*domain.FiscalPeriod, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	locked, err := s.periodRepo.IsPeriodLocked(ctx, branchID, startDate)
	if err != nil {
		return nil, fmt.Errorf("check period lock: %w", err)
	}
	if locked {
		return nil, fmt.Errorf("period %d-%02d is already closed or locked", year, month)
	}

	now := time.Now()
	period := &domain.FiscalPeriod{
		BaseEntity:   shared.BaseEntity{ID: uuid.New()},
		BranchScoped: shared.BranchScoped{BranchID: branchID},
		Year:         year,
		Month:        month,
		Status:       "LOCKED",
		ClosedAt:     &now,
		ClosedBy:     closedBy,
		Notes:        notes,
	}

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.periodRepo.ClosePeriod(txCtx, period)
	})
	if err != nil {
		return nil, fmt.Errorf("execute close period uow: %w", err)
	}

	return period, nil
}

// ReopenPeriod reopens a closed period with authorization note.
func (s *Service) ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error {
	return s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.periodRepo.ReopenPeriod(txCtx, periodID, reason)
	})
}

func (s *Service) generateJournalNumber(jType domain.JournalType, date time.Time) string {
	prefix := "JU"
	switch jType {
	case domain.JournalTypeCashReceipt:
		prefix = "JKM"
	case domain.JournalTypeCashDisbursement:
		prefix = "JKK"
	case domain.JournalTypeBank:
		prefix = "JB"
	case domain.JournalTypeSales:
		prefix = "JPEN"
	case domain.JournalTypePurchase:
		prefix = "JPEM"
	case domain.JournalTypeAdjusting:
		prefix = "JA"
	case domain.JournalTypeClosing:
		prefix = "JC"
	case domain.JournalTypeReversal:
		prefix = "JR"
	}

	nano := time.Now().UnixNano() % 100000
	return fmt.Sprintf("%s-%s-%05d", prefix, date.Format("20060102"), nano)
}
