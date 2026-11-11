package accounting_test

import (
	"context"
	"testing"
	"time"

	app "clinic-api/internal/application/accounting"
	domain "clinic-api/internal/domain/accounting"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockReportRepo struct {
	ledgerReport  *domain.AccountGeneralLedger
	tbReport      *domain.TrialBalanceReport
	isReport      *domain.IncomeStatement
	bsReport      *domain.BalanceSheet
	cfReport      *domain.CashFlowStatement
}

func (m *mockReportRepo) GetGeneralLedger(ctx context.Context, filter domain.GeneralLedgerFilter) (*domain.AccountGeneralLedger, error) {
	return m.ledgerReport, nil
}

func (m *mockReportRepo) GetTrialBalance(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.TrialBalanceReport, error) {
	return m.tbReport, nil
}

func (m *mockReportRepo) GetIncomeStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.IncomeStatement, error) {
	return m.isReport, nil
}

func (m *mockReportRepo) GetBalanceSheet(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*domain.BalanceSheet, error) {
	return m.bsReport, nil
}

func (m *mockReportRepo) GetCashFlowStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*domain.CashFlowStatement, error) {
	return m.cfReport, nil
}

func TestReportService_GeneralLedger(t *testing.T) {
	now := time.Now()
	mockRepo := &mockReportRepo{
		ledgerReport: &domain.AccountGeneralLedger{
			AccountCode:    "1-1111",
			AccountName:    "Kas Kasir",
			NormalBalance:  domain.NormalBalanceDebit,
			OpeningBalance: decimal.FromInt(500000),
			TotalDebit:     decimal.FromInt(250000),
			TotalCredit:    decimal.FromInt(100000),
			ClosingBalance: decimal.FromInt(650000), // 500k + 250k - 100k
			Transactions: []domain.LedgerTransactionLine{
				{
					Date:           now,
					JournalNumber:  "JKM-001",
					Debit:          decimal.FromInt(250000),
					Credit:         decimal.Zero,
					RunningBalance: decimal.FromInt(750000),
				},
				{
					Date:           now,
					JournalNumber:  "JKK-001",
					Debit:          decimal.Zero,
					Credit:         decimal.FromInt(100000),
					RunningBalance: decimal.FromInt(650000),
				},
			},
		},
	}

	svc := app.NewReportService(mockRepo)
	res, err := svc.GetGeneralLedger(context.Background(), domain.GeneralLedgerFilter{
		AccountCode: "1-1111",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.ClosingBalance.Equal(decimal.FromInt(650000)) {
		t.Errorf("expected closing balance 650000, got %s", res.ClosingBalance)
	}
	if len(res.Transactions) != 2 {
		t.Errorf("expected 2 transaction lines, got %d", len(res.Transactions))
	}
}

func TestReportService_TrialBalance_Balanced(t *testing.T) {
	mockRepo := &mockReportRepo{
		tbReport: &domain.TrialBalanceReport{
			TotalEndingDebit:  decimal.FromInt(15000000),
			TotalEndingCredit: decimal.FromInt(15000000),
			IsBalanced:         true,
			Items: []domain.TrialBalanceItem{
				{AccountCode: "1-1111", AccountName: "Kas Kasir", EndingDebit: decimal.FromInt(5000000)},
				{AccountCode: "1-1310", AccountName: "Persediaan Obat", EndingDebit: decimal.FromInt(10000000)},
				{AccountCode: "3-1100", AccountName: "Modal Disetor", EndingCredit: decimal.FromInt(15000000)},
			},
		},
	}

	svc := app.NewReportService(mockRepo)
	tb, err := svc.GetTrialBalance(context.Background(), nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !tb.IsBalanced {
		t.Errorf("expected Trial Balance to be balanced")
	}
	if !tb.TotalEndingDebit.Equal(tb.TotalEndingCredit) {
		t.Errorf("debit and credit mismatch: %s != %s", tb.TotalEndingDebit, tb.TotalEndingCredit)
	}
}

func TestReportService_IncomeStatement_Calculations(t *testing.T) {
	mockRepo := &mockReportRepo{
		isReport: &domain.IncomeStatement{
			TotalRevenue:  decimal.FromInt(10000000),
			TotalCOGS:     decimal.FromInt(4000000),
			GrossProfit:   decimal.FromInt(6000000),  // 10M - 4M
			TotalExpenses: decimal.FromInt(2500000),
			NetIncome:     decimal.FromInt(3500000),  // 6M - 2.5M
		},
	}

	svc := app.NewReportService(mockRepo)
	is, err := svc.GetIncomeStatement(context.Background(), nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedGross := is.TotalRevenue.Sub(is.TotalCOGS)
	if !is.GrossProfit.Equal(expectedGross) {
		t.Errorf("expected gross profit %s, got %s", expectedGross, is.GrossProfit)
	}

	expectedNet := is.GrossProfit.Sub(is.TotalExpenses)
	if !is.NetIncome.Equal(expectedNet) {
		t.Errorf("expected net income %s, got %s", expectedNet, is.NetIncome)
	}
}

func TestReportService_BalanceSheet_Equation(t *testing.T) {
	// Equation: Assets = Liabilities + Equity
	mockRepo := &mockReportRepo{
		bsReport: &domain.BalanceSheet{
			TotalAssets:            decimal.FromInt(50000000),
			TotalLiabilities:       decimal.FromInt(20000000),
			TotalEquity:            decimal.FromInt(30000000),
			TotalLiabilitiesEquity: decimal.FromInt(50000000),
			IsBalanced:             true,
		},
	}

	svc := app.NewReportService(mockRepo)
	bs, err := svc.GetBalanceSheet(context.Background(), nil, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bs.IsBalanced {
		t.Errorf("expected Balance Sheet to be balanced")
	}
	if !bs.TotalAssets.Equal(bs.TotalLiabilitiesEquity) {
		t.Errorf("accounting equation violated: Assets %s != Liabilities+Equity %s", bs.TotalAssets, bs.TotalLiabilitiesEquity)
	}
}

