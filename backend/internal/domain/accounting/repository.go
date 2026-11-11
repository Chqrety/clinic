package accounting

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByCode(ctx context.Context, code string) (*Account, error)
	List(ctx context.Context, filter AccountFilter) ([]Account, int64, error)
	Create(ctx context.Context, acc *Account) error
	Update(ctx context.Context, acc *Account) error
}

type AccountFilter struct {
	Type      *AccountType
	IsActive  *bool
	IsHeader  *bool
	Search    string
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}

type JournalRepository interface {
	Create(ctx context.Context, entry *JournalEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*JournalEntry, error)
	GetByJournalNumber(ctx context.Context, number string) (*JournalEntry, error)
	GetByReference(ctx context.Context, refType, refID string) (*JournalEntry, error)
	List(ctx context.Context, filter JournalFilter) ([]JournalEntry, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status JournalStatus, postedBy *uuid.UUID) error
}

type JournalFilter struct {
	BranchID  *uuid.UUID
	Type      *JournalType
	Status    *JournalStatus
	StartDate *time.Time
	EndDate   *time.Time
	Search    string
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}

type FiscalPeriodRepository interface {
	GetPeriod(ctx context.Context, branchID uuid.UUID, year, month int) (*FiscalPeriod, error)
	ClosePeriod(ctx context.Context, period *FiscalPeriod) error
	ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error
	IsPeriodLocked(ctx context.Context, branchID uuid.UUID, date time.Time) (bool, error)
}

type ReportRepository interface {
	GetGeneralLedger(ctx context.Context, filter GeneralLedgerFilter) (*AccountGeneralLedger, error)
	GetTrialBalance(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*TrialBalanceReport, error)
	GetIncomeStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*IncomeStatement, error)
	GetBalanceSheet(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*BalanceSheet, error)
	GetCashFlowStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*CashFlowStatement, error)
}

