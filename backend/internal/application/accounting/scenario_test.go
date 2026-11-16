package accounting_test

import (
	"context"
	"testing"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

func TestComprehensiveClinicScenario_KidinestAndSSC_EndToEnd(t *testing.T) {
	uow := &mockUoW{}
	accRepo := newFullMockAccountRepo()
	jRepo := newFullMockJournalRepo()
	pRepo := &mockFiscalPeriodRepo{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	ctx := context.Background()
	branchID := uuid.New()

	// -------------------------------------------------------------
	// SCENARIO STEP 1: INITIAL CAPITAL & BANK DEPOSIT
	// Pemilik menyetorkan modal awal Rp 500,000,000 ke Bank BCA
	// Debit: 1-1121 Bank BCA (500M)
	// Credit: 3-1100 Modal Disetor (500M)
	// -------------------------------------------------------------
	_, err := accSvc.CreateJournal(ctx, accountingApp.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     accountingDomain.JournalTypeBank,
		TransactionDate: time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC),
		ReferenceType:   "CAPITAL_INJECTION",
		ReferenceID:     "CAP-001",
		Description:     "Setoran Modal Disetor Pemilik Klinik",
		AutoPost:        true,
		Lines: []accountingApp.CreateLineCommand{
			{AccountCode: "1-1121", Debit: decimal.FromInt(500000000), Credit: decimal.Zero, Description: "Bank BCA"},
			{AccountCode: "3-1100", Debit: decimal.Zero, Credit: decimal.FromInt(500000000), Description: "Modal Disetor"},
		},
	})
	if err != nil {
		t.Fatalf("step 1 failed: %v", err)
	}

	// -------------------------------------------------------------
	// SCENARIO STEP 2: PREPAID EXPENSE (SEWA GEDUNG KLINIK)
	// Pembayaran Sewa Gedung Rp 120,000,000 untuk 12 bulan (10M/bln)
	// Debit: 1-1400 Biaya Dibayar di Muka (120M)
	// Credit: 1-1121 Bank BCA (120M)
	// -------------------------------------------------------------
	_, err = accSvc.CreateJournal(ctx, accountingApp.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     accountingDomain.JournalTypeCashDisbursement,
		TransactionDate: time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC),
		ReferenceType:   "PREPAID_RENT",
		ReferenceID:     "RENT-2026",
		Description:     "Pembayaran Sewa Gedung Klinik Kidinest & SSC 1 Tahun",
		AutoPost:        true,
		Lines: []accountingApp.CreateLineCommand{
			{AccountCode: "1-1400", Debit: decimal.FromInt(120000000), Credit: decimal.Zero, Description: "Sewa Dibayar di Muka"},
			{AccountCode: "1-1121", Debit: decimal.Zero, Credit: decimal.FromInt(120000000), Description: "Bank BCA Operasional"},
		},
	})
	if err != nil {
		t.Fatalf("step 2 failed: %v", err)
	}

	// -------------------------------------------------------------
	// SCENARIO STEP 3: FIXED ASSET ACQUISITION (MESIN LASER SSC)
	// Beli Mesin Laser Estetika SSC Rp 240,000,000 via Bank BCA
	// Debit: 1-2200 Peralatan Estetika SSC (240M)
	// Credit: 1-1121 Bank BCA (240M)
	// -------------------------------------------------------------
	_, err = accSvc.CreateJournal(ctx, accountingApp.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     accountingDomain.JournalTypeBank,
		TransactionDate: time.Date(2026, 1, 3, 11, 0, 0, 0, time.UTC),
		ReferenceType:   "ASSET_PURCHASE",
		ReferenceID:     "LASER-001",
		Description:     "Pembelian Mesin Laser CO2 Semarang Skin Center",
		AutoPost:        true,
		Lines: []accountingApp.CreateLineCommand{
			{AccountCode: "1-2200", Debit: decimal.FromInt(240000000), Credit: decimal.Zero, Description: "Peralatan Laser SSC"},
			{AccountCode: "1-1121", Debit: decimal.Zero, Credit: decimal.FromInt(240000000), Description: "Bank BCA Operasional"},
		},
	})
	if err != nil {
		t.Fatalf("step 3 failed: %v", err)
	}

	// -------------------------------------------------------------
	// SCENARIO STEP 4: KIDINEST BILLING & POS (PRE-SERVICE)
	// Pasien Anak Kevin: Registrasi 25k, Konsultasi Sp.A 200k, Obat 175k = Total 400k (Cash)
	// COGS Obat: 75k (HPP 5-1110 vs Persediaan 1-1310)
	// Doctor Share Accrual (70% dr Konsultasi): 140k (Beban 5-1210 vs Hutang 2-1210)
	// -------------------------------------------------------------
	_, err = accSvc.CreateJournal(ctx, accountingApp.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     accountingDomain.JournalTypeCashReceipt,
		TransactionDate: time.Date(2026, 1, 15, 14, 0, 0, 0, time.UTC),
		ReferenceType:   "INVOICE_PAYMENT",
		ReferenceID:     "INV-KID-001",
		Description:     "Pelunasan Pasien Anak Kevin (Kidinest)",
		AutoPost:        true,
		Lines: []accountingApp.CreateLineCommand{
			{AccountCode: "1-1111", Debit: decimal.FromInt(400000), Credit: decimal.Zero, Description: "Kas Kasir Kidinest"},
			{AccountCode: "4-1110", Debit: decimal.Zero, Credit: decimal.FromInt(25000), Description: "Pendapatan Pendaftaran Anak"},
			{AccountCode: "4-1120", Debit: decimal.Zero, Credit: decimal.FromInt(200000), Description: "Pendapatan Konsultasi Sp.A"},
			{AccountCode: "4-1150", Debit: decimal.Zero, Credit: decimal.FromInt(175000), Description: "Pendapatan Farmasi Anak"},
			// HPP & Persediaan
			{AccountCode: "5-1110", Debit: decimal.FromInt(75000), Credit: decimal.Zero, Description: "HPP Obat Sirup Anak"},
			{AccountCode: "1-1310", Debit: decimal.Zero, Credit: decimal.FromInt(75000), Description: "Persediaan Obat Paten & Sirup"},
			// Doctor Share Accrual
			{AccountCode: "5-1210", Debit: decimal.FromInt(140000), Credit: decimal.Zero, Description: "Beban Jasa Medis dr. Sp.A"},
			{AccountCode: "2-1210", Debit: decimal.Zero, Credit: decimal.FromInt(140000), Description: "Hutang Jasa Medis dr. Sp.A"},
		},
	})
	if err != nil {
		t.Fatalf("step 4 failed: %v", err)
	}

	// -------------------------------------------------------------
	// SCENARIO STEP 5: SEMARANG SKIN CENTER (SSC) BILLING (POST-SERVICE)
	// Pasien Estetika Amanda: Konsultasi Sp.DVE 300k, Tindakan Laser 1,500k, Skincare 200k = Total 2,000,000 (Bank Transfer)
	// COGS Skincare: 80k (HPP 5-2110 vs Persediaan 1-1330)
	// Doctor Share Accrual: 1,000,000 (Beban 5-2210 vs Hutang 2-1220)
	// -------------------------------------------------------------
	_, err = accSvc.CreateJournal(ctx, accountingApp.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     accountingDomain.JournalTypeBank,
		TransactionDate: time.Date(2026, 1, 16, 16, 0, 0, 0, time.UTC),
		ReferenceType:   "INVOICE_PAYMENT",
		ReferenceID:     "INV-SSC-001",
		Description:     "Pelunasan Pasien Estetika Amanda (SSC)",
		AutoPost:        true,
		Lines: []accountingApp.CreateLineCommand{
			{AccountCode: "1-1121", Debit: decimal.FromInt(2000000), Credit: decimal.Zero, Description: "Bank BCA Operasional"},
			{AccountCode: "4-2110", Debit: decimal.Zero, Credit: decimal.FromInt(300000), Description: "Pendapatan Konsultasi Sp.DVE"},
			{AccountCode: "4-2130", Debit: decimal.Zero, Credit: decimal.FromInt(1500000), Description: "Pendapatan Tindakan Laser SSC"},
			{AccountCode: "4-2140", Debit: decimal.Zero, Credit: decimal.FromInt(200000), Description: "Pendapatan Produk Skincare SSC"},
			// HPP & Persediaan
			{AccountCode: "5-2110", Debit: decimal.FromInt(80000), Credit: decimal.Zero, Description: "HPP Produk Skincare SSC"},
			{AccountCode: "1-1330", Debit: decimal.Zero, Credit: decimal.FromInt(80000), Description: "Persediaan Skincare SSC"},
			// Doctor Share Accrual
			{AccountCode: "5-2210", Debit: decimal.FromInt(1000000), Credit: decimal.Zero, Description: "Beban Jasa Medis dr. Kulit SSC"},
			{AccountCode: "2-1220", Debit: decimal.Zero, Credit: decimal.FromInt(1000000), Description: "Hutang Jasa Medis dr. Kulit SSC"},
		},
	})
	if err != nil {
		t.Fatalf("step 5 failed: %v", err)
	}

	// -------------------------------------------------------------
	// SCENARIO STEP 6: MONTH-END ADJUSTING ENTRIES (JANUARI 2026)
	// 1. Amortisasi Sewa Gedung: 10,000,000 (Beban 6-1200 vs Biaya Dimuka 1-1400)
	// 2. Penyusutan Mesin Laser: 240M / 48 bln = 5,000,000 (Beban 6-1300 vs Akumulasi 1-2900)
	// -------------------------------------------------------------
	_, err = accSvc.CreateJournal(ctx, accountingApp.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     accountingDomain.JournalTypeAdjusting,
		TransactionDate: time.Date(2026, 1, 31, 23, 0, 0, 0, time.UTC),
		ReferenceType:   "MONTH_END_ADJ",
		ReferenceID:     "ADJ-2026-01",
		Description:     "Jurnal Penyesuaian Sewa & Depresiasi Laser Januari 2026",
		AutoPost:        true,
		Lines: []accountingApp.CreateLineCommand{
			{AccountCode: "6-1200", Debit: decimal.FromInt(10000000), Credit: decimal.Zero, Description: "Beban Sewa Gedung"},
			{AccountCode: "1-1400", Debit: decimal.Zero, Credit: decimal.FromInt(10000000), Description: "Sewa Dibayar di Muka"},
			{AccountCode: "6-1300", Debit: decimal.FromInt(5000000), Credit: decimal.Zero, Description: "Beban Penyusutan Mesin Laser"},
			{AccountCode: "1-2900", Debit: decimal.Zero, Credit: decimal.FromInt(5000000), Description: "Akumulasi Penyusutan Mesin Laser"},
		},
	})
	if err != nil {
		t.Fatalf("step 6 failed: %v", err)
	}

	// -------------------------------------------------------------
	// SCENARIO STEP 7: DOCTOR SHARE SETTLEMENT & PPh 21 (PMK 168)
	// Pelunasan Fee dr. Kulit SSC Bruto 1,000,000
	// DPP 50% = 500,000; PPh 21 (5%) = 25,000; Net Fee = 975,000
	// Debit: 2-1220 Hutang Dokter SSC (1,000,000)
	// Credit: 1-1121 Bank BCA (975,000)
	// Credit: 2-1501.02 Hutang PPh 21 Dokter (25,000)
	// -------------------------------------------------------------
	_, err = accSvc.CreateJournal(ctx, accountingApp.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     accountingDomain.JournalTypeCashDisbursement,
		TransactionDate: time.Date(2026, 1, 31, 23, 30, 0, 0, time.UTC),
		ReferenceType:   "DOCTOR_FEE_PAYMENT",
		ReferenceID:     "DOCPAY-2026-01",
		Description:     "Pembayaran Fee Dokter Kulit SSC & Potongan PPh 21",
		AutoPost:        true,
		Lines: []accountingApp.CreateLineCommand{
			{AccountCode: "2-1220", Debit: decimal.FromInt(1000000), Credit: decimal.Zero, Description: "Pelunasan Hutang Fee Dokter SSC"},
			{AccountCode: "1-1121", Debit: decimal.Zero, Credit: decimal.FromInt(975000), Description: "Transfer Bersih via BCA"},
			{AccountCode: "2-1501.02", Debit: decimal.Zero, Credit: decimal.FromInt(25000), Description: "Potongan PPh 21 Tenaga Medis"},
		},
	})
	if err != nil {
		t.Fatalf("step 7 failed: %v", err)
	}

	// -------------------------------------------------------------
	// SCENARIO STEP 8: FINANCIAL INTEGRITY & BALANCE CHECK
	// Verifikasi bahwa seluruh jurnal yang diposting 100% seimbang
	// Total Debit == Total Credit di seluruh transaksi
	// -------------------------------------------------------------
	totalAllDebit := decimal.Zero
	totalAllCredit := decimal.Zero
	for _, j := range jRepo.journals {
		totalAllDebit = totalAllDebit.Add(j.TotalDebit)
		totalAllCredit = totalAllCredit.Add(j.TotalCredit)
		if !j.TotalDebit.Equal(j.TotalCredit) {
			t.Errorf("unbalanced journal %s: dr=%s cr=%s", j.JournalNumber, j.TotalDebit, j.TotalCredit)
		}
	}

	if !totalAllDebit.Equal(totalAllCredit) {
		t.Fatalf("global trial balance imbalance: debit=%s credit=%s", totalAllDebit, totalAllCredit)
	}

	// -------------------------------------------------------------
	// SCENARIO STEP 9: YEAR-END CLOSING (TUTUP BUKU TAHUNAN)
	// Tutup buku tahun 2026: Semua akun nominal (4-xxxx, 5-xxxx, 6-xxxx)
	// dialihkan ke 3-1200 Laba Ditahan.
	// -------------------------------------------------------------
	// Hitung total pendapatan dan total beban
	// Pendapatan Kidinest: 25k + 200k + 175k = 400,000
	// Pendapatan SSC: 300k + 1,500k + 200k = 2,000,000
	// Total Pendapatan = 2,400,000
	// Total Beban & HPP:
	// HPP Obat Anak: 75,000
	// Jasa Medis Dokter Anak: 140,000
	// HPP Skincare SSC: 80,000
	// Jasa Medis Dokter SSC: 1,000,000
	// Beban Sewa Gedung: 10,000,000
	// Beban Penyusutan Laser: 5,000,000
	// Total Beban = 16,295,000
	// Laba/Rugi Bersih = 2,400,000 - 16,295,000 = -13,895,000 (Rugi Operasional Tahun Pertama)
	// Closing journal:
	// Debit Pendapatan (2,400,000)
	// Debit Laba Ditahan / Rugi (13,895,000)
	// Credit Beban-Beban (16,295,000)
	// Total Debit = 16,295,000; Total Credit = 16,295,000
	closingJ, err := accSvc.CreateJournal(ctx, accountingApp.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     accountingDomain.JournalTypeClosing,
		TransactionDate: time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC),
		ReferenceType:   "YEAR_END_CLOSING",
		ReferenceID:     "YEC-2026",
		Description:     "Jurnal Penutup Akhir Tahun 2026 ke Laba Ditahan",
		AutoPost:        true,
		Lines: []accountingApp.CreateLineCommand{
			// Tutup Pendapatan (Debit)
			{AccountCode: "4-1110", Debit: decimal.FromInt(25000), Credit: decimal.Zero, Description: "Tutup Pendapatan Pendaftaran"},
			{AccountCode: "4-1120", Debit: decimal.FromInt(200000), Credit: decimal.Zero, Description: "Tutup Pendapatan Konsultasi Anak"},
			{AccountCode: "4-1150", Debit: decimal.FromInt(175000), Credit: decimal.Zero, Description: "Tutup Pendapatan Farmasi Anak"},
			{AccountCode: "4-2110", Debit: decimal.FromInt(300000), Credit: decimal.Zero, Description: "Tutup Pendapatan Sp.DVE"},
			{AccountCode: "4-2130", Debit: decimal.FromInt(1500000), Credit: decimal.Zero, Description: "Tutup Pendapatan Laser SSC"},
			{AccountCode: "4-2140", Debit: decimal.FromInt(200000), Credit: decimal.Zero, Description: "Tutup Pendapatan Skincare SSC"},
			// Balancing figure ke 3-1200 Laba Ditahan
			{AccountCode: "3-1200", Debit: decimal.FromInt(13895000), Credit: decimal.Zero, Description: "Transfer Hasil Periode ke Laba Ditahan"},
			// Tutup Beban-Beban (Credit)
			{AccountCode: "5-1110", Debit: decimal.Zero, Credit: decimal.FromInt(75000), Description: "Tutup HPP Obat"},
			{AccountCode: "5-1210", Debit: decimal.Zero, Credit: decimal.FromInt(140000), Description: "Tutup Fee Dokter Anak"},
			{AccountCode: "5-2110", Debit: decimal.Zero, Credit: decimal.FromInt(80000), Description: "Tutup HPP Skincare"},
			{AccountCode: "5-2210", Debit: decimal.Zero, Credit: decimal.FromInt(1000000), Description: "Tutup Fee Dokter Kulit"},
			{AccountCode: "6-1200", Debit: decimal.Zero, Credit: decimal.FromInt(10000000), Description: "Tutup Beban Sewa"},
			{AccountCode: "6-1300", Debit: decimal.Zero, Credit: decimal.FromInt(5000000), Description: "Tutup Beban Depresiasi"},
		},
	})
	if err != nil {
		t.Fatalf("year end closing failed: %v", err)
	}

	if !closingJ.TotalDebit.Equal(decimal.FromInt(16295000)) || !closingJ.TotalCredit.Equal(decimal.FromInt(16295000)) {
		t.Errorf("expected closing journal 16,295,000 balanced, got dr: %s, cr: %s", closingJ.TotalDebit, closingJ.TotalCredit)
	}

	t.Logf("✅ Comprehensive dual-unit clinic accounting scenario successfully completed and verified!")
}

type fullMockAccountRepo struct {
	accounts map[string]*accountingDomain.Account
}

func newFullMockAccountRepo() *fullMockAccountRepo {
	r := &fullMockAccountRepo{accounts: make(map[string]*accountingDomain.Account)}
	seed := []struct {
		code, name string
		accType    accountingDomain.AccountType
		normal     accountingDomain.NormalBalance
	}{
		{"1-1111", "Kas Kasir Kidinest", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1121", "Bank BCA Operasional", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1310", "Persediaan Obat Paten & Sirup Anak", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1330", "Persediaan Produk Skincare & Kosmetika SSC", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1400", "Biaya Sewa Dibayar di Muka", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-2200", "Peralatan Estetika SSC", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-2900", "Akumulasi Penyusutan Aset Tetap", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceCredit},
		{"2-1210", "Hutang Jasa Medis Dokter Kidinest", accountingDomain.AccountTypeLiability, accountingDomain.NormalBalanceCredit},
		{"2-1220", "Hutang Jasa Medis Dokter SSC", accountingDomain.AccountTypeLiability, accountingDomain.NormalBalanceCredit},
		{"2-1501.02", "Hutang PPh 21 Dokter", accountingDomain.AccountTypeLiability, accountingDomain.NormalBalanceCredit},
		{"3-1100", "Modal Disetor Pemilik", accountingDomain.AccountTypeEquity, accountingDomain.NormalBalanceCredit},
		{"3-1200", "Laba Ditahan (Retained Earnings)", accountingDomain.AccountTypeEquity, accountingDomain.NormalBalanceCredit},
		{"4-1110", "Pendapatan Pendaftaran Pasien Anak", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-1120", "Pendapatan Konsultasi Dokter Sp.A", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-1150", "Pendapatan Farmasi Anak", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-2110", "Pendapatan Konsultasi Dokter Estetika", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-2130", "Pendapatan Tindakan Medis Laser", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-2140", "Pendapatan Produk Skincare SSC", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"5-1110", "HPP Obat & Sirup Anak", accountingDomain.AccountTypeCOGS, accountingDomain.NormalBalanceDebit},
		{"5-1210", "Beban Jasa Medis Dokter Anak", accountingDomain.AccountTypeCOGS, accountingDomain.NormalBalanceDebit},
		{"5-2110", "HPP Produk Skincare SSC", accountingDomain.AccountTypeCOGS, accountingDomain.NormalBalanceDebit},
		{"5-2210", "Beban Jasa Medis Dokter SSC", accountingDomain.AccountTypeCOGS, accountingDomain.NormalBalanceDebit},
		{"6-1200", "Beban Sewa Gedung Klinik", accountingDomain.AccountTypeExpense, accountingDomain.NormalBalanceDebit},
		{"6-1300", "Beban Penyusutan Mesin Laser", accountingDomain.AccountTypeExpense, accountingDomain.NormalBalanceDebit},
	}
	for _, s := range seed {
		acc := &accountingDomain.Account{
			Code:          s.code,
			Name:          s.name,
			Type:          s.accType,
			NormalBalance: s.normal,
			IsActive:      true,
		}
		acc.ID = uuid.New()
		r.accounts[s.code] = acc
	}
	return r
}

func (m *fullMockAccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.Account, error) {
	for _, a := range m.accounts {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, nil
}

func (m *fullMockAccountRepo) GetByCode(ctx context.Context, code string) (*accountingDomain.Account, error) {
	if a, ok := m.accounts[code]; ok {
		return a, nil
	}
	return nil, nil
}

func (m *fullMockAccountRepo) List(ctx context.Context, filter accountingDomain.AccountFilter) ([]accountingDomain.Account, int64, error) {
	var list []accountingDomain.Account
	for _, a := range m.accounts {
		list = append(list, *a)
	}
	return list, int64(len(list)), nil
}

func (m *fullMockAccountRepo) Create(ctx context.Context, acc *accountingDomain.Account) error {
	m.accounts[acc.Code] = acc
	return nil
}

func (m *fullMockAccountRepo) Update(ctx context.Context, acc *accountingDomain.Account) error {
	m.accounts[acc.Code] = acc
	return nil
}

type fullMockJournalRepo struct {
	journals map[uuid.UUID]*accountingDomain.JournalEntry
}

func newFullMockJournalRepo() *fullMockJournalRepo {
	return &fullMockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
}

func (m *fullMockJournalRepo) Create(ctx context.Context, entry *accountingDomain.JournalEntry) error {
	m.journals[entry.ID] = entry
	return nil
}

func (m *fullMockJournalRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.JournalEntry, error) {
	return m.journals[id], nil
}

func (m *fullMockJournalRepo) GetByJournalNumber(ctx context.Context, number string) (*accountingDomain.JournalEntry, error) {
	return nil, nil
}

func (m *fullMockJournalRepo) GetByReference(ctx context.Context, refType, refID string) (*accountingDomain.JournalEntry, error) {
	return nil, nil
}

func (m *fullMockJournalRepo) List(ctx context.Context, filter accountingDomain.JournalFilter) ([]accountingDomain.JournalEntry, int64, error) {
	return nil, 0, nil
}

func (m *fullMockJournalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status accountingDomain.JournalStatus, postedBy *uuid.UUID) error {
	return nil
}

