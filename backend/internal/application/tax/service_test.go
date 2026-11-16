package tax_test

import (
	"context"
	"strings"
	"testing"

	taxApp "clinic-api/internal/application/tax"
	"clinic-api/internal/domain/tax"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockTaxRepo struct {
	withholdings []tax.TaxWithholding
}

func (m *mockTaxRepo) CreateWithholding(ctx context.Context, tw *tax.TaxWithholding) error {
	m.withholdings = append(m.withholdings, *tw)
	return nil
}

func (m *mockTaxRepo) GetWithholdingByID(ctx context.Context, id uuid.UUID) (*tax.TaxWithholding, error) {
	for _, w := range m.withholdings {
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, nil
}

func (m *mockTaxRepo) GetCumulativeDPP(ctx context.Context, branchID uuid.UUID, taxpayerNIK string, year int) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

func (m *mockTaxRepo) ListWithholdings(ctx context.Context, filter tax.TaxFilter) ([]tax.TaxWithholding, int64, error) {
	var list []tax.TaxWithholding
	for _, w := range m.withholdings {
		if filter.TaxType != nil && w.TaxType != *filter.TaxType {
			continue
		}
		list = append(list, w)
	}
	return list, int64(len(list)), nil
}

func (m *mockTaxRepo) MarkAsReported(ctx context.Context, id uuid.UUID) error {
	for i := range m.withholdings {
		if m.withholdings[i].ID == id {
			m.withholdings[i].IsReported = true
		}
	}
	return nil
}

func TestTaxService_CreateAndExportEBupot21(t *testing.T) {
	repo := &mockTaxRepo{}
	svc := taxApp.NewService(repo)

	ctx := context.Background()
	branchID := uuid.New()

	// 1. Create Tax Withholding (PPh 21 Dokter Tenaga Ahli)
	tw, err := svc.CreateWithholding(ctx, taxApp.CreateWithholdingCommand{
		BranchID:     branchID,
		TaxType:      tax.TaxTypePPh21Doctor,
		TaxpayerNIK:  "3374012345678001",
		TaxpayerNPWP: "09.123.456.7-001.000",
		TaxpayerName: "dr. Kartika Sp.DVE",
		ObjectCode:   "21-100-07",
		GrossAmount:  decimal.FromInt(15000000),
		DPPAmount:    decimal.FromInt(7500000),
		TaxRate:      decimal.FromFloat(5.00),
		TaxAmount:    decimal.FromInt(375000),
		PeriodYear:   2026,
		PeriodMonth:  11,
	})
	if err != nil {
		t.Fatalf("failed to create tax withholding: %v", err)
	}

	if tw.BupotNumber == "" {
		t.Errorf("expected bupot number to be populated")
	}

	// 2. Export DJP e-Bupot 21/26 CSV
	csvBytes, err := svc.ExportEBupot21CSV(ctx, &branchID, 2026, 11)
	if err != nil {
		t.Fatalf("failed to export ebupot csv: %v", err)
	}

	csvStr := string(csvBytes)
	if !strings.Contains(csvStr, "MasaPajak,TahunPajak") {
		t.Errorf("expected CSV header, got %s", csvStr)
	}
	if !strings.Contains(csvStr, "dr. Kartika Sp.DVE") {
		t.Errorf("expected taxpayer name in CSV, got %s", csvStr)
	}
	if !strings.Contains(csvStr, "21-100-07") {
		t.Errorf("expected object code in CSV, got %s", csvStr)
	}

	// 3. Verify Tax Summary
	summary, err := svc.GetTaxSummary(ctx, &branchID, 2026, 11)
	if err != nil {
		t.Fatalf("failed to get tax summary: %v", err)
	}
	if len(summary) != 1 {
		t.Fatalf("expected 1 summary category, got %d", len(summary))
	}
	if !summary[0].TotalTax.Equal(decimal.FromInt(375000)) {
		t.Errorf("expected total tax 375000, got %s", summary[0].TotalTax)
	}
}

