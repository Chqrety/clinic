package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"clinic-api/internal/domain/accounting"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresReportRepo struct {
	db *sql.DB
}

func NewPostgresReportRepo(db *sql.DB) *PostgresReportRepo {
	return &PostgresReportRepo{db: db}
}

// GetGeneralLedger computes running balance and transaction lines for a specific account.
func (r *PostgresReportRepo) GetGeneralLedger(ctx context.Context, filter accounting.GeneralLedgerFilter) (*accounting.AccountGeneralLedger, error) {
	exec := GetExecutor(ctx, r.db)

	// 1. Get Account Details
	var acc accounting.Account
	var queryAcc string
	var arg any
	if filter.AccountID != nil && *filter.AccountID != uuid.Nil {
		queryAcc = `SELECT id, code, name, type, normal_balance FROM accounts WHERE id = $1`
		arg = *filter.AccountID
	} else if filter.AccountCode != "" {
		queryAcc = `SELECT id, code, name, type, normal_balance FROM accounts WHERE code = $1`
		arg = filter.AccountCode
	} else {
		return nil, fmt.Errorf("account_id or account_code must be specified for general ledger")
	}

	err := exec.QueryRowContext(ctx, queryAcc, arg).Scan(&acc.ID, &acc.Code, &acc.Name, &acc.Type, &acc.NormalBalance)
	if err != nil {
		return nil, fmt.Errorf("lookup account for ledger: %w", err)
	}

	// 2. Compute Opening Balance (before StartDate)
	queryOpening := `
		SELECT COALESCE(SUM(jl.debit), 0), COALESCE(SUM(jl.credit), 0)
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		WHERE jl.account_id = $1 AND je.status = 'POSTED' AND je.transaction_date < $2
		  AND ($3::uuid IS NULL OR je.branch_id = $3)
	`
	var opDebit, opCredit decimal.Decimal
	err = exec.QueryRowContext(ctx, queryOpening, acc.ID, filter.StartDate, filter.BranchID).Scan(&opDebit, &opCredit)
	if err != nil {
		return nil, fmt.Errorf("query opening balance: %w", err)
	}

	openingBalance := decimal.Zero
	if acc.NormalBalance == accounting.NormalBalanceDebit {
		openingBalance = opDebit.Sub(opCredit)
	} else {
		openingBalance = opCredit.Sub(opDebit)
	}

	// 3. Query Transactions between StartDate and EndDate
	queryTx := `
		SELECT je.transaction_date, je.id, je.journal_number, je.journal_type,
		       je.reference_type, je.reference_id, COALESCE(jl.description, je.description),
		       jl.debit, jl.credit
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		WHERE jl.account_id = $1 AND je.status = 'POSTED'
		  AND je.transaction_date >= $2 AND je.transaction_date <= $3
		  AND ($4::uuid IS NULL OR je.branch_id = $4)
		ORDER BY je.transaction_date ASC, je.created_at ASC, jl.line_order ASC
	`
	rows, err := exec.QueryContext(ctx, queryTx, acc.ID, filter.StartDate, filter.EndDate, filter.BranchID)
	if err != nil {
		return nil, fmt.Errorf("query ledger transactions: %w", err)
	}
	defer rows.Close()

	running := openingBalance
	totalDebit := decimal.Zero
	totalCredit := decimal.Zero
	var txLines []accounting.LedgerTransactionLine

	for rows.Next() {
		var line accounting.LedgerTransactionLine
		var debit, credit decimal.Decimal

		if err := rows.Scan(
			&line.Date, &line.JournalID, &line.JournalNumber, &line.JournalType,
			&line.ReferenceType, &line.ReferenceID, &line.Description,
			&debit, &credit,
		); err != nil {
			return nil, fmt.Errorf("scan ledger row: %w", err)
		}

		line.Debit = debit
		line.Credit = credit
		totalDebit = totalDebit.Add(debit)
		totalCredit = totalCredit.Add(credit)

		if acc.NormalBalance == accounting.NormalBalanceDebit {
			running = running.Add(debit).Sub(credit)
		} else {
			running = running.Add(credit).Sub(debit)
		}
		line.RunningBalance = running

		txLines = append(txLines, line)
	}

	return &accounting.AccountGeneralLedger{
		AccountID:      acc.ID,
		AccountCode:    acc.Code,
		AccountName:    acc.Name,
		AccountType:    acc.Type,
		NormalBalance:  acc.NormalBalance,
		OpeningBalance: openingBalance,
		TotalDebit:     totalDebit,
		TotalCredit:    totalCredit,
		ClosingBalance: running,
		Transactions:   txLines,
	}, nil
}

// GetTrialBalance aggregates opening, period, and ending balances across all active accounts.
func (r *PostgresReportRepo) GetTrialBalance(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*accounting.TrialBalanceReport, error) {
	exec := GetExecutor(ctx, r.db)

	query := `
		SELECT a.id, a.code, a.name, a.type, a.normal_balance,
		       COALESCE(SUM(CASE WHEN je.transaction_date < $1 THEN jl.debit ELSE 0 END), 0) as op_debit,
		       COALESCE(SUM(CASE WHEN je.transaction_date < $1 THEN jl.credit ELSE 0 END), 0) as op_credit,
		       COALESCE(SUM(CASE WHEN je.transaction_date >= $1 AND je.transaction_date <= $2 THEN jl.debit ELSE 0 END), 0) as period_debit,
		       COALESCE(SUM(CASE WHEN je.transaction_date >= $1 AND je.transaction_date <= $2 THEN jl.credit ELSE 0 END), 0) as period_credit
		FROM accounts a
		LEFT JOIN journal_lines jl ON a.id = jl.account_id
		LEFT JOIN journal_entries je ON jl.journal_entry_id = je.id AND je.status = 'POSTED' AND ($3::uuid IS NULL OR je.branch_id = $3)
		WHERE a.is_active = TRUE AND a.is_header = FALSE
		GROUP BY a.id, a.code, a.name, a.type, a.normal_balance
		ORDER BY a.code ASC
	`

	rows, err := exec.QueryContext(ctx, query, startDate, endDate, branchID)
	if err != nil {
		return nil, fmt.Errorf("query trial balance: %w", err)
	}
	defer rows.Close()

	var items []accounting.TrialBalanceItem
	totalOpDebit := decimal.Zero
	totalOpCredit := decimal.Zero
	totalPerDebit := decimal.Zero
	totalPerCredit := decimal.Zero
	totalEndDebit := decimal.Zero
	totalEndCredit := decimal.Zero

	for rows.Next() {
		var it accounting.TrialBalanceItem
		var rawOpDebit, rawOpCredit, perDebit, perCredit decimal.Decimal

		if err := rows.Scan(
			&it.AccountID, &it.AccountCode, &it.AccountName, &it.AccountType, &it.NormalBalance,
			&rawOpDebit, &rawOpCredit, &perDebit, &perCredit,
		); err != nil {
			return nil, fmt.Errorf("scan trial balance row: %w", err)
		}

		it.PeriodDebit = perDebit
		it.PeriodCredit = perCredit

		// Compute Net Opening
		if it.NormalBalance == accounting.NormalBalanceDebit {
			netOp := rawOpDebit.Sub(rawOpCredit)
			if netOp.GreaterThanOrEqual(decimal.Zero) {
				it.OpeningDebit = netOp
				it.OpeningCredit = decimal.Zero
			} else {
				it.OpeningDebit = decimal.Zero
				it.OpeningCredit = netOp.Abs()
			}
			// Ending = (OpeningDebit - OpeningCredit) + (PeriodDebit - PeriodCredit)
			netEnd := it.OpeningDebit.Sub(it.OpeningCredit).Add(perDebit).Sub(perCredit)
			if netEnd.GreaterThanOrEqual(decimal.Zero) {
				it.EndingDebit = netEnd
				it.EndingCredit = decimal.Zero
			} else {
				it.EndingDebit = decimal.Zero
				it.EndingCredit = netEnd.Abs()
			}
		} else {
			netOp := rawOpCredit.Sub(rawOpDebit)
			if netOp.GreaterThanOrEqual(decimal.Zero) {
				it.OpeningCredit = netOp
				it.OpeningDebit = decimal.Zero
			} else {
				it.OpeningCredit = decimal.Zero
				it.OpeningDebit = netOp.Abs()
			}
			// Ending = (OpeningCredit - OpeningDebit) + (PeriodCredit - PeriodDebit)
			netEnd := it.OpeningCredit.Sub(it.OpeningDebit).Add(perCredit).Sub(perDebit)
			if netEnd.GreaterThanOrEqual(decimal.Zero) {
				it.EndingCredit = netEnd
				it.EndingDebit = decimal.Zero
			} else {
				it.EndingCredit = decimal.Zero
				it.EndingDebit = netEnd.Abs()
			}
		}

		// Accumulate totals
		totalOpDebit = totalOpDebit.Add(it.OpeningDebit)
		totalOpCredit = totalOpCredit.Add(it.OpeningCredit)
		totalPerDebit = totalPerDebit.Add(it.PeriodDebit)
		totalPerCredit = totalPerCredit.Add(it.PeriodCredit)
		totalEndDebit = totalEndDebit.Add(it.EndingDebit)
		totalEndCredit = totalEndCredit.Add(it.EndingCredit)

		items = append(items, it)
	}

	isBalanced := totalEndDebit.Equal(totalEndCredit)

	return &accounting.TrialBalanceReport{
		BranchID:           branchID,
		StartDate:          startDate,
		EndDate:            endDate,
		TotalOpeningDebit:  totalOpDebit,
		TotalOpeningCredit: totalOpCredit,
		TotalPeriodDebit:   totalPerDebit,
		TotalPeriodCredit:  totalPerCredit,
		TotalEndingDebit:   totalEndDebit,
		TotalEndingCredit:  totalEndCredit,
		IsBalanced:         isBalanced,
		Items:              items,
	}, nil
}

// GetIncomeStatement computes Laba Rugi for a given date range.
func (r *PostgresReportRepo) GetIncomeStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*accounting.IncomeStatement, error) {
	exec := GetExecutor(ctx, r.db)

	query := `
		SELECT a.code, a.name, a.type,
		       COALESCE(SUM(jl.debit), 0) as total_debit,
		       COALESCE(SUM(jl.credit), 0) as total_credit
		FROM accounts a
		JOIN journal_lines jl ON a.id = jl.account_id
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		WHERE je.status = 'POSTED'
		  AND je.transaction_date >= $1 AND je.transaction_date <= $2
		  AND ($3::uuid IS NULL OR je.branch_id = $3)
		  AND a.type IN ('REVENUE', 'COGS', 'EXPENSE')
		GROUP BY a.code, a.name, a.type
		ORDER BY a.code ASC
	`
	rows, err := exec.QueryContext(ctx, query, startDate, endDate, branchID)
	if err != nil {
		return nil, fmt.Errorf("query income statement data: %w", err)
	}
	defer rows.Close()

	var revItems, cogsItems, expItems []accounting.FinancialReportItem
	totalRev := decimal.Zero
	totalCOGS := decimal.Zero
	totalExp := decimal.Zero

	for rows.Next() {
		var code, name, accType string
		var debit, credit decimal.Decimal

		if err := rows.Scan(&code, &name, &accType, &debit, &credit); err != nil {
			return nil, fmt.Errorf("scan income statement row: %w", err)
		}

		switch accType {
		case "REVENUE":
			// Revenue normal balance is CREDIT: amount = credit - debit
			amount := credit.Sub(debit)
			if !amount.IsZero() {
				revItems = append(revItems, accounting.FinancialReportItem{AccountCode: code, AccountName: name, Amount: amount})
				totalRev = totalRev.Add(amount)
			}
		case "COGS":
			// COGS normal balance is DEBIT: amount = debit - credit
			amount := debit.Sub(credit)
			if !amount.IsZero() {
				cogsItems = append(cogsItems, accounting.FinancialReportItem{AccountCode: code, AccountName: name, Amount: amount})
				totalCOGS = totalCOGS.Add(amount)
			}
		case "EXPENSE":
			// Expense normal balance is DEBIT: amount = debit - credit
			amount := debit.Sub(credit)
			if !amount.IsZero() {
				expItems = append(expItems, accounting.FinancialReportItem{AccountCode: code, AccountName: name, Amount: amount})
				totalExp = totalExp.Add(amount)
			}
		}
	}

	grossProfit := totalRev.Sub(totalCOGS)
	netIncome := grossProfit.Sub(totalExp)

	return &accounting.IncomeStatement{
		BranchID:  branchID,
		StartDate: startDate,
		EndDate:   endDate,
		RevenueCategories: []accounting.FinancialReportCategory{
			{CategoryName: "Pendapatan Operasional Klinik", Total: totalRev, Items: revItems},
		},
		TotalRevenue: totalRev,
		COGSCategories: []accounting.FinancialReportCategory{
			{CategoryName: "Harga Pokok Penjualan (HPP)", Total: totalCOGS, Items: cogsItems},
		},
		TotalCOGS:   totalCOGS,
		GrossProfit: grossProfit,
		ExpenseCategories: []accounting.FinancialReportCategory{
			{CategoryName: "Beban Operasional Klinik", Total: totalExp, Items: expItems},
		},
		TotalExpenses: totalExp,
		NetIncome:     netIncome,
	}, nil
}

// GetBalanceSheet computes Neraca as of a specific date.
func (r *PostgresReportRepo) GetBalanceSheet(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*accounting.BalanceSheet, error) {
	exec := GetExecutor(ctx, r.db)

	// 1. Query Assets, Liabilities, and Equity balances
	query := `
		SELECT a.code, a.name, a.type, a.normal_balance,
		       COALESCE(SUM(jl.debit), 0) as total_debit,
		       COALESCE(SUM(jl.credit), 0) as total_credit
		FROM accounts a
		LEFT JOIN journal_lines jl ON a.id = jl.account_id
		LEFT JOIN journal_entries je ON jl.journal_entry_id = je.id AND je.status = 'POSTED' AND je.transaction_date <= $1 AND ($2::uuid IS NULL OR je.branch_id = $2)
		WHERE a.is_active = TRUE AND a.is_header = FALSE AND a.type IN ('ASSET', 'LIABILITY', 'EQUITY')
		GROUP BY a.code, a.name, a.type, a.normal_balance
		ORDER BY a.code ASC
	`
	rows, err := exec.QueryContext(ctx, query, asOfDate, branchID)
	if err != nil {
		return nil, fmt.Errorf("query balance sheet items: %w", err)
	}
	defer rows.Close()

	var curAssets, nonCurAssets, curLiab, longLiab, eqItems []accounting.FinancialReportItem
	totalCurAsset := decimal.Zero
	totalNonCurAsset := decimal.Zero
	totalCurLiab := decimal.Zero
	totalLongLiab := decimal.Zero
	totalEquity := decimal.Zero

	for rows.Next() {
		var code, name, accType, normalBalance string
		var debit, credit decimal.Decimal

		if err := rows.Scan(&code, &name, &accType, &normalBalance, &debit, &credit); err != nil {
			return nil, fmt.Errorf("scan balance sheet row: %w", err)
		}

		switch accType {
		case "ASSET":
			amount := debit.Sub(credit)
			if strings.HasPrefix(code, "1-1") { // Aset Lancar
				curAssets = append(curAssets, accounting.FinancialReportItem{AccountCode: code, AccountName: name, Amount: amount})
				totalCurAsset = totalCurAsset.Add(amount)
			} else { // Aset Tidak Lancar / Tetap
				nonCurAssets = append(nonCurAssets, accounting.FinancialReportItem{AccountCode: code, AccountName: name, Amount: amount})
				totalNonCurAsset = totalNonCurAsset.Add(amount)
			}
		case "LIABILITY":
			amount := credit.Sub(debit)
			if strings.HasPrefix(code, "2-1") { // Liabilitas Jangka Pendek
				curLiab = append(curLiab, accounting.FinancialReportItem{AccountCode: code, AccountName: name, Amount: amount})
				totalCurLiab = totalCurLiab.Add(amount)
			} else { // Liabilitas Jangka Panjang
				longLiab = append(longLiab, accounting.FinancialReportItem{AccountCode: code, AccountName: name, Amount: amount})
				totalLongLiab = totalLongLiab.Add(amount)
			}
		case "EQUITY":
			amount := credit.Sub(debit)
			eqItems = append(eqItems, accounting.FinancialReportItem{AccountCode: code, AccountName: name, Amount: amount})
			totalEquity = totalEquity.Add(amount)
		}
	}

	// 2. Compute Current Period Profit (cumulative net income from inception up to AsOfDate)
	queryProfit := `
		SELECT 
			COALESCE(SUM(CASE WHEN a.type = 'REVENUE' THEN (jl.credit - jl.debit) ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN a.type IN ('COGS', 'EXPENSE') THEN (jl.debit - jl.credit) ELSE 0 END), 0)
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN accounts a ON jl.account_id = a.id
		WHERE je.status = 'POSTED' AND je.transaction_date <= $1
		  AND ($2::uuid IS NULL OR je.branch_id = $2)
	`
	var currentPeriodProfit decimal.Decimal
	_ = exec.QueryRowContext(ctx, queryProfit, asOfDate, branchID).Scan(&currentPeriodProfit)

	totalAssets := totalCurAsset.Add(totalNonCurAsset)
	totalLiabilities := totalCurLiab.Add(totalLongLiab)
	grandTotalEquity := totalEquity.Add(currentPeriodProfit)
	totalLiabEquity := totalLiabilities.Add(grandTotalEquity)

	isBalanced := totalAssets.Equal(totalLiabEquity)

	return &accounting.BalanceSheet{
		BranchID: branchID,
		AsOfDate: asOfDate,
		CurrentAssets: accounting.FinancialReportCategory{
			CategoryName: "Aset Lancar", Total: totalCurAsset, Items: curAssets,
		},
		NonCurrentAssets: accounting.FinancialReportCategory{
			CategoryName: "Aset Tidak Lancar & Aset Tetap", Total: totalNonCurAsset, Items: nonCurAssets,
		},
		TotalAssets: totalAssets,
		CurrentLiabilities: accounting.FinancialReportCategory{
			CategoryName: "Liabilitas Jangka Pendek", Total: totalCurLiab, Items: curLiab,
		},
		LongTermLiabilities: accounting.FinancialReportCategory{
			CategoryName: "Liabilitas Jangka Panjang", Total: totalLongLiab, Items: longLiab,
		},
		TotalLiabilities: totalLiabilities,
		EquityCategories: accounting.FinancialReportCategory{
			CategoryName: "Ekuitas Pemilik & Laba Ditahan", Total: totalEquity, Items: eqItems,
		},
		CurrentPeriodProfit:    currentPeriodProfit,
		TotalEquity:            grandTotalEquity,
		TotalLiabilitiesEquity: totalLiabEquity,
		IsBalanced:             isBalanced,
	}, nil
}

// GetCashFlowStatement computes cash flows from operating, investing, and financing activities.
func (r *PostgresReportRepo) GetCashFlowStatement(ctx context.Context, branchID *uuid.UUID, startDate, endDate time.Time) (*accounting.CashFlowStatement, error) {
	exec := GetExecutor(ctx, r.db)

	// 1. Beginning Cash & Setara Kas (Accounts 1-11xx before startDate)
	queryBeginning := `
		SELECT COALESCE(SUM(jl.debit - jl.credit), 0)
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN accounts a ON jl.account_id = a.id
		WHERE je.status = 'POSTED' AND je.transaction_date < $1
		  AND a.code LIKE '1-11%'
		  AND ($2::uuid IS NULL OR je.branch_id = $2)
	`
	var beginningCash decimal.Decimal
	if err := exec.QueryRowContext(ctx, queryBeginning, startDate, branchID).Scan(&beginningCash); err != nil {
		return nil, fmt.Errorf("query beginning cash: %w", err)
	}

	// 2. Ending Cash & Setara Kas (Accounts 1-11xx up to endDate)
	queryEnding := `
		SELECT COALESCE(SUM(jl.debit - jl.credit), 0)
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN accounts a ON jl.account_id = a.id
		WHERE je.status = 'POSTED' AND je.transaction_date <= $1
		  AND a.code LIKE '1-11%'
		  AND ($2::uuid IS NULL OR je.branch_id = $2)
	`
	var endingCash decimal.Decimal
	if err := exec.QueryRowContext(ctx, queryEnding, endDate, branchID).Scan(&endingCash); err != nil {
		return nil, fmt.Errorf("query ending cash: %w", err)
	}

	netCashChange := endingCash.Sub(beginningCash)

	// 3. Operating Activities: Inflows (Revenues & AR collections into Cash/Bank)
	queryOperatingIn := `
		SELECT COALESCE(SUM(jl.debit), 0)
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN accounts a ON jl.account_id = a.id
		WHERE je.status = 'POSTED' AND je.transaction_date >= $1 AND je.transaction_date <= $2
		  AND a.code LIKE '1-11%'
		  AND je.journal_type IN ('CASH_RECEIPT', 'BANK')
		  AND ($3::uuid IS NULL OR je.branch_id = $3)
	`
	var opIn decimal.Decimal
	_ = exec.QueryRowContext(ctx, queryOperatingIn, startDate, endDate, branchID).Scan(&opIn)

	// Operating Outflows (Disbursements for suppliers, doctor shares, operating expenses)
	queryOperatingOut := `
		SELECT COALESCE(SUM(jl.credit), 0)
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN accounts a ON jl.account_id = a.id
		WHERE je.status = 'POSTED' AND je.transaction_date >= $1 AND je.transaction_date <= $2
		  AND a.code LIKE '1-11%'
		  AND je.journal_type IN ('CASH_DISBURSEMENT', 'BANK')
		  AND ($3::uuid IS NULL OR je.branch_id = $3)
	`
	var opOut decimal.Decimal
	_ = exec.QueryRowContext(ctx, queryOperatingOut, startDate, endDate, branchID).Scan(&opOut)

	netOperating := opIn.Sub(opOut)

	opItems := []accounting.FinancialReportItem{
		{AccountCode: "INFLOW-OP", AccountName: "Penerimaan Kas Layanan & Piutang Pasien", Amount: opIn},
		{AccountCode: "OUTFLOW-OP", AccountName: "Pengeluaran Kas Operasional, Obat & Dokter", Amount: opOut.Neg()},
	}

	// 4. Investing Activities (Purchase/Disposal of Fixed Assets 1-2xxx)
	queryInvesting := `
		SELECT COALESCE(SUM(jl.debit - jl.credit), 0)
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN accounts a ON jl.account_id = a.id
		WHERE je.status = 'POSTED' AND je.transaction_date >= $1 AND je.transaction_date <= $2
		  AND a.code LIKE '1-2%' AND a.code NOT LIKE '1-29%'
		  AND ($3::uuid IS NULL OR je.branch_id = $3)
	`
	var invCapex decimal.Decimal
	_ = exec.QueryRowContext(ctx, queryInvesting, startDate, endDate, branchID).Scan(&invCapex)
	netInvesting := invCapex.Neg() // Capex is cash outflow

	invItems := []accounting.FinancialReportItem{
		{AccountCode: "CAPEX-ASSET", AccountName: "Perolehan Mesin Medis & Aset Tetap", Amount: netInvesting},
	}

	// 5. Financing Activities (Equity contributions & long-term loans)
	queryFinancing := `
		SELECT COALESCE(SUM(jl.credit - jl.debit), 0)
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN accounts a ON jl.account_id = a.id
		WHERE je.status = 'POSTED' AND je.transaction_date >= $1 AND je.transaction_date <= $2
		  AND (a.code LIKE '3-11%' OR a.code LIKE '2-2%')
		  AND ($3::uuid IS NULL OR je.branch_id = $3)
	`
	var finFlow decimal.Decimal
	_ = exec.QueryRowContext(ctx, queryFinancing, startDate, endDate, branchID).Scan(&finFlow)
	netFinancing := finFlow

	finItems := []accounting.FinancialReportItem{
		{AccountCode: "FIN-EQUITY", AccountName: "Setoran Modal & Pembiayaan", Amount: netFinancing},
	}

	return &accounting.CashFlowStatement{
		BranchID:  branchID,
		StartDate: startDate,
		EndDate:   endDate,
		OperatingActivities: accounting.FinancialReportCategory{
			CategoryName: "Arus Kas dari Aktivitas Operasi",
			Total:        netOperating,
			Items:        opItems,
		},
		NetOperatingCash: netOperating,
		InvestingActivities: accounting.FinancialReportCategory{
			CategoryName: "Arus Kas dari Aktivitas Investasi",
			Total:        netInvesting,
			Items:        invItems,
		},
		NetInvestingCash: netInvesting,
		FinancingActivities: accounting.FinancialReportCategory{
			CategoryName: "Arus Kas dari Aktivitas Pendanaan",
			Total:        netFinancing,
			Items:        finItems,
		},
		NetFinancingCash: netFinancing,
		BeginningCash:    beginningCash,
		NetCashChange:    netCashChange,
		EndingCash:       endingCash,
	}, nil
}
