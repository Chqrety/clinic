package accounting

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domain "clinic-api/internal/domain/accounting"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrUnpostedDraftsExist   = errors.New("cannot close year: unposted draft journals still exist in the period")
	ErrYearAlreadyClosed     = errors.New("fiscal year has already been closed")
	ErrInvalidYear           = errors.New("invalid fiscal year")
)

type ClosingService struct {
	db                *sql.DB
	uow               UnitOfWork
	accountingService *Service
	periodRepo        domain.FiscalPeriodRepository
	accountRepo       domain.AccountRepository
}

func NewClosingService(
	db *sql.DB,
	uow UnitOfWork,
	accountingService *Service,
	periodRepo domain.FiscalPeriodRepository,
	accountRepo domain.AccountRepository,
) *ClosingService {
	return &ClosingService{
		db:                db,
		uow:               uow,
		accountingService: accountingService,
		periodRepo:        periodRepo,
		accountRepo:       accountRepo,
	}
}

type YearEndClosingResult struct {
	FiscalYear       int                    `json:"fiscal_year"`
	BranchID         uuid.UUID              `json:"branch_id"`
	TotalRevenue     decimal.Decimal        `json:"total_revenue"`
	TotalExpenses    decimal.Decimal        `json:"total_expenses"`
	NetIncome        decimal.Decimal        `json:"net_income"`
	RetainedEarnings decimal.Decimal        `json:"retained_earnings_transferred"`
	ClosingJournalID uuid.UUID              `json:"closing_journal_id"`
	JournalNumber    string                 `json:"journal_number"`
	ClosedAt         time.Time              `json:"closed_at"`
}

// ExecuteYearEndClosing zeroes out all nominal accounts (Revenue, COGS, Expense) and transfers net profit to Retained Earnings (3-1200).
func (s *ClosingService) ExecuteYearEndClosing(ctx context.Context, branchID uuid.UUID, year int, closedBy *uuid.UUID) (*YearEndClosingResult, error) {
	if year < 2000 || year > 2100 {
		return nil, ErrInvalidYear
	}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC)

	// 1. Check for unposted DRAFT journals in the fiscal year
	var draftCount int
	draftQuery := `
		SELECT COUNT(*)
		FROM journal_entries
		WHERE branch_id = $1 AND transaction_date >= $2 AND transaction_date <= $3
		  AND status = 'DRAFT'
	`
	if err := s.db.QueryRowContext(ctx, draftQuery, branchID, startDate, endDate).Scan(&draftCount); err != nil {
		return nil, fmt.Errorf("check draft journals: %w", err)
	}
	if draftCount > 0 {
		return nil, fmt.Errorf("%w (%d draft journals found)", ErrUnpostedDraftsExist, draftCount)
	}

	// 2. Query all nominal account balances (REVENUE, COGS, EXPENSE)
	nominalQuery := `
		SELECT a.code, a.name, a.type,
		       COALESCE(SUM(jl.debit), 0) as total_debit,
		       COALESCE(SUM(jl.credit), 0) as total_credit
		FROM accounts a
		JOIN journal_lines jl ON a.id = jl.account_id
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		WHERE je.status = 'POSTED'
		  AND je.journal_type != 'CLOSING'
		  AND je.transaction_date >= $1 AND je.transaction_date <= $2
		  AND je.branch_id = $3
		  AND a.type IN ('REVENUE', 'COGS', 'EXPENSE')
		GROUP BY a.code, a.name, a.type
		HAVING COALESCE(SUM(jl.debit), 0) != COALESCE(SUM(jl.credit), 0)
	`
	rows, err := s.db.QueryContext(ctx, nominalQuery, startDate, endDate, branchID)
	if err != nil {
		return nil, fmt.Errorf("query nominal accounts for closing: %w", err)
	}
	defer rows.Close()

	var closingLines []CreateLineCommand
	totalRevenue := decimal.Zero
	totalExpense := decimal.Zero

	for rows.Next() {
		var code, name string
		var accType domain.AccountType
		var drStr, crStr string
		if err := rows.Scan(&code, &name, &accType, &drStr, &crStr); err != nil {
			return nil, fmt.Errorf("scan nominal row: %w", err)
		}
		dr, _ := decimal.NewFromString(drStr)
		cr, _ := decimal.NewFromString(crStr)

		if accType == domain.AccountTypeRevenue {
			netRevenue := cr.Sub(dr)
			if netRevenue.IsPositive() {
				// To close revenue with credit balance: DEBIT revenue account
				closingLines = append(closingLines, CreateLineCommand{
					AccountCode: code,
					Debit:       netRevenue,
					Credit:      decimal.Zero,
					Description: fmt.Sprintf("Tutup Akun Pendapatan %s Tahun %d", name, year),
				})
				totalRevenue = totalRevenue.Add(netRevenue)
			}
		} else { // COGS or EXPENSE
			netExp := dr.Sub(cr)
			if netExp.IsPositive() {
				// To close expense with debit balance: CREDIT expense account
				closingLines = append(closingLines, CreateLineCommand{
					AccountCode: code,
					Debit:       decimal.Zero,
					Credit:      netExp,
					Description: fmt.Sprintf("Tutup Akun Beban %s Tahun %d", name, year),
				})
				totalExpense = totalExpense.Add(netExp)
			}
		}
	}

	netIncome := totalRevenue.Sub(totalExpense)

	// Balance figure goes to 3-1200 Laba Ditahan (Retained Earnings)
	retainedEarningsAcc := "3-1200"
	if netIncome.IsPositive() {
		// Profit: Credit Retained Earnings
		closingLines = append(closingLines, CreateLineCommand{
			AccountCode: retainedEarningsAcc,
			Debit:       decimal.Zero,
			Credit:      netIncome,
			Description: fmt.Sprintf("Transfer Laba Bersih Tahun %d ke Laba Ditahan", year),
		})
	} else if netIncome.IsNegative() {
		// Loss: Debit Retained Earnings
		loss := netIncome.Abs()
		closingLines = append(closingLines, CreateLineCommand{
			AccountCode: retainedEarningsAcc,
			Debit:       loss,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Transfer Rugi Bersih Tahun %d ke Laba Ditahan", year),
		})
	}

	if len(closingLines) == 0 {
		return nil, errors.New("no nominal transactions found to close for specified year")
	}

	// 3. Post Closing Journal (JT) in Atomic UoW and Lock Fiscal Year
	var closingJournal *domain.JournalEntry
	now := time.Now()

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		jEntry, err := s.accountingService.CreateJournal(txCtx, CreateJournalCommand{
			BranchID:        branchID,
			JournalType:     domain.JournalTypeClosing,
			TransactionDate: endDate,
			ReferenceType:   "YEAR_END_CLOSING",
			ReferenceID:     fmt.Sprintf("YEC-%d", year),
			Description:     fmt.Sprintf("Jurnal Penutup Akhir Tahun Buku %d", year),
			CreatedBy:       closedBy,
			AutoPost:        true,
			Lines:           closingLines,
		})
		if err != nil {
			return fmt.Errorf("create closing journal: %w", err)
		}
		closingJournal = jEntry

		// Lock all 12 monthly periods
		for m := 1; m <= 12; m++ {
			p, err := s.periodRepo.GetPeriod(txCtx, branchID, year, m)
			if err == nil && p != nil && p.Status == domain.PeriodStatusOpen {
				_ = s.periodRepo.ClosePeriod(txCtx, p)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("execute year-end closing uow: %w", err)
	}

	return &YearEndClosingResult{
		FiscalYear:       year,
		BranchID:         branchID,
		TotalRevenue:     totalRevenue,
		TotalExpenses:    totalExpense,
		NetIncome:        netIncome,
		RetainedEarnings: netIncome,
		ClosingJournalID: closingJournal.ID,
		JournalNumber:    closingJournal.JournalNumber,
		ClosedAt:         now,
	}, nil
}

