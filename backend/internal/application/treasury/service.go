package treasury

import (
	"context"
	"errors"
	"fmt"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/shared"
	"clinic-api/internal/domain/treasury"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrShiftAlreadyOpen = errors.New("a cashier shift is already open for this cash account")
	ErrShiftNotFound    = errors.New("cashier shift not found")
	ErrShiftNotClosed   = errors.New("shift must be closed before depositing to bank")
	ErrShiftAlreadyDeposited = errors.New("shift has already been deposited")
	ErrInvalidDepositAmount  = errors.New("deposit amount must be positive")
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	uow               UnitOfWork
	treasuryRepo      treasury.Repository
	accountingService *accountingApp.Service
}

func NewService(
	uow UnitOfWork,
	treasuryRepo treasury.Repository,
	accountingService *accountingApp.Service,
) *Service {
	return &Service{
		uow:               uow,
		treasuryRepo:      treasuryRepo,
		accountingService: accountingService,
	}
}

type OpenShiftCommand struct {
	BranchID        uuid.UUID       `json:"branch_id"`
	CostCenterID    uuid.UUID       `json:"cost_center_id"`
	CashierUserID   *uuid.UUID      `json:"cashier_user_id,omitempty"`
	CashierName     string          `json:"cashier_name"`
	CashAccountCode string          `json:"cash_account_code"` // '1-1111' Kidinest, '1-1112' SSC
	OpeningCash     decimal.Decimal `json:"opening_cash"`
	Notes           string          `json:"notes,omitempty"`
}

func (s *Service) OpenShift(ctx context.Context, cmd OpenShiftCommand) (*treasury.CashierShift, error) {
	// Check if already open
	existing, err := s.treasuryRepo.GetActiveShift(ctx, cmd.BranchID, cmd.CashAccountCode)
	if err == nil && existing != nil {
		return nil, ErrShiftAlreadyOpen
	}

	shift := &treasury.CashierShift{
		BaseEntity:       shared.BaseEntity{ID: uuid.New()},
		BranchScoped:     shared.BranchScoped{BranchID: cmd.BranchID},
		CostCenterID:     cmd.CostCenterID,
		CashierUserID:    cmd.CashierUserID,
		CashierName:      cmd.CashierName,
		CashAccountCode:  cmd.CashAccountCode,
		StartTime:        time.Now(),
		OpeningCash:      cmd.OpeningCash,
		SystemCashIn:     decimal.Zero,
		SystemCashOut:    decimal.Zero,
		ExpectedCash:     cmd.OpeningCash,
		ActualCash:       decimal.Zero,
		DifferenceAmount: decimal.Zero,
		DifferenceType:   treasury.DifferenceNone,
		Status:           treasury.ShiftStatusOpen,
		Notes:            cmd.Notes,
	}

	if err := s.treasuryRepo.OpenShift(ctx, shift); err != nil {
		return nil, fmt.Errorf("open cashier shift: %w", err)
	}

	return shift, nil
}

type DenominationInput struct {
	DenominationValue int `json:"denomination_value"` // e.g. 100000, 50000, 20000, 10000, etc.
	Quantity          int `json:"quantity"`
}

type CloseShiftCommand struct {
	ShiftID       uuid.UUID           `json:"shift_id"`
	SystemCashIn  decimal.Decimal     `json:"system_cash_in"`
	SystemCashOut decimal.Decimal     `json:"system_cash_out"`
	Denominations []DenominationInput `json:"denominations"`
	Notes         string              `json:"notes,omitempty"`
}

func (s *Service) CloseShift(ctx context.Context, cmd CloseShiftCommand) (*treasury.CashierShift, *accountingDomain.JournalEntry, error) {
	shift, err := s.treasuryRepo.GetShiftByID(ctx, cmd.ShiftID)
	if err != nil {
		return nil, nil, ErrShiftNotFound
	}
	if shift.Status != treasury.ShiftStatusOpen {
		return nil, nil, errors.New("shift is not currently open")
	}

	// Compute actual cash from physical count denominations
	actualCash := decimal.Zero
	var denoms []treasury.ShiftDenomination
	for _, d := range cmd.Denominations {
		val := decimal.FromInt(int64(d.DenominationValue))
		qty := decimal.FromInt(int64(d.Quantity))
		subtotal := val.Mul(qty)
		actualCash = actualCash.Add(subtotal)

		denoms = append(denoms, treasury.ShiftDenomination{
			ID:                uuid.New(),
			ShiftID:           shift.ID,
			DenominationValue: d.DenominationValue,
			Quantity:          d.Quantity,
			Subtotal:          subtotal,
		})
	}

	// Calculate expected cash = opening + in - out
	expectedCash := shift.OpeningCash.Add(cmd.SystemCashIn).Sub(cmd.SystemCashOut)

	shift.SystemCashIn = cmd.SystemCashIn
	shift.SystemCashOut = cmd.SystemCashOut
	shift.ExpectedCash = expectedCash
	shift.ActualCash = actualCash
	now := time.Now()
	shift.EndTime = &now
	shift.Notes = cmd.Notes

	var journalEntry *accountingDomain.JournalEntry

	if actualCash.GreaterThan(expectedCash) {
		// Overage (Selisih Kas Lebih) -> 4-9200
		diff := actualCash.Sub(expectedCash)
		shift.DifferenceAmount = diff
		shift.DifferenceType = treasury.DifferenceOverage

		jLines := []accountingApp.CreateLineCommand{
			{
				AccountCode: shift.CashAccountCode,
				Debit:       diff,
				Credit:      decimal.Zero,
				Description: fmt.Sprintf("Penyesuaian Selisih Kas Lebih Shift Kasir %s", shift.CashierName),
			},
			{
				AccountCode: "4-9200", // Pendapatan Selisih Kas Lebih
				Debit:       decimal.Zero,
				Credit:      diff,
				Description: fmt.Sprintf("Pendapatan Selisih Kas Lebih Shift %s", shift.CashierName),
			},
		}

		err = s.uow.Do(ctx, func(txCtx context.Context) error {
			jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
				BranchID:        shift.BranchID,
				JournalType:     accountingDomain.JournalTypeAdjusting,
				TransactionDate: now,
				ReferenceType:   "SHIFT_RECONCILIATION",
				ReferenceID:     shift.ID.String(),
				Description:     fmt.Sprintf("Rekonsiliasi Kasir Shift Lebih %s", shift.CashierName),
				AutoPost:        true,
				Lines:           jLines,
			})
			if err != nil {
				return err
			}
			journalEntry = jEntry
			shift.ReconciliationJournalID = &jEntry.ID
			shift.Status = treasury.ShiftStatusClosed
			return s.treasuryRepo.CloseShift(txCtx, shift, denoms)
		})
	} else if actualCash.LessThan(expectedCash) {
		// Shortage (Selisih Kas Kurang) -> 6-1700
		diff := expectedCash.Sub(actualCash)
		shift.DifferenceAmount = diff
		shift.DifferenceType = treasury.DifferenceShortage

		jLines := []accountingApp.CreateLineCommand{
			{
				AccountCode: "6-1700", // Beban Selisih Kas Kurang
				Debit:       diff,
				Credit:      decimal.Zero,
				Description: fmt.Sprintf("Beban Selisih Kas Kurang Shift %s", shift.CashierName),
			},
			{
				AccountCode: shift.CashAccountCode,
				Debit:       decimal.Zero,
				Credit:      diff,
				Description: fmt.Sprintf("Penyesuaian Selisih Kas Kurang Shift Kasir %s", shift.CashierName),
			},
		}

		err = s.uow.Do(ctx, func(txCtx context.Context) error {
			jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
				BranchID:        shift.BranchID,
				JournalType:     accountingDomain.JournalTypeAdjusting,
				TransactionDate: now,
				ReferenceType:   "SHIFT_RECONCILIATION",
				ReferenceID:     shift.ID.String(),
				Description:     fmt.Sprintf("Rekonsiliasi Kasir Shift Kurang %s", shift.CashierName),
				AutoPost:        true,
				Lines:           jLines,
			})
			if err != nil {
				return err
			}
			journalEntry = jEntry
			shift.ReconciliationJournalID = &jEntry.ID
			shift.Status = treasury.ShiftStatusClosed
			return s.treasuryRepo.CloseShift(txCtx, shift, denoms)
		})
	} else {
		// Exactly matches
		shift.DifferenceAmount = decimal.Zero
		shift.DifferenceType = treasury.DifferenceNone
		shift.Status = treasury.ShiftStatusClosed

		err = s.uow.Do(ctx, func(txCtx context.Context) error {
			return s.treasuryRepo.CloseShift(txCtx, shift, denoms)
		})
	}

	if err != nil {
		return nil, nil, fmt.Errorf("close shift in uow: %w", err)
	}

	shift.Denominations = denoms
	return shift, journalEntry, nil
}

type DepositCommand struct {
	ShiftID        uuid.UUID       `json:"shift_id"`
	DepositAmount  decimal.Decimal `json:"deposit_amount"`
	BankAccCode    string          `json:"bank_account_code"` // Default '1-1121' Bank BCA
	DepositedBy    *uuid.UUID      `json:"deposited_by,omitempty"`
	DepositSlipRef string          `json:"deposit_slip_ref,omitempty"`
}

func (s *Service) DepositCashToBank(ctx context.Context, cmd DepositCommand) (*treasury.CashierShift, *accountingDomain.JournalEntry, error) {
	shift, err := s.treasuryRepo.GetShiftByID(ctx, cmd.ShiftID)
	if err != nil {
		return nil, nil, ErrShiftNotFound
	}
	if shift.Status == treasury.ShiftStatusOpen {
		return nil, nil, ErrShiftNotClosed
	}
	if shift.Status == treasury.ShiftStatusDeposited {
		return nil, nil, ErrShiftAlreadyDeposited
	}
	if !cmd.DepositAmount.IsPositive() {
		return nil, nil, ErrInvalidDepositAmount
	}

	bankAcc := cmd.BankAccCode
	if bankAcc == "" {
		bankAcc = "1-1121" // Bank BCA Operasional
	}

	now := time.Now()
	// Journal JB:
	// Debit: Bank BCA
	// Credit: Kas Kasir
	jLines := []accountingApp.CreateLineCommand{
		{
			AccountCode: bankAcc,
			Debit:       cmd.DepositAmount,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Setor Kas Kasir Shift %s ke Bank (%s)", shift.CashierName, cmd.DepositSlipRef),
		},
		{
			AccountCode: shift.CashAccountCode,
			Debit:       decimal.Zero,
			Credit:      cmd.DepositAmount,
			Description: fmt.Sprintf("Pengeluaran Kas Setoran Bank Shift %s", shift.CashierName),
		},
	}

	var journalEntry *accountingDomain.JournalEntry

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        shift.BranchID,
			JournalType:     accountingDomain.JournalTypeBank,
			TransactionDate: now,
			ReferenceType:   "CASHIER_BANK_DEPOSIT",
			ReferenceID:     shift.ID.String(),
			Description:     fmt.Sprintf("Setoran Kas Kasir ke %s - Slip %s", bankAcc, cmd.DepositSlipRef),
			CreatedBy:       cmd.DepositedBy,
			AutoPost:        true,
			Lines:           jLines,
		})
		if err != nil {
			return fmt.Errorf("create deposit journal: %w", err)
		}
		journalEntry = jEntry

		return s.treasuryRepo.RecordDeposit(txCtx, shift.ID, jEntry.ID)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("deposit to bank uow: %w", err)
	}

	shift.DepositJournalID = &journalEntry.ID
	shift.Status = treasury.ShiftStatusDeposited

	return shift, journalEntry, nil
}

func (s *Service) GetActiveShift(ctx context.Context, branchID uuid.UUID, cashAccountCode string) (*treasury.CashierShift, error) {
	return s.treasuryRepo.GetActiveShift(ctx, branchID, cashAccountCode)
}

func (s *Service) GetShiftByID(ctx context.Context, id uuid.UUID) (*treasury.CashierShift, error) {
	return s.treasuryRepo.GetShiftByID(ctx, id)
}

func (s *Service) ListShifts(ctx context.Context, filter treasury.ShiftFilter) ([]treasury.CashierShift, int64, error) {
	return s.treasuryRepo.ListShifts(ctx, filter)
}

