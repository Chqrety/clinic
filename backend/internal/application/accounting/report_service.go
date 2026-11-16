package accounting

import (
	"context"
	"fmt"
	"time"

	domain "clinic-api/internal/domain/accounting"

	"github.com/google/uuid"
)

type ReportRepository interface {
	GetGeneralLedger(ctx context.Context, filter domain.GeneralLedgerFilter) (*domain.AccountGeneralLedger, error)
	GetTrialBalance(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.TrialBalanceReport, error)
	GetIncomeStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.IncomeStatement, error)
	GetBalanceSheet(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*domain.BalanceSheet, error)
	GetCashFlowStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.CashFlowStatement, error)
}

type ReportService struct {
	reportRepo ReportRepository
}

func NewReportService(reportRepo ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

func (s *ReportService) GetGeneralLedger(ctx context.Context, filter domain.GeneralLedgerFilter) (*domain.AccountGeneralLedger, error) {
	if filter.StartDate.IsZero() {
		filter.StartDate = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	if filter.EndDate.IsZero() {
		filter.EndDate = time.Now()
	}
	if filter.StartDate.After(filter.EndDate) {
		return nil, fmt.Errorf("start_date cannot be after end_date")
	}
	return s.reportRepo.GetGeneralLedger(ctx, filter)
}

func (s *ReportService) GetTrialBalance(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.TrialBalanceReport, error) {
	if startDate.IsZero() {
		startDate = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start_date cannot be after end_date")
	}
	return s.reportRepo.GetTrialBalance(ctx, branchID, startDate, endDate)
}

func (s *ReportService) GetIncomeStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.IncomeStatement, error) {
	if startDate.IsZero() {
		startDate = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start_date cannot be after end_date")
	}
	return s.reportRepo.GetIncomeStatement(ctx, branchID, startDate, endDate)
}

func (s *ReportService) GetBalanceSheet(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*domain.BalanceSheet, error) {
	if asOfDate.IsZero() {
		asOfDate = time.Now()
	}
	return s.reportRepo.GetBalanceSheet(ctx, branchID, asOfDate)
}

func (s *ReportService) GetCashFlowStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.CashFlowStatement, error) {
	if startDate.IsZero() {
		startDate = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start_date cannot be after end_date")
	}
	return s.reportRepo.GetCashFlowStatement(ctx, branchID, startDate, endDate)
}

