package reporting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"clinic-api/internal/domain/accounting"

	"github.com/google/uuid"
)

var (
	ErrInvalidDateRange = errors.New("start date cannot be after end date")
	ErrDateRequired     = errors.New("date parameters are required")
)

type Service struct {
	reportRepo accounting.ReportRepository
}

func NewService(reportRepo accounting.ReportRepository) *Service {
	return &Service{
		reportRepo: reportRepo,
	}
}

// GetGeneralLedger produces a detailed ledger with opening balance, chronological transaction lines, and running balances.
func (s *Service) GetGeneralLedger(ctx context.Context, filter accounting.GeneralLedgerFilter) (*accounting.AccountGeneralLedger, error) {
	if filter.StartDate.After(filter.EndDate) {
		return nil, ErrInvalidDateRange
	}
	if filter.AccountID == nil && filter.AccountCode == "" {
		return nil, fmt.Errorf("account_id or account_code must be provided")
	}

	return s.reportRepo.GetGeneralLedger(ctx, filter)
}

// GetTrialBalance produces the 6-column Neraca Saldo (Opening Dr/Cr, Period Dr/Cr, Ending Dr/Cr).
func (s *Service) GetTrialBalance(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*accounting.TrialBalanceReport, error) {
	if startDate.After(endDate) {
		return nil, ErrInvalidDateRange
	}

	return s.reportRepo.GetTrialBalance(ctx, branchID, startDate, endDate)
}

// GetIncomeStatement produces Laporan Laba Rugi (Pendapatan - HPP = Laba Kotor - Beban = Laba Bersih).
func (s *Service) GetIncomeStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*accounting.IncomeStatement, error) {
	if startDate.After(endDate) {
		return nil, ErrInvalidDateRange
	}

	return s.reportRepo.GetIncomeStatement(ctx, branchID, startDate, endDate)
}

// GetBalanceSheet produces Laporan Neraca (Aset = Liabilitas + Ekuitas) as of a specific date.
func (s *Service) GetBalanceSheet(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*accounting.BalanceSheet, error) {
	if asOfDate.IsZero() {
		asOfDate = time.Now()
	}

	return s.reportRepo.GetBalanceSheet(ctx, branchID, asOfDate)
}

// GetCashFlowStatement produces Laporan Arus Kas (Operating, Investing, Financing).
func (s *Service) GetCashFlowStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*accounting.CashFlowStatement, error) {
	if startDate.After(endDate) {
		return nil, ErrInvalidDateRange
	}

	return s.reportRepo.GetCashFlowStatement(ctx, branchID, startDate, endDate)
}


