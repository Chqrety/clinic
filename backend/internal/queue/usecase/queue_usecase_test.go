package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"clinic-api/internal/queue/domain"
)

type mockPatientRepo struct {
	patients map[string]*domain.Patient
	nextID   int64
}

func newMockPatientRepo() *mockPatientRepo {
	return &mockPatientRepo{
		patients: make(map[string]*domain.Patient),
		nextID:   1,
	}
}

func (m *mockPatientRepo) FindOrCreateByPhone(ctx context.Context, name, phone string) (*domain.Patient, error) {
	if p, ok := m.patients[phone]; ok {
		return p, nil
	}
	p := &domain.Patient{
		ID:        m.nextID,
		Name:      name,
		Phone:     phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.patients[phone] = p
	m.nextID++
	return p, nil
}

func (m *mockPatientRepo) GetPatientByID(ctx context.Context, id int64) (*domain.Patient, error) {
	for _, p := range m.patients {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, domain.ErrPatientNotFound
}

type mockQueueRepo struct {
	queues  map[int64]*domain.QueueItem
	counters map[string]int
	nextID  int64
}

func newMockQueueRepo() *mockQueueRepo {
	return &mockQueueRepo{
		queues:  make(map[int64]*domain.QueueItem),
		counters: make(map[string]int),
		nextID:  1,
	}
}

func (m *mockQueueRepo) GenerateNextQueueNumber(ctx context.Context, date time.Time, unit domain.ClinicUnit, category domain.QueueCategory) (string, error) {
	prefix := "A"
	if unit == domain.UnitSSC {
		prefix = "S"
	} else if category == domain.CategorySakit {
		prefix = "B"
	}
	key := string(unit) + "-" + prefix
	m.counters[key]++
	return fmt.Sprintf("%s-%03d", prefix, m.counters[key]), nil
}

func (m *mockQueueRepo) Create(ctx context.Context, q *domain.QueueItem) error {
	q.ID = m.nextID
	q.CreatedAt = time.Now()
	m.queues[q.ID] = q
	m.nextID++
	return nil
}

func (m *mockQueueRepo) GetByID(ctx context.Context, id int64) (*domain.QueueItem, error) {
	q, ok := m.queues[id]
	if !ok {
		return nil, domain.ErrQueueNotFound
	}
	return q, nil
}

func (m *mockQueueRepo) UpdatePaymentSuccess(ctx context.Context, id int64, queueNumber string) (*domain.QueueItem, error) {
	q, ok := m.queues[id]
	if !ok {
		return nil, domain.ErrQueueNotFound
	}
	now := time.Now()
	q.PaymentStatus = domain.PaymentPaid
	q.QueueStatus = domain.QueueWaiting
	q.QueueNumber = queueNumber
	q.PaidAt = &now
	return q, nil
}

func (m *mockQueueRepo) UpdateQueueStatus(ctx context.Context, id int64, status domain.QueueStatus) error {
	q, ok := m.queues[id]
	if !ok {
		return domain.ErrQueueNotFound
	}
	q.QueueStatus = status
	return nil
}

func (m *mockQueueRepo) GetCallingQueue(ctx context.Context, unit string) (*domain.QueueItem, error) {
	for _, q := range m.queues {
		if q.QueueStatus == domain.QueueCalling {
			return q, nil
		}
	}
	return nil, nil
}

func (m *mockQueueRepo) ListUpcomingQueues(ctx context.Context, unit string, limit int) ([]domain.QueueItem, error) {
	var list []domain.QueueItem
	for _, q := range m.queues {
		if q.QueueStatus == domain.QueueWaiting {
			list = append(list, *q)
		}
	}
	return list, nil
}

func TestQueueUsecase_Lifecycle(t *testing.T) {
	pRepo := newMockPatientRepo()
	qRepo := newMockQueueRepo()
	wa := NewMockWhatsAppAdapter()
	uc := NewQueueUsecase(pRepo, qRepo, wa, "http://localhost:5173")

	t.Run("Pendaftaran Kidinest harus berstatus PENDING dan belum dapat nomor antrian resmi", func(t *testing.T) {
		res, err := uc.RegisterQueue(context.Background(), domain.RegisterInput{
			ClinicUnit:    domain.UnitKidinest,
			PatientName:   "Budi Pratama",
			PatientPhone:  "08123456789",
			Category:      domain.CategoryVaksin,
			PaymentMethod: domain.PaymentQRIS,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.PaymentStatus != domain.PaymentPending {
			t.Errorf("expected PENDING, got %s", res.PaymentStatus)
		}
		if res.QueueNumber != "MENUNGGU-BAYAR" {
			t.Errorf("expected MENUNGGU-BAYAR, got %s", res.QueueNumber)
		}

		// Verifikasi pembayaran untuk menerbitkan nomor urut
		verified, err := uc.VerifyPayment(context.Background(), res.QueueID)
		if err != nil {
			t.Fatalf("unexpected error verifying payment: %v", err)
		}
		if verified.PaymentStatus != domain.PaymentPaid {
			t.Errorf("expected PAID, got %s", verified.PaymentStatus)
		}
		if verified.QueueNumber != "A-001" {
			t.Errorf("expected A-001, got %s", verified.QueueNumber)
		}
	})

	t.Run("Pendaftaran SSC booking langsung dapat nomor antrian S-001 tanpa bayar awal", func(t *testing.T) {
		res, err := uc.RegisterQueue(context.Background(), domain.RegisterInput{
			ClinicUnit:    domain.UnitSSC,
			PatientName:   "Clarissa Wijaya",
			PatientPhone:  "08987654321",
			Category:      domain.CategoryTreatment,
			PaymentMethod: domain.PaymentUnpaidBooking,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.RegistrationFee != 0 {
			t.Errorf("expected fee 0, got %d", res.RegistrationFee)
		}
		if res.QueueNumber != "S-001" {
			t.Errorf("expected S-001, got %s", res.QueueNumber)
		}
	})

	t.Run("Masking nama pasien di monitor display harus menyamarkan nama belakang", func(t *testing.T) {
		masked := domain.MaskPatientName("Budi Pratama")
		if masked != "Budi P****" {
			t.Errorf("expected 'Budi P****', got '%s'", masked)
		}
	})
}
