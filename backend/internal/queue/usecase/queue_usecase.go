package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"clinic-api/internal/queue/domain"
)

type QueueUsecase struct {
	patientRepo domain.PatientRepository
	queueRepo   domain.QueueRepository
	waAdapter   WhatsAppAdapter
	baseURL     string
}

func NewQueueUsecase(
	patientRepo domain.PatientRepository,
	queueRepo domain.QueueRepository,
	waAdapter WhatsAppAdapter,
	baseURL string,
) *QueueUsecase {
	if baseURL == "" {
		baseURL = "http://localhost:5173"
	}
	return &QueueUsecase{
		patientRepo: patientRepo,
		queueRepo:   queueRepo,
		waAdapter:   waAdapter,
		baseURL:     baseURL,
	}
}

func (u *QueueUsecase) RegisterQueue(ctx context.Context, input domain.RegisterInput) (*domain.RegisterResult, error) {
	name := strings.TrimSpace(input.PatientName)
	if name == "" {
		return nil, domain.ErrEmptyPatientName
	}
	phone := strings.TrimSpace(input.PatientPhone)
	if phone == "" {
		return nil, domain.ErrEmptyPatientPhone
	}

	if input.ClinicUnit != domain.UnitKidinest && input.ClinicUnit != domain.UnitSSC {
		return nil, domain.ErrInvalidClinicUnit
	}

	// 1. Dapatkan atau buat data pasien (Privacy-First, tanpa NIK)
	patient, err := u.patientRepo.FindOrCreateByPhone(ctx, name, phone)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses data pasien: %w", err)
	}

	docName := strings.TrimSpace(input.DoctorName)
	if docName == "" {
		if input.ClinicUnit == domain.UnitKidinest {
			docName = "dr. Hendra Pratama, Sp.A"
		} else {
			docName = "dr. Nadia Kusuma, Sp.KK"
		}
	}

	now := time.Now()

	// 2. Cabang Logika Unit Bisnis
	if input.ClinicUnit == domain.UnitKidinest {
		// KIDINEST: Bayar registrasi di muka (Rp 50.000). Nomor antrian baru terbit setelah pembayaran lunas (FIFO Invariant).
		q := &domain.QueueItem{
			PatientID:       patient.ID,
			PatientName:     patient.Name,
			PatientPhone:    patient.Phone,
			ClinicUnit:      domain.UnitKidinest,
			Category:        input.Category,
			DoctorID:        input.DoctorID,
			DoctorName:      docName,
			QueueNumber:     "MENUNGGU-BAYAR",
			RegistrationFee: 50000,
			PaymentMethod:   input.PaymentMethod,
			PaymentStatus:   domain.PaymentPending,
			QueueStatus:     domain.QueueWaiting,
		}

		if err := u.queueRepo.Create(ctx, q); err != nil {
			return nil, fmt.Errorf("gagal membuat draft antrian: %w", err)
		}

		return &domain.RegisterResult{
			QueueID:         q.ID,
			QueueNumber:     q.QueueNumber,
			ClinicUnit:      q.ClinicUnit,
			Category:        q.Category,
			DoctorName:      q.DoctorName,
			RegistrationFee: q.RegistrationFee,
			PaymentMethod:   q.PaymentMethod,
			PaymentStatus:   q.PaymentStatus,
			QueueStatus:     q.QueueStatus,
			QRISPayload:     fmt.Sprintf("00020101021226680016ID.CO.KIDINEST.WWW01189360000201100000015204581253033605802ID5914KIDINEST POLI6008SEMARANG62070703A016304%d", q.ID),
			BankTransferAcc: "BCA 828-192-8811 (a.n Kidinest Care)",
			Message:         "Pendaftaran berhasil dicatat. Silakan selesaikan pembayaran registrasi untuk menerbitkan nomor antrian resmi.",
		}, nil
	}

	// SEMARANG SKIN CENTER (SSC): Booking jadwal Rp0 di awal. Nomor antrian terbit seketika, bayar di kasir setelah tindakan.
	queueNum, err := u.queueRepo.GenerateNextQueueNumber(ctx, now, domain.UnitSSC, domain.CategoryTreatment)
	if err != nil {
		return nil, fmt.Errorf("gagal menerbitkan nomor antrian SSC: %w", err)
	}

	q := &domain.QueueItem{
		PatientID:       patient.ID,
		PatientName:     patient.Name,
		PatientPhone:    patient.Phone,
		ClinicUnit:      domain.UnitSSC,
		Category:        domain.CategoryTreatment,
		DoctorID:        input.DoctorID,
		DoctorName:      docName,
		QueueNumber:     queueNum,
		RegistrationFee: 0,
		PaymentMethod:   domain.PaymentUnpaidBooking,
		PaymentStatus:   domain.PaymentPaid, // Bebas bayar registrasi awal
		QueueStatus:     domain.QueueWaiting,
		PaidAt:          &now,
	}

	if err := u.queueRepo.Create(ctx, q); err != nil {
		return nil, fmt.Errorf("gagal membuat tiket antrian SSC: %w", err)
	}

	// Kirim pesan WhatsApp tiket
	trackingURL := fmt.Sprintf("%s/kiosk?ticket=%d", u.baseURL, q.ID)
	_ = u.waAdapter.SendQueueTicket(patient.Phone, patient.Name, q.QueueNumber, string(q.ClinicUnit), trackingURL)

	return &domain.RegisterResult{
		QueueID:         q.ID,
		QueueNumber:     q.QueueNumber,
		ClinicUnit:      q.ClinicUnit,
		Category:        q.Category,
		DoctorName:      q.DoctorName,
		RegistrationFee: 0,
		PaymentMethod:   q.PaymentMethod,
		PaymentStatus:   q.PaymentStatus,
		QueueStatus:     q.QueueStatus,
		Message:         "Tiket reservasi SSC berhasil terbit dan dikirimkan via WhatsApp. Silakan menuju kasir setelah tindakan medis/perawatan selesai.",
	}, nil
}

func (u *QueueUsecase) VerifyPayment(ctx context.Context, queueID int64) (*domain.QueueItem, error) {
	q, err := u.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	if q.PaymentStatus == domain.PaymentPaid {
		return nil, domain.ErrPaymentAlreadyProcessed
	}

	// Generate nomor antrian sequential dengan atomic counter locking
	queueNumber, err := u.queueRepo.GenerateNextQueueNumber(ctx, time.Now(), q.ClinicUnit, q.Category)
	if err != nil {
		return nil, fmt.Errorf("gagal menerbitkan nomor antrian FIFO: %w", err)
	}

	updatedQueue, err := u.queueRepo.UpdatePaymentSuccess(ctx, queueID, queueNumber)
	if err != nil {
		return nil, fmt.Errorf("gagal memperbarui status pembayaran: %w", err)
	}

	// Kirim WhatsApp tiket resmi
	trackingURL := fmt.Sprintf("%s/kiosk?ticket=%d", u.baseURL, updatedQueue.ID)
	_ = u.waAdapter.SendQueueTicket(updatedQueue.PatientPhone, updatedQueue.PatientName, updatedQueue.QueueNumber, string(updatedQueue.ClinicUnit), trackingURL)

	return updatedQueue, nil
}

func (u *QueueUsecase) GetMonitorData(ctx context.Context, unit string) (*domain.MonitorResponse, error) {
	calling, err := u.queueRepo.GetCallingQueue(ctx, unit)
	if err != nil {
		return nil, err
	}

	upcoming, err := u.queueRepo.ListUpcomingQueues(ctx, unit, 5)
	if err != nil {
		return nil, err
	}

	// Privacy-first masking pada monitor display publik
	if calling != nil {
		calling.PatientName = domain.MaskPatientName(calling.PatientName)
	}

	maskedUpcoming := make([]domain.QueueItem, len(upcoming))
	for i, item := range upcoming {
		item.PatientName = domain.MaskPatientName(item.PatientName)
		maskedUpcoming[i] = item
	}

	return &domain.MonitorResponse{
		CurrentlyCalling: calling,
		UpcomingQueue:    maskedUpcoming,
		ClinicUnit:       unit,
	}, nil
}

func (u *QueueUsecase) GetTicket(ctx context.Context, queueID int64) (*domain.QueueItem, error) {
	return u.queueRepo.GetByID(ctx, queueID)
}

func (u *QueueUsecase) CallNext(ctx context.Context, queueID int64) error {
	return u.queueRepo.UpdateQueueStatus(ctx, queueID, domain.QueueCalling)
}

func (u *QueueUsecase) CompleteExam(ctx context.Context, queueID int64) error {
	return u.queueRepo.UpdateQueueStatus(ctx, queueID, domain.QueueCompleted)
}
