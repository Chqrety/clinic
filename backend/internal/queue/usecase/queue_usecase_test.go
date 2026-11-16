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

func (m *mockPatientRepo) LookupPatient(ctx context.Context, query string) (*domain.PatientLookupResult, error) {
	for _, p := range m.patients {
		if p.Phone == query || p.Name == query {
			return &domain.PatientLookupResult{
				Patient: *p,
			}, nil
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

func (m *mockQueueRepo) UpdatePaymentDetails(ctx context.Context, id int64, reference, paymentURL string) error {
	q, ok := m.queues[id]
	if !ok {
		return domain.ErrQueueNotFound
	}
	q.PaymentReference = reference
	q.PaymentURL = paymentURL
	return nil
}

func (m *mockQueueRepo) GetByPaymentReference(ctx context.Context, reference string) (*domain.QueueItem, error) {
	for _, q := range m.queues {
		if q.PaymentReference == reference || fmt.Sprintf("KID-%d", q.ID) == reference {
			return q, nil
		}
	}
	return nil, domain.ErrQueueNotFound
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

func (m *mockQueueRepo) ListReservations(ctx context.Context, unit string) ([]domain.QueueItem, error) {
	var list []domain.QueueItem
	for _, q := range m.queues {
		if q.QueueStatus == domain.QueueBooked {
			list = append(list, *q)
		}
	}
	return list, nil
}

func (m *mockQueueRepo) ListRefundRequests(ctx context.Context, unit string) ([]domain.QueueItem, error) {
	var list []domain.QueueItem
	for _, q := range m.queues {
		if q.RefundStatus != nil && *q.RefundStatus != "NONE" {
			list = append(list, *q)
		}
	}
	return list, nil
}

func (m *mockQueueRepo) ListActiveDoctorQueues(ctx context.Context) ([]domain.QueueItem, error) {
	var list []domain.QueueItem
	for _, q := range m.queues {
		if q.QueueStatus == domain.QueueWaiting || q.QueueStatus == domain.QueueCalling || q.QueueStatus == domain.QueueInExam {
			list = append(list, *q)
		}
	}
	return list, nil
}

func (m *mockQueueRepo) Reschedule(ctx context.Context, id int64, newDate, newTime string) error {
	q, ok := m.queues[id]
	if !ok {
		return domain.ErrQueueNotFound
	}
	q.AppointmentDate = &newDate
	q.AppointmentTime = &newTime
	q.RescheduleCount++
	return nil
}

func (m *mockQueueRepo) RequestRefund(ctx context.Context, id int64, bank, accNumber, accName, notes string) error {
	q, ok := m.queues[id]
	if !ok {
		return domain.ErrQueueNotFound
	}
	status := "REQUESTED"
	q.RefundStatus = &status
	q.RefundBank = &bank
	q.RefundAccountNumber = &accNumber
	q.RefundAccountName = &accName
	q.RefundNotes = &notes
	return nil
}

func (m *mockQueueRepo) ProcessRefund(ctx context.Context, id int64, status string, notes string) error {
	q, ok := m.queues[id]
	if !ok {
		return domain.ErrQueueNotFound
	}
	q.RefundStatus = &status
	if status == "APPROVED" {
		q.QueueStatus = domain.QueueCancelled
		q.PaymentStatus = domain.PaymentRefunded
	}
	return nil
}

func TestQueueUsecase_Lifecycle(t *testing.T) {
	pRepo := newMockPatientRepo()
	qRepo := newMockQueueRepo()
	wa := NewMockWhatsAppAdapter()
	uc := NewQueueUsecase(pRepo, qRepo, wa, nil, nil, "http://localhost:5173")

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

	t.Run("Pendaftaran SSC booking langsung dapat nomor antrian S-001 dengan status BOOKED, lalu WAITING setelah check-in GRO", func(t *testing.T) {
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
		if res.QueueStatus != domain.QueueBooked {
			t.Errorf("expected BOOKED status before check-in, got %s", res.QueueStatus)
		}

		// Check-in oleh GRO
		checkedIn, err := uc.CheckIn(context.Background(), res.QueueID)
		if err != nil {
			t.Fatalf("unexpected error during check-in: %v", err)
		}
		if checkedIn.QueueStatus != domain.QueueWaiting {
			t.Errorf("expected WAITING status after check-in, got %s", checkedIn.QueueStatus)
		}

		// Verify Doctor Rooms Status shows ready and waiting count
		rooms, err := uc.GetDoctorRoomsStatus(context.Background())
		if err != nil {
			t.Fatalf("unexpected error getting rooms status: %v", err)
		}
		if len(rooms) == 0 {
			t.Errorf("expected doctor rooms list, got empty")
		}
	})

	t.Run("Masking nama pasien di monitor display harus menyamarkan nama belakang", func(t *testing.T) {
		masked := domain.MaskPatientName("Budi Pratama")
		if masked != "Budi P****" {
			t.Errorf("expected 'Budi P****', got '%s'", masked)
		}
	})

	t.Run("Tripay Webhook: harus memproses status PAID dan menerbitkan nomor antrian resmi", func(t *testing.T) {
		// 1. Register queue
		res, err := uc.RegisterQueue(context.Background(), domain.RegisterInput{
			ClinicUnit:    domain.UnitKidinest,
			PatientName:   "Ahmad Rifai",
			PatientPhone:  "08198765432",
			Category:      domain.CategorySakit,
			PaymentMethod: domain.PaymentQRIS,
		})
		if err != nil {
			t.Fatalf("failed register: %v", err)
		}

		// 2. Simulate Webhook callback
		callbackJSON := fmt.Sprintf(`{
			"reference": "DEV-T12345",
			"merchant_ref": "KID-%d",
			"payment_method": "QRIS2",
			"status": "PAID",
			"total_amount": 50000
		}`, res.QueueID)

		err = uc.ProcessTripayWebhook(context.Background(), []byte(callbackJSON), "")
		if err != nil {
			t.Fatalf("failed process webhook: %v", err)
		}

		// Verify queue status updated
		q, err := qRepo.GetByID(context.Background(), res.QueueID)
		if err != nil {
			t.Fatalf("failed get queue: %v", err)
		}
		if q.PaymentStatus != domain.PaymentPaid {
			t.Errorf("expected PAID status, got %s", q.PaymentStatus)
		}
		if q.QueueNumber != "B-001" {
			t.Errorf("expected queue number B-001, got %s", q.QueueNumber)
		}

		// Test Idempotency: second webhook call should succeed without re-generating
		err = uc.ProcessTripayWebhook(context.Background(), []byte(callbackJSON), "")
		if err != nil {
			t.Fatalf("webhook re-call failed: %v", err)
		}
		q2, _ := qRepo.GetByID(context.Background(), res.QueueID)
		if q2.QueueNumber != "B-001" {
			t.Errorf("expected queue number to remain B-001, got %s", q2.QueueNumber)
		}
	})

	t.Run("Webhook Midtrans berhasil menyelesaikan pembayaran dan menerbitkan nomor antrian", func(t *testing.T) {
		res, err := uc.RegisterQueue(context.Background(), domain.RegisterInput{
			ClinicUnit:    domain.UnitKidinest,
			PatientName:   "Dewi Lestari",
			PatientPhone:  "08198765432",
			Category:      domain.CategorySakit,
			PaymentMethod: domain.PaymentQRIS,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		midtransJSON := fmt.Sprintf(`{
			"order_id": "KID-%d",
			"transaction_id": "mid-tx-9999",
			"transaction_status": "settlement",
			"status_code": "200",
			"gross_amount": "50000.00",
			"signature_key": "dummy"
		}`, res.QueueID)

		err = uc.ProcessMidtransWebhook(context.Background(), []byte(midtransJSON))
		if err != nil {
			t.Fatalf("failed process midtrans webhook: %v", err)
		}

		q, err := qRepo.GetByID(context.Background(), res.QueueID)
		if err != nil {
			t.Fatalf("failed get queue: %v", err)
		}
		if q.PaymentStatus != domain.PaymentPaid {
			t.Errorf("expected PAID status, got %s", q.PaymentStatus)
		}
		if q.QueueNumber == "" {
			t.Errorf("expected generated queue number, got empty")
		}

		// Idempotency
		err = uc.ProcessMidtransWebhook(context.Background(), []byte(midtransJSON))
		if err != nil {
			t.Fatalf("second midtrans webhook call failed: %v", err)
		}
	})
}
