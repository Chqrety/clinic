package tax

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/internal/domain/tax"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type Service struct {
	taxRepo tax.Repository
}

func NewService(taxRepo tax.Repository) *Service {
	return &Service{
		taxRepo: taxRepo,
	}
}

type CreateWithholdingCommand struct {
	BranchID     uuid.UUID       `json:"branch_id"`
	TaxType      tax.TaxType     `json:"tax_type"`
	TaxpayerNIK  string          `json:"taxpayer_nik"`
	TaxpayerNPWP string          `json:"taxpayer_npwp"`
	TaxpayerName string          `json:"taxpayer_name"`
	ObjectCode   string          `json:"object_code"`
	GrossAmount  decimal.Decimal `json:"gross_amount"`
	DPPAmount    decimal.Decimal `json:"dpp_amount"`
	TaxRate      decimal.Decimal `json:"tax_rate"`
	TaxAmount    decimal.Decimal `json:"tax_amount"`
	PeriodYear   int             `json:"period_year"`
	PeriodMonth  int             `json:"period_month"`
	JournalID    *uuid.UUID      `json:"journal_id,omitempty"`
}

func (s *Service) CreateWithholding(ctx context.Context, cmd CreateWithholdingCommand) (*tax.TaxWithholding, error) {
	bupotNumber := fmt.Sprintf("BP-%s-%d%02d-%04d", cmd.TaxType, cmd.PeriodYear, cmd.PeriodMonth, time.Now().UnixNano()%10000)

	tw := &tax.TaxWithholding{
		BaseEntity:   shared.BaseEntity{ID: uuid.New()},
		BranchScoped: shared.BranchScoped{BranchID: cmd.BranchID},
		TaxType:      cmd.TaxType,
		BupotNumber:  bupotNumber,
		TaxpayerNIK:  cmd.TaxpayerNIK,
		TaxpayerNPWP: cmd.TaxpayerNPWP,
		TaxpayerName: cmd.TaxpayerName,
		ObjectCode:   cmd.ObjectCode,
		GrossAmount:  cmd.GrossAmount,
		DPPAmount:    cmd.DPPAmount,
		TaxRate:      cmd.TaxRate,
		TaxAmount:    cmd.TaxAmount,
		PeriodYear:   cmd.PeriodYear,
		PeriodMonth:  cmd.PeriodMonth,
		IsReported:   false,
		JournalID:    cmd.JournalID,
	}

	if err := s.taxRepo.CreateWithholding(ctx, tw); err != nil {
		return nil, fmt.Errorf("create tax withholding: %w", err)
	}

	return tw, nil
}

func (s *Service) ListWithholdings(ctx context.Context, branchID *uuid.UUID, year, month int, taxType *tax.TaxType) ([]tax.TaxWithholding, error) {
	return s.taxRepo.ListWithholdings(ctx, branchID, year, month, taxType)
}

func (s *Service) MarkAsReported(ctx context.Context, id uuid.UUID) error {
	return s.taxRepo.MarkAsReported(ctx, id)
}

type TaxSummaryItem struct {
	TaxType     tax.TaxType     `json:"tax_type"`
	Count       int             `json:"count"`
	TotalGross  decimal.Decimal `json:"total_gross"`
	TotalDPP    decimal.Decimal `json:"total_dpp"`
	TotalTax    decimal.Decimal `json:"total_tax"`
}

func (s *Service) GetTaxSummary(ctx context.Context, branchID *uuid.UUID, year, month int) ([]TaxSummaryItem, error) {
	list, err := s.taxRepo.ListWithholdings(ctx, branchID, year, month, nil)
	if err != nil {
		return nil, fmt.Errorf("list withholdings for summary: %w", err)
	}

	summaryMap := make(map[tax.TaxType]*TaxSummaryItem)

	for _, item := range list {
		summary, exists := summaryMap[item.TaxType]
		if !exists {
			summary = &TaxSummaryItem{
				TaxType:    item.TaxType,
				Count:      0,
				TotalGross: decimal.Zero,
				TotalDPP:   decimal.Zero,
				TotalTax:   decimal.Zero,
			}
			summaryMap[item.TaxType] = summary
		}
		summary.Count++
		summary.TotalGross = summary.TotalGross.Add(item.GrossAmount)
		summary.TotalDPP = summary.TotalDPP.Add(item.DPPAmount)
		summary.TotalTax = summary.TotalTax.Add(item.TaxAmount)
	}

	var result []TaxSummaryItem
	for _, v := range summaryMap {
		result = append(result, *v)
	}
	return result, nil
}

// ExportEBupot21CSV exports withholdings in the standard CSV format required by DJP e-Bupot 21/26.
func (s *Service) ExportEBupot21CSV(ctx context.Context, branchID *uuid.UUID, year, month int) ([]byte, error) {
	p21 := tax.TaxTypePPh21Doctor
	records, err := s.taxRepo.ListWithholdings(ctx, branchID, year, month, &p21)
	if err != nil {
		return nil, fmt.Errorf("fetch pph21 withholdings: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// DJP e-Bupot 21/26 Header
	header := []string{
		"MasaPajak", "TahunPajak", "Pembetulan", "NomorBuktiPotong", "TanggalBuktiPotong",
		"NPWP", "NIK", "NamaWajibPajak", "KodeObjekPajak",
		"JumlahPenghasilanBruto", "JumlahDPP", "Tarif", "JumlahPPh",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	nowStr := time.Now().Format("2006-01-02")
	for _, r := range records {
		npwp := r.TaxpayerNPWP
		if npwp == "" {
			npwp = "0000000000000000"
		}
		row := []string{
			fmt.Sprintf("%02d", r.PeriodMonth),
			fmt.Sprintf("%d", r.PeriodYear),
			"0", // Normal (Bukan Pembetulan)
			r.BupotNumber,
			nowStr,
			npwp,
			r.TaxpayerNIK,
			r.TaxpayerName,
			r.ObjectCode,
			r.GrossAmount.String(),
			r.DPPAmount.String(),
			r.TaxRate.String(),
			r.TaxAmount.String(),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("csv writer error: %w", err)
	}

	return buf.Bytes(), nil
}

