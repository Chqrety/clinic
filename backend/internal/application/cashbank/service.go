package cashbank

import (
	"context"
	"errors"
	"fmt"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/cashbank"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrAccountNotFound       = errors.New("cash/bank account not found")
	ErrCounterpartNotFound   = errors.New("counterpart account not found in Chart of Accounts")
	ErrInvalidAmount         = errors.New("transaction amount must be positive")
	ErrSameAccountTransfer   = errors.New("source and destination accounts cannot be identical")
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	uow               UnitOfWork
	repo              cashbank.Repository
	accountingService *accountingApp.Service
}

func NewService(
	uow UnitOfWork,
	repo cashbank.Repository,
	accountingService *accountingApp.Service,
) *Service {
	return &Service{
		uow:               uow,
		repo:              repo,
		accountingService: accountingService,
	}
}

type RecordCashInCommand struct {
	BranchID             uuid.UUID       `json:"branch_id"`
	CashBankAccountID    uuid.UUID       `json:"cash_bank_account_id"`
	Amount               decimal.Decimal `json:"amount"`
	TransactionDate      time.Time       `json:"transaction_date"`
	CounterpartAccountID uuid.UUID       `json:"counterpart_account_id"` // Kredit: Pendapatan Lain / Ekuitas / Piutang
	Description          string          `json:"description"`
	ReferenceNo          string          `json:"reference_no,omitempty"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
}

type RecordCashOutCommand struct {
	BranchID             uuid.UUID       `json:"branch_id"`
	CashBankAccountID    uuid.UUID       `json:"cash_bank_account_id"`
	Amount               decimal.Decimal `json:"amount"`
	TransactionDate      time.Time       `json:"transaction_date"`
	CounterpartAccountID uuid.UUID       `json:"counterpart_account_id"` // Debit: Beban Listrik / Beban Operasional / Hutang
	Description          string          `json:"description"`
	ReferenceNo          string          `json:"reference_no,omitempty"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
}

type TransferCommand struct {
	BranchID        uuid.UUID       `json:"branch_id"`
	SourceAccountID uuid.UUID       `json:"source_account_id"` // Sumber (e.g. Kas Kasir)
	TargetAccountID uuid.UUID       `json:"target_account_id"` // Tujuan (e.g. Bank BCA)
	Amount          decimal.Decimal `json:"amount"`
	TransactionDate time.Time       `json:"transaction_date"`
	Description     string          `json:"description"`
	ReferenceNo     string          `json:"reference_no,omitempty"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
}

// RecordCashIn records a Cash Receipt (JKM) and automatically posts a balanced double-entry journal.
func (s *Service) RecordCashIn(ctx context.Context, cmd RecordCashInCommand) (*cashbank.CashBankTransaction, *accountingDomain.JournalEntry, error) {
	if !cmd.Amount.IsPositive() {
		return nil, nil, ErrInvalidAmount
	}
	if cmd.TransactionDate.IsZero() {
		cmd.TransactionDate = time.Now()
	}

	cbAcc, err := s.repo.GetAccountByID(ctx, cmd.CashBankAccountID)
	if err != nil {
		return nil, nil, ErrAccountNotFound
	}

	counterAcc, err := s.accountingService.GetAccountByID(ctx, cmd.CounterpartAccountID)
	if err != nil {
		return nil, nil, ErrCounterpartNotFound
	}

	txNumber := fmt.Sprintf("BKM-%s-%04d", cmd.TransactionDate.Format("20060102"), time.Now().UnixNano()%10000)

	// Double-entry: Debit Cash/Bank Account, Credit Counterpart Account
	journalLines := []accountingApp.CreateLineCommand{
		{
			AccountID:   cbAcc.AccountID,
			Debit:       cmd.Amount,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Kas Masuk: %s", cmd.Description),
		},
		{
			AccountID:   counterAcc.ID,
			Debit:       decimal.Zero,
			Credit:      cmd.Amount,
			Description: fmt.Sprintf("Lawan Kas Masuk: %s", counterAcc.Name),
		},
	}

	journalType := accountingDomain.JournalTypeCashReceipt
	if cbAcc.Type == cashbank.AccountTypeBank {
		journalType = accountingDomain.JournalTypeBank
	}

	var journalEntry *accountingDomain.JournalEntry
	txItem := &cashbank.CashBankTransaction{
		BaseEntity:           shared.BaseEntity{ID: uuid.New()},
		BranchScoped:         shared.BranchScoped{BranchID: cmd.BranchID},
		CashBankAccountID:    cbAcc.ID,
		TransactionNumber:    txNumber,
		TransactionType:      cashbank.TransactionTypeCashIn,
		Amount:               cmd.Amount,
		TransactionDate:      cmd.TransactionDate,
		CounterpartAccountID: counterAcc.ID,
		Description:          cmd.Description,
		ReferenceNo:          cmd.ReferenceNo,
	}

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		// Post journal
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        cmd.BranchID,
			JournalType:     journalType,
			TransactionDate: cmd.TransactionDate,
			ReferenceType:   "CASH_IN",
			ReferenceID:     txNumber,
			Description:     cmd.Description,
			CreatedBy:       cmd.CreatedBy,
			AutoPost:        true,
			Lines:           journalLines,
		})
		if err != nil {
			return fmt.Errorf("auto-post cash in journal: %w", err)
		}
		journalEntry = jEntry
		txItem.JournalID = &jEntry.ID

		return s.repo.CreateTransaction(txCtx, txItem)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("execute cash in uow: %w", err)
	}

	return txItem, journalEntry, nil
}

// RecordCashOut records a Cash Disbursement (JKK) and automatically posts a balanced double-entry journal.
func (s *Service) RecordCashOut(ctx context.Context, cmd RecordCashOutCommand) (*cashbank.CashBankTransaction, *accountingDomain.JournalEntry, error) {
	if !cmd.Amount.IsPositive() {
		return nil, nil, ErrInvalidAmount
	}
	if cmd.TransactionDate.IsZero() {
		cmd.TransactionDate = time.Now()
	}

	cbAcc, err := s.repo.GetAccountByID(ctx, cmd.CashBankAccountID)
	if err != nil {
		return nil, nil, ErrAccountNotFound
	}

	counterAcc, err := s.accountingService.GetAccountByID(ctx, cmd.CounterpartAccountID)
	if err != nil {
		return nil, nil, ErrCounterpartNotFound
	}

	txNumber := fmt.Sprintf("BKK-%s-%04d", cmd.TransactionDate.Format("20060102"), time.Now().UnixNano()%10000)

	// Double-entry: Debit Expense/Liability Account, Credit Cash/Bank Account
	journalLines := []accountingApp.CreateLineCommand{
		{
			AccountID:   counterAcc.ID,
			Debit:       cmd.Amount,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Beban/Pengeluaran: %s", cmd.Description),
		},
		{
			AccountID:   cbAcc.AccountID,
			Debit:       decimal.Zero,
			Credit:      cmd.Amount,
			Description: fmt.Sprintf("Pengeluaran Kas/Bank: %s", cbAcc.Name),
		},
	}

	journalType := accountingDomain.JournalTypeCashDisbursement
	if cbAcc.Type == cashbank.AccountTypeBank {
		journalType = accountingDomain.JournalTypeBank
	}

	var journalEntry *accountingDomain.JournalEntry
	txItem := &cashbank.CashBankTransaction{
		BaseEntity:           shared.BaseEntity{ID: uuid.New()},
		BranchScoped:         shared.BranchScoped{BranchID: cmd.BranchID},
		CashBankAccountID:    cbAcc.ID,
		TransactionNumber:    txNumber,
		TransactionType:      cashbank.TransactionTypeCashOut,
		Amount:               cmd.Amount,
		TransactionDate:      cmd.TransactionDate,
		CounterpartAccountID: counterAcc.ID,
		Description:          cmd.Description,
		ReferenceNo:          cmd.ReferenceNo,
	}

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        cmd.BranchID,
			JournalType:     journalType,
			TransactionDate: cmd.TransactionDate,
			ReferenceType:   "CASH_OUT",
			ReferenceID:     txNumber,
			Description:     cmd.Description,
			CreatedBy:       cmd.CreatedBy,
			AutoPost:        true,
			Lines:           journalLines,
		})
		if err != nil {
			return fmt.Errorf("auto-post cash out journal: %w", err)
		}
		journalEntry = jEntry
		txItem.JournalID = &jEntry.ID

		return s.repo.CreateTransaction(txCtx, txItem)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("execute cash out uow: %w", err)
	}

	return txItem, journalEntry, nil
}

// TransferCashBank moves funds between cash/bank accounts (e.g. Kasir -> Bank BCA) without affecting revenues or expenses.
func (s *Service) TransferCashBank(ctx context.Context, cmd TransferCommand) (*cashbank.CashBankTransaction, *accountingDomain.JournalEntry, error) {
	if cmd.SourceAccountID == cmd.TargetAccountID {
		return nil, nil, ErrSameAccountTransfer
	}
	if !cmd.Amount.IsPositive() {
		return nil, nil, ErrInvalidAmount
	}
	if cmd.TransactionDate.IsZero() {
		cmd.TransactionDate = time.Now()
	}

	srcAcc, err := s.repo.GetAccountByID(ctx, cmd.SourceAccountID)
	if err != nil {
		return nil, nil, fmt.Errorf("source account: %w", ErrAccountNotFound)
	}

	targetAcc, err := s.repo.GetAccountByID(ctx, cmd.TargetAccountID)
	if err != nil {
		return nil, nil, fmt.Errorf("target account: %w", ErrAccountNotFound)
	}

	txNumber := fmt.Sprintf("TRF-%s-%04d", cmd.TransactionDate.Format("20060102"), time.Now().UnixNano()%10000)

	// Double-entry: Debit Destination Account, Credit Source Account
	journalLines := []accountingApp.CreateLineCommand{
		{
			AccountID:   targetAcc.AccountID,
			Debit:       cmd.Amount,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Penerimaan transfer dari %s", srcAcc.Name),
		},
		{
			AccountID:   srcAcc.AccountID,
			Debit:       decimal.Zero,
			Credit:      cmd.Amount,
			Description: fmt.Sprintf("Pengiriman transfer ke %s", targetAcc.Name),
		},
	}

	txItem := &cashbank.CashBankTransaction{
		BaseEntity:           shared.BaseEntity{ID: uuid.New()},
		BranchScoped:         shared.BranchScoped{BranchID: cmd.BranchID},
		CashBankAccountID:    srcAcc.ID,
		TransactionNumber:    txNumber,
		TransactionType:      cashbank.TransactionTypeTransfer,
		Amount:               cmd.Amount,
		TransactionDate:      cmd.TransactionDate,
		CounterpartAccountID: targetAcc.AccountID,
		Description:          fmt.Sprintf("Transfer dari %s ke %s: %s", srcAcc.Name, targetAcc.Name, cmd.Description),
		ReferenceNo:          cmd.ReferenceNo,
	}

	var journalEntry *accountingDomain.JournalEntry
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        cmd.BranchID,
			JournalType:     accountingDomain.JournalTypeGeneral,
			TransactionDate: cmd.TransactionDate,
			ReferenceType:   "CASH_TRANSFER",
			ReferenceID:     txNumber,
			Description:     txItem.Description,
			CreatedBy:       cmd.CreatedBy,
			AutoPost:        true,
			Lines:           journalLines,
		})
		if err != nil {
			return fmt.Errorf("auto-post transfer journal: %w", err)
		}
		journalEntry = jEntry
		txItem.JournalID = &jEntry.ID

		return s.repo.CreateTransaction(txCtx, txItem)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("execute transfer uow: %w", err)
	}

	return txItem, journalEntry, nil
}

func (s *Service) CreateAccount(ctx context.Context, acc *cashbank.CashBankAccount) error {
	return s.repo.CreateAccount(ctx, acc)
}

func (s *Service) GetAccountByID(ctx context.Context, id uuid.UUID) (*cashbank.CashBankAccount, error) {
	return s.repo.GetAccountByID(ctx, id)
}

func (s *Service) ListAccounts(ctx context.Context, branchID *uuid.UUID, accType *cashbank.AccountType) ([]cashbank.CashBankAccount, error) {
	return s.repo.ListAccounts(ctx, branchID, accType)
}

func (s *Service) ListTransactions(ctx context.Context, filter cashbank.TransactionFilter) ([]cashbank.CashBankTransaction, int64, error) {
	return s.repo.ListTransactions(ctx, filter)
}

