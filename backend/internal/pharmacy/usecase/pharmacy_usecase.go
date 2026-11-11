package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"clinic-api/internal/pharmacy/domain"
)

type PharmacyUsecase struct {
	repo domain.PharmacyRepository
}

func NewPharmacyUsecase(repo domain.PharmacyRepository) *PharmacyUsecase {
	return &PharmacyUsecase{repo: repo}
}

func (uc *PharmacyUsecase) ListIncomingPrescriptions(ctx context.Context, clinicUnit string) ([]*domain.PharmacyPrescription, error) {
	return uc.repo.ListIncomingPrescriptions(ctx, clinicUnit)
}

func (uc *PharmacyUsecase) GetPrescription(ctx context.Context, id int64) (*domain.PharmacyPrescription, error) {
	return uc.repo.GetPrescriptionByID(ctx, id)
}

func (uc *PharmacyUsecase) StartPreparing(ctx context.Context, id int64) error {
	rx, err := uc.repo.GetPrescriptionByID(ctx, id)
	if err != nil {
		return err
	}
	if rx.Status == "DISPENSED" {
		return domain.ErrAlreadyDispensed
	}
	if rx.Status == "CANCELLED" {
		return domain.ErrAlreadyCancelled
	}
	return uc.repo.UpdatePrescriptionStatus(ctx, id, "PREPARING")
}

func (uc *PharmacyUsecase) Dispense(ctx context.Context, id int64) error {
	return uc.repo.DispensePrescriptionTx(ctx, id)
}

func (uc *PharmacyUsecase) Cancel(ctx context.Context, id int64) error {
	return uc.repo.CancelPrescriptionTx(ctx, id)
}

func (uc *PharmacyUsecase) ListDrugs(ctx context.Context, clinicUnit string) ([]*domain.Drug, error) {
	return uc.repo.ListDrugs(ctx, clinicUnit)
}

func (uc *PharmacyUsecase) CreateDrug(ctx context.Context, in *domain.CreateDrugInput) (*domain.Drug, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Unit = strings.TrimSpace(in.Unit)
	in.Category = strings.TrimSpace(in.Category)

	if in.Name == "" {
		return nil, domain.ErrEmptyDrugName
	}
	if in.Unit == "" {
		in.Unit = "Botol"
	}
	if in.Category == "" {
		in.Category = "Umum"
	}

	return uc.repo.CreateDrug(ctx, in)
}

func (uc *PharmacyUsecase) Restock(ctx context.Context, in *domain.RestockInput) error {
	in.BatchNumber = strings.TrimSpace(in.BatchNumber)
	in.ExpiredDate = strings.TrimSpace(in.ExpiredDate)

	if in.DrugID <= 0 {
		return domain.ErrDrugNotFound
	}
	if in.BatchNumber == "" {
		return domain.ErrInvalidBatchNumber
	}
	if in.Quantity <= 0 {
		return domain.ErrInvalidQuantity
	}
	if _, err := time.Parse("2006-01-02", in.ExpiredDate); err != nil {
		return domain.ErrInvalidExpiredDate
	}

	return uc.repo.RestockBatch(ctx, in)
}

// ExportStockOpnameCSV: generates spreadsheet template for physical audit
func (uc *PharmacyUsecase) ExportStockOpnameCSV(ctx context.Context, clinicUnit string) (string, error) {
	drugs, err := uc.repo.ListDrugs(ctx, clinicUnit)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// CSV Header
	header := []string{"Drug_ID", "Nama_Obat", "Unit", "Kategori", "Batch_ID", "Nomor_Batch", "Expired_Date", "Stok_Sistem", "Stok_Fisik_Nyata", "Catatan_Selisih"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	for _, d := range drugs {
		if len(d.Batches) == 0 {
			// Row without batch
			row := []string{
				strconv.FormatInt(d.ID, 10),
				d.Name,
				d.Unit,
				d.Category,
				"0",
				"-",
				"-",
				strconv.Itoa(d.PhysicalStock),
				strconv.Itoa(d.PhysicalStock), // default suggested
				"",
			}
			_ = writer.Write(row)
		} else {
			for _, b := range d.Batches {
				row := []string{
					strconv.FormatInt(d.ID, 10),
					d.Name,
					d.Unit,
					d.Category,
					strconv.FormatInt(b.ID, 10),
					b.BatchNumber,
					b.ExpiredDate,
					strconv.Itoa(b.PhysicalStock),
					strconv.Itoa(b.PhysicalStock),
					"",
				}
				_ = writer.Write(row)
			}
		}
	}

	writer.Flush()
	return buf.String(), nil
}

// ImportStockOpnameCSV: parses spreadsheet and reconciles physical stock
func (uc *PharmacyUsecase) ImportStockOpnameCSV(ctx context.Context, csvContent string, adjustedBy string) (int, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("gagal membaca format CSV: %w", err)
	}

	if len(records) <= 1 {
		return 0, domain.ErrEmptyOpnameItems
	}

	if adjustedBy == "" {
		adjustedBy = "Apoteker Bertugas"
	}

	var items []domain.StockOpnameItem

	// Skip header (row 0)
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < 8 {
			continue
		}

		drugID, err := strconv.ParseInt(strings.TrimSpace(row[0]), 10, 64)
		if err != nil || drugID <= 0 {
			continue
		}

		drugName := strings.TrimSpace(row[1])
		batchID, _ := strconv.ParseInt(strings.TrimSpace(row[4]), 10, 64)
		batchNum := strings.TrimSpace(row[5])
		expDate := strings.TrimSpace(row[6])
		systemStock, _ := strconv.Atoi(strings.TrimSpace(row[7]))

		actualStockStr := ""
		if len(row) >= 9 {
			actualStockStr = strings.TrimSpace(row[8])
		}
		actualStock, err := strconv.Atoi(actualStockStr)
		if err != nil {
			actualStock = systemStock
		}

		notes := ""
		if len(row) >= 10 {
			notes = strings.TrimSpace(row[9])
		}

		discrepancy := actualStock - systemStock

		items = append(items, domain.StockOpnameItem{
			DrugID:      drugID,
			DrugName:    drugName,
			BatchID:     batchID,
			BatchNumber: batchNum,
			ExpiredDate: expDate,
			SystemStock: systemStock,
			ActualStock: actualStock,
			Discrepancy: discrepancy,
			Notes:       notes,
		})
	}

	if len(items) == 0 {
		return 0, domain.ErrEmptyOpnameItems
	}

	if err := uc.repo.ApplyStockOpnameTx(ctx, items, adjustedBy); err != nil {
		return 0, fmt.Errorf("gagal eksekusi rekonsiliasi stok opnam: %w", err)
	}

	return len(items), nil
}

func (uc *PharmacyUsecase) ListOpnameLogs(ctx context.Context, limit int) ([]*domain.StockOpnameLog, error) {
	return uc.repo.ListOpnameLogs(ctx, limit)
}
