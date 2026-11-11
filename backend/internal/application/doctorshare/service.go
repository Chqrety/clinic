package doctorshare

import (
	"context"
	"errors"
	"fmt"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/doctorshare"
	"clinic-api/internal/domain/shared"
	"clinic-api/internal/domain/tax"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrSettlementNotFound     = errors.New("doctor fee settlement not found")
	ErrSettlementNotApproved = errors.New("settlement must be approved before payment")
	ErrSettlementAlreadyPaid = errors.New("settlement has already been paid")
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	uow               UnitOfWork
	doctorShareRepo   doctorshare.Repository
	taxRepo           tax.Repository
	accountingService *accountingApp.Service
}

func NewService(
	uow UnitOfWork,
	doctorShareRepo doctorshare.Repository,
	taxRepo tax.Repository,
	accountingService *accountingApp.Service,
) *Service {
	return &Service{
		uow:               uow,
		doctorShareRepo:   doctorShareRepo,
		taxRepo:           taxRepo,
		accountingService: accountingService,
	}
}

type CreateRuleCommand struct {
	BranchID          uuid.UUID       `json:"branch_id"`
	CostCenterID      uuid.UUID       `json:"cost_center_id"`
	DoctorID          *uuid.UUID      `json:"doctor_id,omitempty"`
	DoctorName        string          `json:"doctor_name"`
	ServiceCategory   string          `json:"service_category"`
	SharingPercentage decimal.Decimal `json:"sharing_percentage"`
	FixedFee          decimal.Decimal `json:"fixed_fee"`
	IsActive          bool            `json:"is_active"`
}

func (s *Service) CreateRule(ctx context.Context, cmd CreateRuleCommand) (*doctorshare.DoctorShareRule, error) {
	rule := &doctorshare.DoctorShareRule{
		ID:                uuid.New(),
		BranchID:          cmd.BranchID,
		CostCenterID:      cmd.CostCenterID,
		DoctorID:          cmd.DoctorID,
		DoctorName:        cmd.DoctorName,
		ServiceCategory:   cmd.ServiceCategory,
		SharingPercentage: cmd.SharingPercentage,
		FixedFee:          cmd.FixedFee,
		IsActive:          cmd.IsActive,
	}
	if err := s.doctorShareRepo.CreateRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("create doctor share rule: %w", err)
	}
	return rule, nil
}

func (s *Service) ListRules(ctx context.Context, filter doctorshare.RuleFilter) ([]doctorshare.DoctorShareRule, int64, error) {
	return s.doctorShareRepo.ListRules(ctx, filter)
}

type FeeItemInput struct {
	InvoiceID          uuid.UUID       `json:"invoice_id"`
	ServiceDate        time.Time       `json:"service_date"`
	PatientName        string          `json:"patient_name"`
	ServiceDescription string          `json:"service_description"`
	PatientFee         decimal.Decimal `json:"patient_fee"`
	DoctorShareAmount  decimal.Decimal `json:"doctor_share_amount"`
}

type GenerateSettlementCommand struct {
	BranchID    uuid.UUID      `json:"branch_id"`
	DoctorID    *uuid.UUID     `json:"doctor_id,omitempty"`
	DoctorName  string         `json:"doctor_name"`
	DoctorNPWP  string         `json:"doctor_npwp,omitempty"`
	TaxpayerNIK string         `json:"taxpayer_nik,omitempty"`
	PeriodYear  int            `json:"period_year"`
	PeriodMonth int            `json:"period_month"`
	Items       []FeeItemInput `json:"items"`
}

func (s *Service) GenerateSettlement(ctx context.Context, cmd GenerateSettlementCommand) (*doctorshare.DoctorFeeSettlement, error) {
	gross := decimal.Zero
	var items []doctorshare.DoctorFeeItem

	for _, it := range cmd.Items {
		gross = gross.Add(it.DoctorShareAmount)
		items = append(items, doctorshare.DoctorFeeItem{
			ID:                 uuid.New(),
			InvoiceID:          it.InvoiceID,
			ServiceDate:        it.ServiceDate,
			PatientName:        it.PatientName,
			ServiceDescription: it.ServiceDescription,
			PatientFee:         it.PatientFee,
			DoctorShareAmount:  it.DoctorShareAmount,
		})
	}

	// Calculate prior year cumulative DPP for progressive tax bracket calculation
	priorDPP := decimal.Zero
	if s.taxRepo != nil && cmd.TaxpayerNIK != "" {
		if cumDPP, err := s.taxRepo.GetCumulativeDPP(ctx, cmd.BranchID, cmd.TaxpayerNIK, cmd.PeriodYear); err == nil {
			priorDPP = cumDPP
		}
	}

	dpp, taxPPh21 := tax.CalculateDoctorPPh21(gross, priorDPP)
	netFee := gross.Sub(taxPPh21)

	settlement := &doctorshare.DoctorFeeSettlement{
		BaseEntity:   shared.BaseEntity{ID: uuid.New()},
		BranchScoped: shared.BranchScoped{BranchID: cmd.BranchID},
		DoctorID:     cmd.DoctorID,
		DoctorName:   cmd.DoctorName,
		DoctorNPWP:   cmd.DoctorNPWP,
		PeriodYear:   cmd.PeriodYear,
		PeriodMonth:  cmd.PeriodMonth,
		GrossFee:     gross,
		DPPAmount:    dpp,
		TaxPPh21:     taxPPh21,
		NetFee:       netFee,
		Status:       doctorshare.StatusDraft,
		Items:        items,
	}

	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.doctorShareRepo.CreateSettlement(txCtx, settlement)
	})
	if err != nil {
		return nil, fmt.Errorf("create settlement in uow: %w", err)
	}

	return settlement, nil
}

func (s *Service) ApproveSettlement(ctx context.Context, id uuid.UUID) (*doctorshare.DoctorFeeSettlement, error) {
	settlement, err := s.doctorShareRepo.GetSettlementByID(ctx, id)
	if err != nil {
		return nil, ErrSettlementNotFound
	}
	if settlement.Status == doctorshare.StatusPaid {
		return nil, ErrSettlementAlreadyPaid
	}

	err = s.doctorShareRepo.UpdateSettlementStatus(ctx, id, doctorshare.StatusApproved, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("approve settlement: %w", err)
	}
	settlement.Status = doctorshare.StatusApproved
	return settlement, nil
}

type PaySettlementCommand struct {
	SettlementID uuid.UUID  `json:"settlement_id"`
	BankAccCode  string     `json:"bank_account_code"` // Default '1-1121' Bank BCA Operasional
	UnitCode     string     `json:"unit_code"`         // 'KIDINEST' or 'SSC'
	PaidBy       *uuid.UUID `json:"paid_by,omitempty"`
}

func (s *Service) PaySettlement(ctx context.Context, cmd PaySettlementCommand) (*doctorshare.DoctorFeeSettlement, *accountingDomain.JournalEntry, error) {
	settlement, err := s.doctorShareRepo.GetSettlementByID(ctx, cmd.SettlementID)
	if err != nil {
		return nil, nil, ErrSettlementNotFound
	}
	if settlement.Status != doctorshare.StatusApproved {
		return nil, nil, ErrSettlementNotApproved
	}

	bankAcc := cmd.BankAccCode
	if bankAcc == "" {
		bankAcc = "1-1121" // Bank BCA Operasional
	}

	payableAcc := "2-1210" // Hutang Jasa Medis Dokter Kidinest
	if cmd.UnitCode == "SSC" {
		payableAcc = "2-1220" // Hutang Jasa Medis Dokter SSC
	}

	taxAcc := "2-1501.02" // Hutang PPh 21 Dokter / Tenaga Medis Ahli

	// Double-entry payment journal:
	// Debit: Hutang Jasa Medis (GrossFee)
	// Credit: Bank BCA (NetFee)
	// Credit: Hutang PPh 21 Dokter (TaxPPh21)
	journalLines := []accountingApp.CreateLineCommand{
		{
			AccountCode: payableAcc,
			Debit:       settlement.GrossFee,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Pembayaran Jasa Medis %s Periode %d-%02d", settlement.DoctorName, settlement.PeriodYear, settlement.PeriodMonth),
		},
		{
			AccountCode: bankAcc,
			Debit:       decimal.Zero,
			Credit:      settlement.NetFee,
			Description: fmt.Sprintf("Transfer Bersih Jasa Medis %s via %s", settlement.DoctorName, bankAcc),
		},
	}

	if settlement.TaxPPh21.IsPositive() {
		journalLines = append(journalLines, accountingApp.CreateLineCommand{
			AccountCode: taxAcc,
			Debit:       decimal.Zero,
			Credit:      settlement.TaxPPh21,
			Description: fmt.Sprintf("Potongan PPh 21 Tenaga Medis Ahli %s (DPP 50%%)", settlement.DoctorName),
		})
	}

	var journalEntry *accountingDomain.JournalEntry
	now := time.Now()

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		// 1. Post JKK (Cash Disbursement Journal)
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        settlement.BranchID,
			JournalType:     accountingDomain.JournalTypeCashDisbursement,
			TransactionDate: now,
			ReferenceType:   "DOCTOR_FEE_SETTLEMENT",
			ReferenceID:     settlement.ID.String(),
			Description:     fmt.Sprintf("Pelunasan Fee Dokter %s Bulan %d-%02d", settlement.DoctorName, settlement.PeriodYear, settlement.PeriodMonth),
			CreatedBy:       cmd.PaidBy,
			AutoPost:        true,
			Lines:           journalLines,
		})
		if err != nil {
			return fmt.Errorf("post doctor settlement journal: %w", err)
		}
		journalEntry = jEntry

		// 2. Register Tax Withholding entry if tax is positive
		if settlement.TaxPPh21.IsPositive() && s.taxRepo != nil {
			bupotNum := fmt.Sprintf("BP-21-%d%02d-%04d", settlement.PeriodYear, settlement.PeriodMonth, time.Now().UnixNano()%10000)
			tw := &tax.TaxWithholding{
				BaseEntity:   shared.BaseEntity{ID: uuid.New()},
				BranchScoped: shared.BranchScoped{BranchID: settlement.BranchID},
				TaxType:      tax.TaxTypePPh21Doctor,
				BupotNumber:  bupotNum,
				TaxpayerNPWP: settlement.DoctorNPWP,
				TaxpayerName: settlement.DoctorName,
				ObjectCode:   "21-100-07", // Tenaga Ahli Dokter DJP
				GrossAmount:  settlement.GrossFee,
				DPPAmount:    settlement.DPPAmount,
				TaxRate:      decimal.FromFloat(5.00), // Progressive base
				TaxAmount:    settlement.TaxPPh21,
				PeriodYear:   settlement.PeriodYear,
				PeriodMonth:  settlement.PeriodMonth,
				IsReported:   false,
				JournalID:    &jEntry.ID,
			}
			if err := s.taxRepo.CreateWithholding(txCtx, tw); err != nil {
				return fmt.Errorf("create tax withholding record: %w", err)
			}
		}

		// 3. Update settlement status to PAID
		return s.doctorShareRepo.UpdateSettlementStatus(txCtx, settlement.ID, doctorshare.StatusPaid, &now, &jEntry.ID)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("pay settlement uow: %w", err)
	}

	settlement.Status = doctorshare.StatusPaid
	settlement.PaidAt = &now
	settlement.JournalID = &journalEntry.ID

	return settlement, journalEntry, nil
}

func (s *Service) GetSettlementByID(ctx context.Context, id uuid.UUID) (*doctorshare.DoctorFeeSettlement, error) {
	return s.doctorShareRepo.GetSettlementByID(ctx, id)
}

func (s *Service) ListSettlements(ctx context.Context, filter doctorshare.SettlementFilter) ([]doctorshare.DoctorFeeSettlement, int64, error) {
	return s.doctorShareRepo.ListSettlements(ctx, filter)
}
