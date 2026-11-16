package prepaid

import (
	"context"
	"errors"
	"fmt"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/prepaid"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrPrepaidNotFound       = errors.New("prepaid expense not found")
	ErrAlreadyFullyAmortized = errors.New("prepaid expense is already fully amortized")
	ErrInvalidPeriod         = errors.New("end date must be after start date")
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	uow               UnitOfWork
	prepaidRepo       prepaid.Repository
	accountingService *accountingApp.Service
}

func NewService(
	uow UnitOfWork,
	prepaidRepo prepaid.Repository,
	accountingService *accountingApp.Service,
) *Service {
	return &Service{
		uow:               uow,
		prepaidRepo:       prepaidRepo,
		accountingService: accountingService,
	}
}

type RegisterPrepaidCommand struct {
	BranchID           uuid.UUID               `json:"branch_id"`
	CostCenterID       uuid.UUID               `json:"cost_center_id"`
	Name               string                  `json:"name"`
	Category           prepaid.PrepaidCategory `json:"category"` // 'RENT', 'INSURANCE', 'SOFTWARE_LICENSE', 'MAINTENANCE'
	TotalAmount        decimal.Decimal         `json:"total_amount"`
	PrepaidAccountCode string                  `json:"prepaid_account_code"` // '1-1400'
	ExpenseAccountCode string                  `json:"expense_account_code"` // '6-1200'
	StartDate          time.Time               `json:"start_date"`
	EndDate            time.Time               `json:"end_date"`
	PaymentAccountCode string                  `json:"payment_account_code,omitempty"` // e.g. '1-1121' Bank BCA
	AutoPostPayment    bool                    `json:"auto_post_payment"`
}

func (s *Service) RegisterPrepaid(ctx context.Context, cmd RegisterPrepaidCommand) (*prepaid.PrepaidExpense, *accountingDomain.JournalEntry, error) {
	if cmd.EndDate.Before(cmd.StartDate) || cmd.EndDate.Equal(cmd.StartDate) {
		return nil, nil, ErrInvalidPeriod
	}

	months := (cmd.EndDate.Year()-cmd.StartDate.Year())*12 + int(cmd.EndDate.Month()-cmd.StartDate.Month())
	if months <= 0 {
		months = 1
	}

	monthly := cmd.TotalAmount.Div(decimal.FromInt(int64(months)))

	prepaidAcc := cmd.PrepaidAccountCode
	if prepaidAcc == "" {
		prepaidAcc = "1-1400" // Biaya Dibayar di Muka
	}
	expenseAcc := cmd.ExpenseAccountCode
	if expenseAcc == "" {
		expenseAcc = "6-1200" // Beban Sewa Gedung & Pemeliharaan
	}

	p := &prepaid.PrepaidExpense{
		BaseEntity:         shared.BaseEntity{ID: uuid.New()},
		BranchScoped:       shared.BranchScoped{BranchID: cmd.BranchID},
		CostCenterID:       cmd.CostCenterID,
		Name:               cmd.Name,
		Category:           cmd.Category,
		TotalAmount:        cmd.TotalAmount,
		PrepaidAccountCode: prepaidAcc,
		ExpenseAccountCode: expenseAcc,
		StartDate:          cmd.StartDate,
		EndDate:            cmd.EndDate,
		PeriodMonths:       months,
		MonthlyAmount:      monthly,
		RemainingAmount:    cmd.TotalAmount,
		Status:             prepaid.StatusActive,
	}

	var journalEntry *accountingDomain.JournalEntry

	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		if err := s.prepaidRepo.CreatePrepaid(txCtx, p); err != nil {
			return err
		}

		if cmd.AutoPostPayment && cmd.PaymentAccountCode != "" {
			jLines := []accountingApp.CreateLineCommand{
				{
					AccountCode: prepaidAcc,
					Debit:       cmd.TotalAmount,
					Credit:      decimal.Zero,
					Description: fmt.Sprintf("Pembayaran Dibayar di Muka: %s", p.Name),
				},
				{
					AccountCode: cmd.PaymentAccountCode,
					Debit:       decimal.Zero,
					Credit:      cmd.TotalAmount,
					Description: fmt.Sprintf("Pengeluaran Kas/Bank untuk %s", p.Name),
				},
			}

			jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
				BranchID:        cmd.BranchID,
				JournalType:     accountingDomain.JournalTypeCashDisbursement,
				TransactionDate: cmd.StartDate,
				ReferenceType:   "PREPAID_PAYMENT",
				ReferenceID:     p.ID.String(),
				Description:     fmt.Sprintf("Pencatatan Biaya Dibayar di Muka %s", p.Name),
				AutoPost:        true,
				Lines:           jLines,
			})
			if err != nil {
				return fmt.Errorf("create prepaid payment journal: %w", err)
			}
			journalEntry = jEntry
			p.JournalID = &jEntry.ID
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("register prepaid in uow: %w", err)
	}

	return p, journalEntry, nil
}

type AmortizeCommand struct {
	PrepaidID   uuid.UUID `json:"prepaid_id"`
	PeriodYear  int       `json:"period_year"`
	PeriodMonth int       `json:"period_month"`
}

func (s *Service) AmortizeMonthly(ctx context.Context, cmd AmortizeCommand) (*prepaid.PrepaidAmortization, *accountingDomain.JournalEntry, error) {
	p, err := s.prepaidRepo.GetPrepaidByID(ctx, cmd.PrepaidID)
	if err != nil {
		return nil, nil, ErrPrepaidNotFound
	}
	if p.Status != prepaid.StatusActive {
		return nil, nil, ErrAlreadyFullyAmortized
	}

	amount := p.MonthlyAmount
	if amount.GreaterThan(p.RemainingAmount) {
		amount = p.RemainingAmount
	}

	newRemaining := p.RemainingAmount.Sub(amount)
	newStatus := prepaid.StatusActive
	if newRemaining.IsZero() || newRemaining.IsNegative() {
		newRemaining = decimal.Zero
		newStatus = prepaid.StatusFullyAmortized
	}

	now := time.Now()
	jLines := []accountingApp.CreateLineCommand{
		{
			AccountCode: p.ExpenseAccountCode,
			Debit:       amount,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Amortisasi Beban %s (%d-%02d)", p.Name, cmd.PeriodYear, cmd.PeriodMonth),
		},
		{
			AccountCode: p.PrepaidAccountCode,
			Debit:       decimal.Zero,
			Credit:      amount,
			Description: fmt.Sprintf("Pengurangan Biaya Dibayar di Muka %s (%d-%02d)", p.Name, cmd.PeriodYear, cmd.PeriodMonth),
		},
	}

	amort := &prepaid.PrepaidAmortization{
		ID:               uuid.New(),
		PrepaidExpenseID: p.ID,
		PeriodYear:       cmd.PeriodYear,
		PeriodMonth:      cmd.PeriodMonth,
		Amount:           amount,
		CreatedAt:        now,
	}

	var journalEntry *accountingDomain.JournalEntry

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        p.BranchID,
			JournalType:     accountingDomain.JournalTypeAdjusting,
			TransactionDate: now,
			ReferenceType:   "PREPAID_AMORTIZATION",
			ReferenceID:     fmt.Sprintf("%s-%d%02d", p.ID.String(), cmd.PeriodYear, cmd.PeriodMonth),
			Description:     fmt.Sprintf("Amortisasi Bulanan %s Periode %d-%02d", p.Name, cmd.PeriodYear, cmd.PeriodMonth),
			AutoPost:        true,
			Lines:           jLines,
		})
		if err != nil {
			return fmt.Errorf("create amortization journal: %w", err)
		}
		journalEntry = jEntry
		amort.JournalID = &jEntry.ID

		return s.prepaidRepo.RecordAmortization(txCtx, amort, newRemaining, newStatus)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("amortize monthly in uow: %w", err)
	}

	return amort, journalEntry, nil
}

func (s *Service) RunMonthlyAmortizationBatch(ctx context.Context, branchID *uuid.UUID, year, month int) ([]prepaid.PrepaidAmortization, error) {
	active, err := s.prepaidRepo.ListActivePrepaid(ctx, branchID)
	if err != nil {
		return nil, fmt.Errorf("list active prepaid: %w", err)
	}

	var results []prepaid.PrepaidAmortization
	for _, p := range active {
		amort, _, err := s.AmortizeMonthly(ctx, AmortizeCommand{
			PrepaidID:   p.ID,
			PeriodYear:  year,
			PeriodMonth: month,
		})
		if err == nil && amort != nil {
			results = append(results, *amort)
		}
	}
	return results, nil
}

func (s *Service) GetPrepaidByID(ctx context.Context, id uuid.UUID) (*prepaid.PrepaidExpense, error) {
	return s.prepaidRepo.GetPrepaidByID(ctx, id)
}

func (s *Service) ListActivePrepaid(ctx context.Context, branchID *uuid.UUID) ([]prepaid.PrepaidExpense, error) {
	return s.prepaidRepo.ListActivePrepaid(ctx, branchID)
}

func (s *Service) ListPrepaid(ctx context.Context, filter prepaid.PrepaidFilter) ([]prepaid.PrepaidExpense, int64, error) {
	return s.prepaidRepo.List(ctx, filter)
}

func (s *Service) GetAmortizationHistory(ctx context.Context, id uuid.UUID) ([]prepaid.PrepaidAmortization, error) {
	return s.prepaidRepo.GetAmortizationHistory(ctx, id)
}

