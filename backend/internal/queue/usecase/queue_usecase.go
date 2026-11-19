package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"clinic-api/internal/queue/domain"
	"clinic-api/internal/shared/midtrans"
	"clinic-api/internal/shared/tripay"
)

type QueueUsecase struct {
	patientRepo    domain.PatientRepository
	queueRepo      domain.QueueRepository
	waAdapter      WhatsAppAdapter
	tripayClient   *tripay.Client
	midtransClient *midtrans.Client
	baseURL        string
}

func NewQueueUsecase(
	patientRepo domain.PatientRepository,
	queueRepo domain.QueueRepository,
	waAdapter WhatsAppAdapter,
	tripayClient *tripay.Client,
	midtransClient *midtrans.Client,
	baseURL string,
) *QueueUsecase {
	if baseURL == "" {
		baseURL = "http://localhost:5173"
	}
	return &QueueUsecase{
		patientRepo:    patientRepo,
		queueRepo:      queueRepo,
		waAdapter:      waAdapter,
		tripayClient:   tripayClient,
		midtransClient: midtransClient,
		baseURL:        baseURL,
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

	pType := strings.ToUpper(strings.TrimSpace(input.PatientType))
	if pType != "LAMA" {
		pType = "BARU"
	}

	now := time.Now()
	today := now.Format("2006-01-02")
	isFuture := false
	if input.AppointmentDate != nil && strings.TrimSpace(*input.AppointmentDate) != "" {
		if strings.TrimSpace(*input.AppointmentDate) > today {
			isFuture = true
		}
	}

	// 2. Cabang Logika Unit Bisnis
	defaultTime := "09:00"
	if input.AppointmentTime == nil || *input.AppointmentTime == "" {
		input.AppointmentTime = &defaultTime
	}

	if input.ClinicUnit == domain.UnitKidinest {
		// KIDINEST: Bayar registrasi di muka (Rp 50.000). Nomor antrian baru terbit setelah pembayaran lunas (FIFO Invariant).
		initStatus := domain.QueueWaiting
		if isFuture {
			initStatus = domain.QueueBooked
		}

		// Sesuai Master Prompt Bagian 3 C:
		// A, B, C: Rp 200.000, D: Rp 250.000
		regFee := 200000
		if input.Category == domain.CategoryPaket || input.Category == "D_PAKET" || input.Category == "D" {
			regFee = 250000
		}

		q := &domain.QueueItem{
			PatientID:       patient.ID,
			PatientName:     patient.Name,
			PatientPhone:    patient.Phone,
			PatientType:     pType,
			ClinicUnit:      domain.UnitKidinest,
			Category:        input.Category,
			DoctorID:        input.DoctorID,
			DoctorName:      docName,
			AppointmentDate: input.AppointmentDate,
			AppointmentTime: input.AppointmentTime,
			QueueNumber:     "MENUNGGU-BAYAR",
			RegistrationFee: regFee,
			PaymentMethod:   input.PaymentMethod,
			PaymentStatus:   domain.PaymentPending,
			QueueStatus:     initStatus,
		}

		if err := u.queueRepo.Create(ctx, q); err != nil {
			return nil, fmt.Errorf("gagal membuat draft antrian: %w", err)
		}

		paymentRef := fmt.Sprintf("SIM-KID-%d", q.ID)
		paymentURL := ""
		snapToken := ""
		clientKey := ""
		qrisPayload := ""

		if input.PaymentMethod == domain.PaymentQRIS {
			orderRef := fmt.Sprintf("KID-%d-%d", q.ID, time.Now().Unix())

			if u.midtransClient != nil && u.midtransClient.IsConfigured() {
				log.Printf("[MIDTRANS] Meminta gambar QRIS langsung untuk %s (Nominal: Rp %d)...", orderRef, q.RegistrationFee)
				midResp, err := u.midtransClient.CreateQRIS(ctx, orderRef, q.RegistrationFee, patient.Name, patient.Phone)
				if err == nil && midResp != nil {
					paymentRef = midResp.OrderID
					if midResp.TransactionID != "" {
						paymentRef = midResp.TransactionID
					}
					paymentURL = midResp.QRURL
					if midResp.QRString != "" {
						qrisPayload = midResp.QRString
					}
					log.Printf("[MIDTRANS] Gambar QRIS berhasil didapatkan kilat! Ref: %s, QRURL: %s", paymentRef, paymentURL)
					_ = u.queueRepo.UpdatePaymentDetails(ctx, q.ID, orderRef, paymentURL)
				} else {
					log.Printf("[MIDTRANS] Core API QRIS error (%v). Mencoba Snap fallback...", err)
					snapResp, snapErr := u.midtransClient.CreateSnapTransaction(ctx, orderRef, q.RegistrationFee, patient.Name, patient.Phone)
					if snapErr == nil && snapResp != nil {
						paymentRef = orderRef
						paymentURL = snapResp.RedirectURL
						snapToken = snapResp.Token
						clientKey = u.midtransClient.ClientKey()
						log.Printf("[MIDTRANS] Snap fallback dibuat! Token: %s, URL: %s", snapToken, paymentURL)

						// Solusi 1: Handshake internal Snap untuk mengekstrak URL gambar QRIS resmi (merchants-app.../qr-code)
						if qrImg, qrErr := u.midtransClient.FetchSnapGoPayQR(ctx, snapResp.Token); qrErr == nil && qrImg != "" {
							paymentURL = qrImg
							log.Printf("[MIDTRANS] Solusi 1 Sukses! Gambar QR asli resmi diperoleh: %s", paymentURL)
						} else {
							log.Printf("[MIDTRANS] Solusi 1 Handshake note: %v", qrErr)
						}

						_ = u.queueRepo.UpdatePaymentDetails(ctx, q.ID, orderRef, paymentURL)
					}
				}
			} else if u.tripayClient != nil && u.tripayClient.IsConfigured() {
				tripayResp, err := u.tripayClient.CreateClosedQRIS(ctx, orderRef, q.RegistrationFee, patient.Name, patient.Phone)
				if err == nil && tripayResp != nil {
					paymentRef = tripayResp.Reference
					paymentURL = tripayResp.QRURL
					if paymentURL == "" {
						paymentURL = tripayResp.CheckoutURL
					}
					if tripayResp.QRString != "" {
						qrisPayload = tripayResp.QRString
					}
					_ = u.queueRepo.UpdatePaymentDetails(ctx, q.ID, paymentRef, paymentURL)
				}
			}

			// Pastikan jika paymentURL belum berupa gambar langsung tetapi qrisPayload ada, sediakan QR code scannable
			if paymentURL == "" && qrisPayload != "" {
				paymentURL = fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s", url.QueryEscape(qrisPayload))
			}
		}

		return &domain.RegisterResult{
			QueueID:          q.ID,
			QueueNumber:      q.QueueNumber,
			ClinicUnit:       q.ClinicUnit,
			Category:         q.Category,
			DoctorName:       q.DoctorName,
			AppointmentDate:  q.AppointmentDate,
			AppointmentTime:  q.AppointmentTime,
			PatientType:      q.PatientType,
			RegistrationFee:  q.RegistrationFee,
			PaymentMethod:    q.PaymentMethod,
			PaymentStatus:    q.PaymentStatus,
			QueueStatus:      q.QueueStatus,
			PaymentReference: paymentRef,
			PaymentURL:       paymentURL,
			SnapToken:        snapToken,
			ClientKey:        clientKey,
			QRISPayload:      qrisPayload,
			BankTransferAcc:  "BCA 828-192-8811 (a.n Kidinest Care)",
			Message:          "Pendaftaran berhasil dicatat. Silakan selesaikan pembayaran registrasi untuk menerbitkan nomor antrian resmi.",
		}, nil
	}

	// SEMARANG SKIN CENTER (SSC): Booking jadwal Rp0 di awal.
	targetDateSSC := now
	if input.AppointmentDate != nil && *input.AppointmentDate != "" {
		if t, err := time.Parse("2006-01-02", *input.AppointmentDate); err == nil {
			targetDateSSC = t
		}
	}
	queueNum, err := u.queueRepo.GenerateNextQueueNumber(ctx, targetDateSSC, domain.UnitSSC, domain.CategoryTreatment)
	if err != nil {
		return nil, fmt.Errorf("gagal menerbitkan nomor antrian SSC: %w", err)
	}

	q := &domain.QueueItem{
		PatientID:       patient.ID,
		PatientName:     patient.Name,
		PatientPhone:    patient.Phone,
		PatientType:     pType,
		ClinicUnit:      domain.UnitSSC,
		Category:        domain.CategoryTreatment,
		DoctorID:        input.DoctorID,
		DoctorName:      docName,
		AppointmentDate: input.AppointmentDate,
		AppointmentTime: input.AppointmentTime,
		QueueNumber:     queueNum,
		RegistrationFee: 0,
		PaymentMethod:   domain.PaymentUnpaidBooking,
		PaymentStatus:   domain.PaymentPaid, // Bebas bayar registrasi awal
		QueueStatus:     domain.QueueBooked,
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
		AppointmentDate: q.AppointmentDate,
		AppointmentTime: q.AppointmentTime,
		PatientType:     q.PatientType,
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

	targetDate := time.Now()
	if q.AppointmentDate != nil && *q.AppointmentDate != "" {
		if t, err := time.Parse("2006-01-02", *q.AppointmentDate); err == nil {
			targetDate = t
		}
	}

	queueNumber, err := u.queueRepo.GenerateNextQueueNumber(ctx, targetDate, q.ClinicUnit, q.Category)
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

func (u *QueueUsecase) completePaymentSuccess(ctx context.Context, refs ...string) error {
	var q *domain.QueueItem
	var err error
	for _, ref := range refs {
		if strings.TrimSpace(ref) == "" {
			continue
		}
		q, err = u.queueRepo.GetByPaymentReference(ctx, ref)
		if err == nil && q != nil {
			break
		}
	}
	if q == nil {
		return fmt.Errorf("antrian tidak ditemukan untuk referensi %v: %w", refs, err)
	}

	// Idempotency Check: jangan terbitkan nomor ganda jika sudah berstatus PAID
	if q.PaymentStatus == domain.PaymentPaid {
		return nil
	}

	targetDate := time.Now()
	if q.AppointmentDate != nil && *q.AppointmentDate != "" {
		if t, err := time.Parse("2006-01-02", *q.AppointmentDate); err == nil {
			targetDate = t
		}
	}

	queueNumber, err := u.queueRepo.GenerateNextQueueNumber(ctx, targetDate, q.ClinicUnit, q.Category)
	if err != nil {
		return fmt.Errorf("gagal menerbitkan nomor antrian FIFO: %w", err)
	}

	updatedQueue, err := u.queueRepo.UpdatePaymentSuccess(ctx, q.ID, queueNumber)
	if err != nil {
		if errors.Is(err, domain.ErrQueueNotFound) {
			return nil // Idempotent exit: already marked paid by concurrent webhook
		}
		return fmt.Errorf("gagal memperbarui status pembayaran: %w", err)
	}

	// Kirim tiket ke WhatsApp
	trackingURL := fmt.Sprintf("%s/kiosk?ticket=%d", u.baseURL, updatedQueue.ID)
	_ = u.waAdapter.SendQueueTicket(updatedQueue.PatientPhone, updatedQueue.PatientName, updatedQueue.QueueNumber, string(updatedQueue.ClinicUnit), trackingURL)

	return nil
}

func (u *QueueUsecase) ProcessTripayWebhook(ctx context.Context, rawBody []byte, signature string) error {
	if u.tripayClient != nil && u.tripayClient.IsConfigured() {
		if !u.tripayClient.ValidateCallbackSignature(rawBody, signature) {
			return errors.New("signature webhook tripay tidak valid")
		}
	}

	var payload tripay.CallbackPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return fmt.Errorf("payload webhook tripay tidak valid: %w", err)
	}

	if payload.Status != "PAID" {
		return nil // Hanya proses jika status PAID
	}

	return u.completePaymentSuccess(ctx, payload.MerchantRef, payload.Reference)
}

func (u *QueueUsecase) ProcessMidtransWebhook(ctx context.Context, rawBody []byte) error {
	var payload midtrans.NotificationPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return fmt.Errorf("payload webhook midtrans tidak valid: %w", err)
	}

	if u.midtransClient != nil && u.midtransClient.IsConfigured() {
		if !u.midtransClient.ValidateNotificationSignature(payload.OrderID, payload.StatusCode, payload.GrossAmount, payload.SignatureKey) {
			return errors.New("signature webhook midtrans tidak valid")
		}
	}

	if !payload.IsPaymentSuccess() {
		return nil // Abaikan jika belum settlement atau capture accept
	}

	return u.completePaymentSuccess(ctx, payload.OrderID, payload.TransactionID)
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
	q, err := u.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	// Sinkronisasi status pembayaran dengan Midtrans secara otomatis jika masih pending
	if q.PaymentStatus == domain.PaymentPending && q.PaymentMethod == domain.PaymentQRIS {
		if u.midtransClient != nil && u.midtransClient.IsConfigured() && q.PaymentReference != "" {
			status, err := u.midtransClient.CheckTransactionStatus(ctx, q.PaymentReference)
			if err == nil && status != nil {
				if status.TransactionStatus == "settlement" || (status.TransactionStatus == "capture" && status.FraudStatus == "accept") {
					log.Printf("[MIDTRANS] Status API mendeteksi pembayaran lunas untuk order %s! Menerbitkan nomor antrian FIFO...", q.PaymentReference)
					updated, err := u.VerifyPayment(ctx, queueID)
					if err == nil && updated != nil {
						return updated, nil
					}
				}
			}
		}
	}

	return q, nil
}

func (u *QueueUsecase) CallNext(ctx context.Context, queueID int64) error {
	return u.queueRepo.UpdateQueueStatus(ctx, queueID, domain.QueueCalling)
}

func (u *QueueUsecase) CompleteExam(ctx context.Context, queueID int64) error {
	return u.queueRepo.UpdateQueueStatus(ctx, queueID, domain.QueueCompleted)
}

func (u *QueueUsecase) CheckIn(ctx context.Context, queueID int64) (*domain.QueueItem, error) {
	q, err := u.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, err
	}
	if q.QueueStatus != domain.QueueBooked && q.QueueStatus != domain.QueueWaiting {
		return nil, fmt.Errorf("antrian tidak dalam status reservasi/menunggu")
	}

	if err := u.queueRepo.UpdateQueueStatus(ctx, queueID, domain.QueueWaiting); err != nil {
		return nil, fmt.Errorf("gagal check-in antrian: %w", err)
	}
	q.QueueStatus = domain.QueueWaiting
	return q, nil
}

func (u *QueueUsecase) ListReservations(ctx context.Context, unit string) ([]domain.QueueItem, error) {
	return u.queueRepo.ListReservations(ctx, unit)
}

func (u *QueueUsecase) ListRefundRequests(ctx context.Context, unit string) ([]domain.QueueItem, error) {
	return u.queueRepo.ListRefundRequests(ctx, unit)
}

func isDoctorMatch(doctorNameInQueue, targetDoctorName string) bool {
	docA := strings.ToLower(strings.TrimSpace(doctorNameInQueue))
	docB := strings.ToLower(strings.TrimSpace(targetDoctorName))
	if docA == docB {
		return true
	}
	if strings.Contains(docA, "patricia") && strings.Contains(docB, "patricia") {
		return true
	}
	if strings.Contains(docA, "hendra") && strings.Contains(docB, "hendra") {
		return true
	}
	if strings.Contains(docA, "dian") && strings.Contains(docB, "dian") {
		return true
	}
	if strings.Contains(docA, "kevin") && strings.Contains(docB, "kevin") {
		return true
	}
	return false
}

func (u *QueueUsecase) GetDoctorRoomsStatus(ctx context.Context) ([]domain.DoctorRoomStatus, error) {
	activeQueues, err := u.queueRepo.ListActiveDoctorQueues(ctx)
	if err != nil {
		return nil, err
	}

	doctors := []struct {
		name string
		unit domain.ClinicUnit
		room string
	}{
		{"dr. Hendra Pratama, Sp.A", domain.UnitKidinest, "Poli Anak 1"},
		{"dr. Patricia, Sp.A", domain.UnitKidinest, "Poli Anak 2"},
		{"dr. Dian Paramita, Sp.D.V.E", domain.UnitSSC, "Poli Estetika 1"},
		{"dr. Kevin Wijaya, Sp.D.V.E", domain.UnitSSC, "Poli Estetika 2"},
	}

	var results []domain.DoctorRoomStatus
	for _, doc := range doctors {
		room := domain.DoctorRoomStatus{
			DoctorName: doc.name,
			ClinicUnit: doc.unit,
			RoomName:   doc.room,
			Status:     "READY",
		}

		for _, q := range activeQueues {
			if isDoctorMatch(q.DoctorName, doc.name) {
				if q.QueueStatus == domain.QueueCalling || q.QueueStatus == domain.QueueInExam {
					room.Status = "BUSY"
					curr := q
					curr.PatientName = domain.MaskPatientName(curr.PatientName)
					room.CurrentPatient = &curr
				} else if q.QueueStatus == domain.QueueWaiting {
					room.WaitingCount++
					if room.NextPatient == nil {
						next := q
						next.PatientName = domain.MaskPatientName(next.PatientName)
						room.NextPatient = &next
					}
				}
			}
		}

		results = append(results, room)
	}

	return results, nil
}

func (u *QueueUsecase) LookupPatient(ctx context.Context, query string) (*domain.PatientLookupResult, error) {
	return u.patientRepo.LookupPatient(ctx, query)
}

func (u *QueueUsecase) RescheduleQueue(ctx context.Context, queueID int64, newDate, newTime string) (*domain.QueueItem, error) {
	q, err := u.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	if q.QueueStatus == domain.QueueCompleted || q.QueueStatus == domain.QueueCancelled {
		return nil, errors.New("antrian sudah selesai atau dibatalkan, tidak dapat di-reschedule")
	}

	if err := u.queueRepo.Reschedule(ctx, queueID, newDate, newTime); err != nil {
		return nil, fmt.Errorf("gagal reschedule jadwal: %w", err)
	}

	updated, err := u.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	// Kirim notifikasi WhatsApp
	msg := fmt.Sprintf("Halo %s, jadwal kunjungan Anda di %s berhasil diubah ke tanggal %s jam %s. Kode Booking: %s.",
		updated.PatientName, updated.ClinicUnit, newDate, newTime, updated.QueueNumber)
	_ = u.waAdapter.SendMessage(updated.PatientPhone, msg)

	return updated, nil
}

func (u *QueueUsecase) RequestRefund(ctx context.Context, queueID int64, bank, accNumber, accName, notes string) (*domain.QueueItem, error) {
	q, err := u.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	if q.PaymentStatus != domain.PaymentPaid {
		return nil, errors.New("hanya antrian berbayar lunas yang dapat mengajukan refund")
	}

	if err := u.queueRepo.RequestRefund(ctx, queueID, bank, accNumber, accName, notes); err != nil {
		return nil, fmt.Errorf("gagal mengajukan refund: %w", err)
	}

	return u.queueRepo.GetByID(ctx, queueID)
}

func (u *QueueUsecase) ProcessRefund(ctx context.Context, queueID int64, status, notes string) (*domain.QueueItem, error) {
	if _, err := u.queueRepo.GetByID(ctx, queueID); err != nil {
		return nil, err
	}

	if err := u.queueRepo.ProcessRefund(ctx, queueID, status, notes); err != nil {
		return nil, fmt.Errorf("gagal memproses refund: %w", err)
	}

	updated, err := u.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	if status == "APPROVED" {
		msg := fmt.Sprintf("Halo %s, pengajuan refund registrasi Anda (%s) telah disetujui kasir. Dana telah ditransfer manual ke rekening Anda. Catatan: %s",
			updated.PatientName, updated.QueueNumber, notes)
		_ = u.waAdapter.SendMessage(updated.PatientPhone, msg)
	}

	return updated, nil
}

