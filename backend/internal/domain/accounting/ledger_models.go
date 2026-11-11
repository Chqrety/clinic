package accounting

import (
	"time"

	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

// LedgerTransactionLine represents a single posting entry in the General Ledger.
type LedgerTransactionLine struct {
	Date          time.Time       `json:"date"`
	JournalID     uuid.UUID       `json:"journal_id"`
	JournalNumber string          `json:"journal_number"`
	JournalType   JournalType     `json:"journal_type"`
	ReferenceType string          `json:"reference_type"`
	ReferenceID   string          `json:"reference_id"`
	Description   string          `json:"description"`
	Debit         decimal.Decimal `json:"debit"`
	Credit        decimal.Decimal `json:"credit"`
	RunningBalance decimal.Decimal `json:"running_balance"`
}

// AccountGeneralLedger represents the General Ledger statement for a specific account.
type AccountGeneralLedger struct {
	AccountID      uuid.UUID               `json:"account_id"`
	AccountCode    string                  `json:"account_code"`
	AccountName    string                  `json:"account_name"`
	AccountType    AccountType             `json:"account_type"`
	NormalBalance  NormalBalance           `json:"normal_balance"`
	OpeningBalance decimal.Decimal         `json:"opening_balance"`
	TotalDebit     decimal.Decimal         `json:"total_debit"`
	TotalCredit    decimal.Decimal         `json:"total_credit"`
	ClosingBalance decimal.Decimal         `json:"closing_balance"`
	Transactions   []LedgerTransactionLine `json:"transactions"`
}

// GeneralLedgerFilter parameters for querying ledger.
type GeneralLedgerFilter struct {
	BranchID    *uuid.UUID
	AccountID   *uuid.UUID
	AccountCode string
	StartDate   time.Time
	EndDate     time.Time
}

// TrialBalanceItem represents a single account line in the Trial Balance.
type TrialBalanceItem struct {
	AccountID     uuid.UUID       `json:"account_id"`
	AccountCode   string          `json:"account_code"`
	AccountName   string          `json:"account_name"`
	AccountType   AccountType     `json:"account_type"`
	NormalBalance NormalBalance   `json:"normal_balance"`
	OpeningDebit  decimal.Decimal `json:"opening_debit"`
	OpeningCredit decimal.Decimal `json:"opening_credit"`
	PeriodDebit   decimal.Decimal `json:"period_debit"`
	PeriodCredit  decimal.Decimal `json:"period_credit"`
	EndingDebit   decimal.Decimal `json:"ending_debit"`
	EndingCredit  decimal.Decimal `json:"ending_credit"`
}

// TrialBalanceReport contains the aggregated Trial Balance.
type TrialBalanceReport struct {
	BranchID          *uuid.UUID         `json:"branch_id,omitempty"`
	StartDate         time.Time          `json:"start_date"`
	EndDate           time.Time          `json:"end_date"`
	TotalOpeningDebit decimal.Decimal    `json:"total_opening_debit"`
	TotalOpeningCredit decimal.Decimal   `json:"total_opening_credit"`
	TotalPeriodDebit  decimal.Decimal    `json:"total_period_debit"`
	TotalPeriodCredit decimal.Decimal    `json:"total_period_credit"`
	TotalEndingDebit  decimal.Decimal    `json:"total_ending_debit"`
	TotalEndingCredit decimal.Decimal    `json:"total_ending_credit"`
	IsBalanced        bool               `json:"is_balanced"`
	Items             []TrialBalanceItem `json:"items"`
}

// FinancialReportItem represents an account or subcategory row in financial statements.
type FinancialReportItem struct {
	AccountCode string          `json:"account_code"`
	AccountName string          `json:"account_name"`
	Amount      decimal.Decimal `json:"amount"`
}

// FinancialReportCategory groups items (e.g. Operating Revenues, Operating Expenses).
type FinancialReportCategory struct {
	CategoryName string                `json:"category_name"`
	Total        decimal.Decimal       `json:"total"`
	Items        []FinancialReportItem `json:"items"`
}

// IncomeStatement represents the Laporan Laba Rugi.
type IncomeStatement struct {
	BranchID           *uuid.UUID                `json:"branch_id,omitempty"`
	StartDate          time.Time                 `json:"start_date"`
	EndDate            time.Time                 `json:"end_date"`
	RevenueCategories  []FinancialReportCategory `json:"revenue_categories"`
	TotalRevenue       decimal.Decimal           `json:"total_revenue"`
	COGSCategories     []FinancialReportCategory `json:"cogs_categories"`
	TotalCOGS          decimal.Decimal           `json:"total_cogs"`
	GrossProfit        decimal.Decimal           `json:"gross_profit"`
	ExpenseCategories  []FinancialReportCategory `json:"expense_categories"`
	TotalExpenses      decimal.Decimal           `json:"total_expenses"`
	NetIncome          decimal.Decimal           `json:"net_income"` // Gross Profit - Total Expenses
}

// BalanceSheet represents the Laporan Neraca / Posisi Keuangan.
type BalanceSheet struct {
	BranchID             *uuid.UUID                `json:"branch_id,omitempty"`
	AsOfDate             time.Time                 `json:"as_of_date"`
	CurrentAssets        FinancialReportCategory   `json:"current_assets"`
	NonCurrentAssets     FinancialReportCategory   `json:"non_current_assets"`
	TotalAssets          decimal.Decimal           `json:"total_assets"`
	CurrentLiabilities   FinancialReportCategory   `json:"current_liabilities"`
	LongTermLiabilities  FinancialReportCategory   `json:"long_term_liabilities"`
	TotalLiabilities     decimal.Decimal           `json:"total_liabilities"`
	EquityCategories     FinancialReportCategory   `json:"equity_categories"`
	CurrentPeriodProfit  decimal.Decimal           `json:"current_period_profit"`
	TotalEquity          decimal.Decimal           `json:"total_equity"`
	TotalLiabilitiesEquity decimal.Decimal         `json:"total_liabilities_and_equity"`
	IsBalanced           bool                      `json:"is_balanced"` // Assets == Liabilities + Equity
}

// CashFlowStatement represents Laporan Arus Kas (Operating, Investing, Financing).
type CashFlowStatement struct {
	BranchID             *uuid.UUID              `json:"branch_id,omitempty"`
	StartDate            time.Time               `json:"start_date"`
	EndDate              time.Time               `json:"end_date"`
	OperatingActivities  FinancialReportCategory `json:"operating_activities"`
	NetOperatingCash     decimal.Decimal         `json:"net_operating_cash"`
	InvestingActivities  FinancialReportCategory `json:"investing_activities"`
	NetInvestingCash     decimal.Decimal         `json:"net_investing_cash"`
	FinancingActivities  FinancialReportCategory `json:"financing_activities"`
	NetFinancingCash     decimal.Decimal         `json:"net_financing_cash"`
	BeginningCash        decimal.Decimal         `json:"beginning_cash"`
	NetCashChange        decimal.Decimal         `json:"net_cash_change"`
	EndingCash           decimal.Decimal         `json:"ending_cash"`
}


