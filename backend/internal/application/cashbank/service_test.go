package cashbank_test

import (
	"context"
	"testing"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	cashbankApp "clinic-api/internal/application/cashbank"
	accountingDomain "clinic-api/internal/domain/accounting"
	cashbankDomain "clinic-api/internal/domain/cashbank"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockCashBankRepo struct {
	accounts     map[uuid.UUID]*cashbankDomain.CashBankAccount
	transactions []*cashbankDomain.CashBankTransaction
}

func newMockCashBankRepo() *mockCashBankRepo {
	return &mockCashBankRepo{
		accounts: make(map[uuid.UUID]*cashbankDomain.CashBankAccount),
	}
}

func (m *mockCashBankRepo) CreateAccount(ctx context.Context, acc *cashbankDomain.CashBankAccount) error {
	if acc.ID == uuid.Nil {
		acc.ID = uuid.New()
	}
	m.accounts[acc.ID] = acc
	return nil
}

func (m *mockCashBankRepo) GetAccountByID(ctx context.Context, id uuid.UUID) (*cashbankDomain.CashBankAccount, error) {
	acc, ok := m.accounts[id]
	if !ok {
		return nil, cashbankApp.ErrAccountNotFound
	}
	return acc, nil
}

func (m *mockCashBankRepo) GetAccountByCode(ctx context.Context, code string) (*cashbankDomain.CashBankAccount, error) {
	for _, a := range m.accounts {
		if a.Code == code {
			return a, nil
		}
	}
	return nil, cashbankApp.ErrAccountNotFound
}

func (m *mockCashBankRepo) ListAccounts(ctx context.Context, filter cashbankDomain.CashAccountFilter) ([]cashbankDomain.CashBankAccount, int64, error) {
	var list []cashbankDomain.CashBankAccount
	for _, a := range m.accounts {
		list = append(list, *a)
	}
	return list, int64(len(list)), nil
}

func (m *mockCashBankRepo) CreateTransaction(ctx context.Context, tx *cashbankDomain.CashBankTransaction) error {
	if tx.ID == uuid.Nil {
		tx.ID = uuid.New()
	}
	m.transactions = append(m.transactions, tx)
	return nil
}

func (m *mockCashBankRepo) GetTransactionByID(ctx context.Context, id uuid.UUID) (*cashbankDomain.CashBankTransaction, error) {
	for _, t := range m.transactions {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, cashbankApp.ErrAccountNotFound
}

func (m *mockCashBankRepo) ListTransactions(ctx context.Context, filter cashbankDomain.TransactionFilter) ([]cashbankDomain.CashBankTransaction, int64, error) {
	var list []cashbankDomain.CashBankTransaction
	for _, t := range m.transactions {
		list = append(list, *t)
	}
	return list, int64(len(list)), nil
}

// Mock Accounting Repositories
type mockAccountRepo struct {
	accounts map[uuid.UUID]*accountingDomain.Account
}

func newMockAccountRepo() *mockAccountRepo {
	return &mockAccountRepo{accounts: make(map[uuid.UUID]*accountingDomain.Account)}
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.Account, error) {
	a, ok := m.accounts[id]
	if !ok {
		return nil, accountingApp.ErrAccountNotFound
	}
	return a, nil
}
func (m *mockAccountRepo) GetByCode(ctx context.Context, code string) (*accountingDomain.Account, error) {
	for _, a := range m.accounts {
		if a.Code == code {
			return a, nil
		}
	}
	return nil, accountingApp.ErrAccountNotFound
}
func (m *mockAccountRepo) List(ctx context.Context, filter accountingDomain.AccountFilter) ([]accountingDomain.Account, int64, error) {
	var list []accountingDomain.Account
	for _, a := range m.accounts {
		list = append(list, *a)
	}
	return list, int64(len(list)), nil
}
func (m *mockAccountRepo) Create(ctx context.Context, acc *accountingDomain.Account) error {
	m.accounts[acc.ID] = acc
	return nil
}
func (m *mockAccountRepo) Update(ctx context.Context, acc *accountingDomain.Account) error {
	m.accounts[acc.ID] = acc
	return nil
}

type mockJournalRepo struct {
	journals map[uuid.UUID]*accountingDomain.JournalEntry
}

func (m *mockJournalRepo) Create(ctx context.Context, entry *accountingDomain.JournalEntry) error {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	m.journals[entry.ID] = entry
	return nil
}
func (m *mockJournalRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.JournalEntry, error) {
	j, ok := m.journals[id]
	if !ok {
		return nil, accountingApp.ErrInvalidJournalStatus
	}
	return j, nil
}
func (m *mockJournalRepo) GetByJournalNumber(ctx context.Context, number string) (*accountingDomain.JournalEntry, error) {
	for _, j := range m.journals {
		if j.JournalNumber == number {
			return j, nil
		}
	}
	return nil, accountingApp.ErrInvalidJournalStatus
}
func (m *mockJournalRepo) GetByReference(ctx context.Context, refType, refID string) (*accountingDomain.JournalEntry, error) {
	for _, j := range m.journals {
		if j.ReferenceType == refType && j.ReferenceID == refID {
			return j, nil
		}
	}
	return nil, accountingApp.ErrInvalidJournalStatus
}
func (m *mockJournalRepo) List(ctx context.Context, filter accountingDomain.JournalFilter) ([]accountingDomain.JournalEntry, int64, error) {
	var list []accountingDomain.JournalEntry
	for _, j := range m.journals {
		list = append(list, *j)
	}
	return list, int64(len(list)), nil
}
func (m *mockJournalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status accountingDomain.JournalStatus, postedBy *uuid.UUID) error {
	j, ok := m.journals[id]
	if !ok {
		return accountingApp.ErrInvalidJournalStatus
	}
	j.Status = status
	return nil
}

type mockFiscalPeriodRepo struct{}

func (m *mockFiscalPeriodRepo) GetPeriod(ctx context.Context, branchID uuid.UUID, year, month int) (*accountingDomain.FiscalPeriod, error) {
	return &accountingDomain.FiscalPeriod{Status: "OPEN"}, nil
}
func (m *mockFiscalPeriodRepo) ClosePeriod(ctx context.Context, period *accountingDomain.FiscalPeriod) error {
	return nil
}
func (m *mockFiscalPeriodRepo) ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error {
	return nil
}
func (m *mockFiscalPeriodRepo) IsPeriodLocked(ctx context.Context, branchID uuid.UUID, date time.Time) (bool, error) {
	return false, nil
}

// --- UNIT TESTS ---

func setupCashBankTest() (*cashbankApp.Service, *mockCashBankRepo, *mockAccountRepo, uuid.UUID, *cashbankDomain.CashBankAccount, *cashbankDomain.CashBankAccount, *accountingDomain.Account) {
	uow := &mockUoW{}
	cbRepo := newMockCashBankRepo()
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	cbSvc := cashbankApp.NewService(uow, cbRepo, accSvc)

	branchID := uuid.New()

	// Seed COA Kas Kasir (1-1111)
	kasCOA := &accountingDomain.Account{
		Code: "1-1111", Name: "Kas Kasir", Type: accountingDomain.AccountTypeAsset,
		NormalBalance: accountingDomain.NormalBalanceDebit, IsActive: true,
	}
	kasCOA.ID = uuid.New()
	_ = accRepo.Create(context.Background(), kasCOA)

	// Seed COA Bank BCA (1-1113)
	bcaCOA := &accountingDomain.Account{
		Code: "1-1113", Name: "Bank BCA", Type: accountingDomain.AccountTypeAsset,
		NormalBalance: accountingDomain.NormalBalanceDebit, IsActive: true,
	}
	bcaCOA.ID = uuid.New()
	_ = accRepo.Create(context.Background(), bcaCOA)

	// Seed COA Beban Listrik (6-1300)
	listrikCOA := &accountingDomain.Account{
		Code: "6-1300", Name: "Beban Listrik", Type: accountingDomain.AccountTypeExpense,
		NormalBalance: accountingDomain.NormalBalanceDebit, IsActive: true,
	}
	listrikCOA.ID = uuid.New()
	_ = accRepo.Create(context.Background(), listrikCOA)

	// Seed Cash Drawer Kasir
	cashDrawer := &cashbankDomain.CashBankAccount{
		AccountID: kasCOA.ID, AccountCode: kasCOA.Code, Code: "KASIR-01", Name: "Kasir Utama",
		Type: cashbankDomain.AccountTypeCash, IsActive: true,
	}
	cashDrawer.ID = uuid.New()
	_ = cbRepo.CreateAccount(context.Background(), cashDrawer)

	// Seed Bank Account BCA
	bankAccount := &cashbankDomain.CashBankAccount{
		AccountID: bcaCOA.ID, AccountCode: bcaCOA.Code, Code: "BCA-OPS", Name: "BCA Operasional",
		Type: cashbankDomain.AccountTypeBank, BankName: "Bank BCA", AccountNumber: "8881234567", IsActive: true,
	}
	bankAccount.ID = uuid.New()
	_ = cbRepo.CreateAccount(context.Background(), bankAccount)

	return cbSvc, cbRepo, accRepo, branchID, cashDrawer, bankAccount, listrikCOA
}

func TestRecordCashIn_AutoJournal(t *testing.T) {
	cbSvc, _, _, branchID, cashDrawer, _, _ := setupCashBankTest()
	ctx := context.Background()

	// Use bank BCA COA as counterpart or create additional revenue account
	counterpartID := uuid.New()
	cbSvc.CreateAccount(ctx, &cashbankDomain.CashBankAccount{}) // ensure initialized

	cmd := cashbankApp.RecordCashInCommand{
		BranchID:             branchID,
		CashBankAccountID:    cashDrawer.ID,
		Amount:               decimal.FromInt(1500000),
		TransactionDate:      time.Now(),
		CounterpartAccountID: cashDrawer.AccountID, // Kas
		Description:          "Setoran modal awal kasir",
	}

	txItem, journal, err := cbSvc.RecordCashIn(ctx, cmd)
	if err != nil {
		t.Fatalf("record cash in failed: %v", err)
	}

	if txItem.TransactionType != cashbankDomain.TransactionTypeCashIn {
		t.Errorf("expected CASH_IN, got %s", txItem.TransactionType)
	}
	if journal == nil || journal.Status != accountingDomain.JournalStatusPosted {
		t.Fatalf("expected POSTED journal entry")
	}
	if !journal.TotalDebit.Equal(decimal.FromInt(1500000)) || !journal.TotalCredit.Equal(decimal.FromInt(1500000)) {
		t.Errorf("unbalanced journal: debit=%s, credit=%s", journal.TotalDebit, journal.TotalCredit)
	}
	_ = counterpartID
}

func TestRecordCashOut_AutoJournal(t *testing.T) {
	cbSvc, _, _, branchID, cashDrawer, _, listrikCOA := setupCashBankTest()
	ctx := context.Background()

	cmd := cashbankApp.RecordCashOutCommand{
		BranchID:             branchID,
		CashBankAccountID:    cashDrawer.ID,
		Amount:               decimal.FromInt(350000),
		TransactionDate:      time.Now(),
		CounterpartAccountID: listrikCOA.ID,
		Description:          "Pembayaran tagihan listrik PLN bulan ini",
	}

	txItem, journal, err := cbSvc.RecordCashOut(ctx, cmd)
	if err != nil {
		t.Fatalf("record cash out failed: %v", err)
	}

	if txItem.TransactionType != cashbankDomain.TransactionTypeCashOut {
		t.Errorf("expected CASH_OUT, got %s", txItem.TransactionType)
	}
	if journal.JournalType != accountingDomain.JournalTypeCashDisbursement {
		t.Errorf("expected CASH_DISBURSEMENT journal type, got %s", journal.JournalType)
	}
	if !journal.TotalDebit.Equal(decimal.FromInt(350000)) {
		t.Errorf("expected debit 350000, got %s", journal.TotalDebit)
	}
}

func TestTransferCashBank_Success(t *testing.T) {
	cbSvc, _, _, branchID, cashDrawer, bankAccount, _ := setupCashBankTest()
	ctx := context.Background()

	// Transfer Rp 1.000.000 from Kasir to Bank BCA
	cmd := cashbankApp.TransferCommand{
		BranchID:        branchID,
		SourceAccountID: cashDrawer.ID,
		TargetAccountID: bankAccount.ID,
		Amount:          decimal.FromInt(1000000),
		TransactionDate: time.Now(),
		Description:     "Setoran tunai kasir ke rekening Bank BCA",
	}

	txItem, journal, err := cbSvc.TransferCashBank(ctx, cmd)
	if err != nil {
		t.Fatalf("transfer failed: %v", err)
	}

	if txItem.TransactionType != cashbankDomain.TransactionTypeTransfer {
		t.Errorf("expected TRANSFER, got %s", txItem.TransactionType)
	}
	if journal.JournalType != accountingDomain.JournalTypeGeneral {
		t.Errorf("expected GENERAL journal, got %s", journal.JournalType)
	}
	// Verify lines: Debit Bank BCA, Credit Kas Kasir
	if !journal.Lines[0].Debit.Equal(decimal.FromInt(1000000)) {
		t.Errorf("expected first line debit 1000000 (Bank BCA), got %s", journal.Lines[0].Debit)
	}
	if !journal.Lines[1].Credit.Equal(decimal.FromInt(1000000)) {
		t.Errorf("expected second line credit 1000000 (Kas Kasir), got %s", journal.Lines[1].Credit)
	}
}

func TestTransferCashBank_SameAccountRejected(t *testing.T) {
	cbSvc, _, _, branchID, cashDrawer, _, _ := setupCashBankTest()
	ctx := context.Background()

	cmd := cashbankApp.TransferCommand{
		BranchID:        branchID,
		SourceAccountID: cashDrawer.ID,
		TargetAccountID: cashDrawer.ID, // Same account!
		Amount:          decimal.FromInt(500000),
	}

	_, _, err := cbSvc.TransferCashBank(ctx, cmd)
	if err == nil {
		t.Fatalf("expected error transferring to same account, got nil")
	}
}

