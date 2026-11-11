package usecase_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"clinic-api/internal/pharmacy/domain"
	"clinic-api/internal/pharmacy/usecase"
)

type mockPharmacyRepo struct {
	prescriptions map[int64]*domain.PharmacyPrescription
	drugs         map[int64]*domain.Drug
	opnameLogs    []*domain.StockOpnameLog
}

func (m *mockPharmacyRepo) ListIncomingPrescriptions(ctx context.Context, clinicUnit string) ([]*domain.PharmacyPrescription, error) {
	var list []*domain.PharmacyPrescription
	for _, p := range m.prescriptions {
		if clinicUnit != "" && clinicUnit != "ALL" && p.ClinicUnit != clinicUnit {
			continue
		}
		list = append(list, p)
	}
	return list, nil
}

func (m *mockPharmacyRepo) GetPrescriptionByID(ctx context.Context, id int64) (*domain.PharmacyPrescription, error) {
	if p, ok := m.prescriptions[id]; ok {
		return p, nil
	}
	return nil, domain.ErrPrescriptionNotFound
}

func (m *mockPharmacyRepo) UpdatePrescriptionStatus(ctx context.Context, id int64, status string) error {
	if p, ok := m.prescriptions[id]; ok {
		p.Status = status
		return nil
	}
	return domain.ErrPrescriptionNotFound
}

func (m *mockPharmacyRepo) DispensePrescriptionTx(ctx context.Context, id int64) error {
	p, ok := m.prescriptions[id]
	if !ok {
		return domain.ErrPrescriptionNotFound
	}
	if p.Status == "DISPENSED" {
		return domain.ErrAlreadyDispensed
	}

	for _, it := range p.Items {
		if !it.IsManual && it.DrugID != nil {
			if d, ok := m.drugs[*it.DrugID]; ok {
				d.PhysicalStock -= it.Quantity
				d.ReservedStock -= it.Quantity
				d.AvailableStock = d.PhysicalStock - d.ReservedStock
			}
		}
	}
	p.Status = "DISPENSED"
	return nil
}

func (m *mockPharmacyRepo) CancelPrescriptionTx(ctx context.Context, id int64) error {
	p, ok := m.prescriptions[id]
	if !ok {
		return domain.ErrPrescriptionNotFound
	}
	if p.Status == "DISPENSED" {
		return domain.ErrAlreadyDispensed
	}

	for _, it := range p.Items {
		if !it.IsManual && it.DrugID != nil {
			if d, ok := m.drugs[*it.DrugID]; ok {
				d.ReservedStock -= it.Quantity
				d.AvailableStock = d.PhysicalStock - d.ReservedStock
			}
		}
	}
	p.Status = "CANCELLED"
	return nil
}

func (m *mockPharmacyRepo) ListDrugs(ctx context.Context, clinicUnit string) ([]*domain.Drug, error) {
	var list []*domain.Drug
	for _, d := range m.drugs {
		if clinicUnit != "" && clinicUnit != "ALL" && d.ClinicUnit != clinicUnit && d.ClinicUnit != "ALL" {
			continue
		}
		list = append(list, d)
	}
	return list, nil
}

func (m *mockPharmacyRepo) GetDrugByID(ctx context.Context, id int64) (*domain.Drug, error) {
	if d, ok := m.drugs[id]; ok {
		return d, nil
	}
	return nil, domain.ErrDrugNotFound
}

func (m *mockPharmacyRepo) CreateDrug(ctx context.Context, in *domain.CreateDrugInput) (*domain.Drug, error) {
	newID := int64(len(m.drugs) + 1)
	d := &domain.Drug{
		ID:            newID,
		ClinicUnit:    in.ClinicUnit,
		Name:          in.Name,
		Unit:          in.Unit,
		Category:      in.Category,
		MinStockAlert: in.MinStockAlert,
		BuyPrice:      in.BuyPrice,
		SellPrice:     in.SellPrice,
		IsActive:      true,
	}
	m.drugs[newID] = d
	return d, nil
}

func (m *mockPharmacyRepo) RestockBatch(ctx context.Context, in *domain.RestockInput) error {
	d, ok := m.drugs[in.DrugID]
	if !ok {
		return domain.ErrDrugNotFound
	}
	d.PhysicalStock += in.Quantity
	d.AvailableStock = d.PhysicalStock - d.ReservedStock
	d.Batches = append(d.Batches, domain.DrugBatch{
		ID:            int64(len(d.Batches) + 1),
		DrugID:        in.DrugID,
		BatchNumber:   in.BatchNumber,
		ExpiredDate:   in.ExpiredDate,
		PhysicalStock: in.Quantity,
	})
	return nil
}

func (m *mockPharmacyRepo) ApplyStockOpnameTx(ctx context.Context, items []domain.StockOpnameItem, adjustedBy string) error {
	for _, it := range items {
		if d, ok := m.drugs[it.DrugID]; ok {
			d.PhysicalStock = it.ActualStock
			d.AvailableStock = d.PhysicalStock - d.ReservedStock
		}
		m.opnameLogs = append(m.opnameLogs, &domain.StockOpnameLog{
			ID:          int64(len(m.opnameLogs) + 1),
			DrugID:      it.DrugID,
			DrugName:    it.DrugName,
			BatchNumber: it.BatchNumber,
			SystemStock: it.SystemStock,
			ActualStock: it.ActualStock,
			Discrepancy: it.Discrepancy,
			AdjustedBy:  adjustedBy,
			Notes:       it.Notes,
			CreatedAt:   time.Now(),
		})
	}
	return nil
}

func (m *mockPharmacyRepo) ListOpnameLogs(ctx context.Context, limit int) ([]*domain.StockOpnameLog, error) {
	return m.opnameLogs, nil
}

func setupPharmacyUsecase() (*usecase.PharmacyUsecase, *mockPharmacyRepo) {
	repo := &mockPharmacyRepo{
		prescriptions: make(map[int64]*domain.PharmacyPrescription),
		drugs:         make(map[int64]*domain.Drug),
		opnameLogs:    make([]*domain.StockOpnameLog, 0),
	}

	repo.drugs[1] = &domain.Drug{
		ID:             1,
		ClinicUnit:     "KIDINEST",
		Name:           "Paracetamol Drop 100mg/ml",
		Unit:           "Botol",
		Category:       "Antipiretik",
		PhysicalStock:  50,
		ReservedStock:  5,
		AvailableStock: 45,
		IsActive:       true,
		Batches: []domain.DrugBatch{
			{ID: 10, DrugID: 1, BatchNumber: "BCH-PCM-01", ExpiredDate: "2027-12-31", PhysicalStock: 50, ReservedStock: 5},
		},
	}

	drugID := int64(1)
	repo.prescriptions[10] = &domain.PharmacyPrescription{
		ID:          10,
		QueueID:     100,
		QueueNumber: "A-001",
		PatientName: "Kenzo",
		ClinicUnit:  "KIDINEST",
		Status:      "RESERVED",
		Items: []domain.PrescriptionItemDetail{
			{DrugID: &drugID, ItemName: "Paracetamol Drop", Quantity: 5, UnitPrice: 35000, Subtotal: 175000},
		},
	}

	uc := usecase.NewPharmacyUsecase(repo)
	return uc, repo
}

func TestDispensePrescriptionInvariant(t *testing.T) {
	uc, repo := setupPharmacyUsecase()
	ctx := context.Background()

	// Initially: Physical = 50, Reserved = 5, Available = 45
	err := uc.Dispense(ctx, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	d := repo.drugs[1]
	// After dispense: Physical = 45, Reserved = 0, Available = 45
	if d.PhysicalStock != 45 {
		t.Errorf("expected physical stock 45, got %d", d.PhysicalStock)
	}
	if d.ReservedStock != 0 {
		t.Errorf("expected reserved stock 0, got %d", d.ReservedStock)
	}
	if d.AvailableStock != 45 {
		t.Errorf("expected available stock 45, got %d", d.AvailableStock)
	}
	if repo.prescriptions[10].Status != "DISPENSED" {
		t.Errorf("expected status DISPENSED, got %s", repo.prescriptions[10].Status)
	}
}

func TestCancelPrescriptionInvariant(t *testing.T) {
	uc, repo := setupPharmacyUsecase()
	ctx := context.Background()

	// Initially: Physical = 50, Reserved = 5, Available = 45
	err := uc.Cancel(ctx, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	d := repo.drugs[1]
	// After cancel: Physical remains 50, Reserved = 0, Available = 50
	if d.PhysicalStock != 50 {
		t.Errorf("expected physical stock 50, got %d", d.PhysicalStock)
	}
	if d.ReservedStock != 0 {
		t.Errorf("expected reserved stock 0, got %d", d.ReservedStock)
	}
	if d.AvailableStock != 50 {
		t.Errorf("expected available stock 50, got %d", d.AvailableStock)
	}
	if repo.prescriptions[10].Status != "CANCELLED" {
		t.Errorf("expected status CANCELLED, got %s", repo.prescriptions[10].Status)
	}
}

func TestRestockValidation(t *testing.T) {
	uc, _ := setupPharmacyUsecase()
	ctx := context.Background()

	// Missing batch number
	err := uc.Restock(ctx, &domain.RestockInput{
		DrugID:      1,
		BatchNumber: "",
		Quantity:    20,
		ExpiredDate: "2027-10-31",
	})
	if err != domain.ErrInvalidBatchNumber {
		t.Errorf("expected ErrInvalidBatchNumber, got %v", err)
	}

	// Invalid date format
	err = uc.Restock(ctx, &domain.RestockInput{
		DrugID:      1,
		BatchNumber: "BCH-NEW",
		Quantity:    20,
		ExpiredDate: "31-10-2027",
	})
	if err != domain.ErrInvalidExpiredDate {
		t.Errorf("expected ErrInvalidExpiredDate, got %v", err)
	}
}

func TestExportAndImportStockOpnameCSV(t *testing.T) {
	uc, repo := setupPharmacyUsecase()
	ctx := context.Background()

	// 1. Export CSV
	csvStr, err := uc.ExportStockOpnameCSV(ctx, "KIDINEST")
	if err != nil {
		t.Fatalf("unexpected error exporting CSV: %v", err)
	}
	if !strings.Contains(csvStr, "Drug_ID,Nama_Obat") {
		t.Errorf("expected CSV header, got %s", csvStr)
	}
	if !strings.Contains(csvStr, "Paracetamol Drop 100mg/ml") {
		t.Errorf("expected Paracetamol Drop in CSV export, got %s", csvStr)
	}

	// 2. Import CSV with actual physical count = 48 (2 bottles damaged/broken)
	mockUploadedCSV := `Drug_ID,Nama_Obat,Unit,Kategori,Batch_ID,Nomor_Batch,Expired_Date,Stok_Sistem,Stok_Fisik_Nyata,Catatan_Selisih
1,Paracetamol Drop 100mg/ml,Botol,Antipiretik,10,BCH-PCM-01,2027-12-31,50,48,2 botol pecah di rak
`
	count, err := uc.ImportStockOpnameCSV(ctx, mockUploadedCSV, "Apt. Dimas")
	if err != nil {
		t.Fatalf("unexpected error importing CSV: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 item reconciled, got %d", count)
	}

	d := repo.drugs[1]
	if d.PhysicalStock != 48 {
		t.Errorf("expected reconciled physical stock 48, got %d", d.PhysicalStock)
	}
	if len(repo.opnameLogs) != 1 {
		t.Fatalf("expected 1 opname log entry, got %d", len(repo.opnameLogs))
	}
	logEntry := repo.opnameLogs[0]
	if logEntry.Discrepancy != -2 {
		t.Errorf("expected discrepancy -2, got %d", logEntry.Discrepancy)
	}
}
