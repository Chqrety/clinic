package usecase_test

import (
	"context"
	"testing"
	"time"

	"clinic-api/internal/consultation/domain"
	"clinic-api/internal/consultation/usecase"
)

type mockConsultationRepo struct {
	queues        map[int64]*domain.DoctorQueuePatient
	drugs         map[int64]*domain.CatalogDrug
	consultations map[int64]*domain.Consultation
}

func (m *mockConsultationRepo) ListDoctorQueues(ctx context.Context, doctorName, clinicUnit string) ([]*domain.DoctorQueuePatient, error) {
	var list []*domain.DoctorQueuePatient
	for _, q := range m.queues {
		if clinicUnit != "" && clinicUnit != "ALL" && q.ClinicUnit != clinicUnit {
			continue
		}
		if doctorName != "" && q.DoctorName != doctorName && q.DoctorName != "" {
			continue
		}
		list = append(list, q)
	}
	return list, nil
}

func (m *mockConsultationRepo) GetQueuePatientByID(ctx context.Context, queueID int64) (*domain.DoctorQueuePatient, error) {
	if q, ok := m.queues[queueID]; ok {
		return q, nil
	}
	return nil, domain.ErrQueueNotFound
}

func (m *mockConsultationRepo) UpdateQueueStatus(ctx context.Context, queueID int64, status string) error {
	if q, ok := m.queues[queueID]; ok {
		q.QueueStatus = status
		now := time.Now()
		if status == "CALLING" {
			q.CalledAt = &now
		}
		return nil
	}
	return domain.ErrQueueNotFound
}

func (m *mockConsultationRepo) GetConsultationByQueueID(ctx context.Context, queueID int64) (*domain.Consultation, error) {
	if c, ok := m.consultations[queueID]; ok {
		return c, nil
	}
	return nil, domain.ErrConsultationNotFound
}

func (m *mockConsultationRepo) SaveConsultationTx(ctx context.Context, c *domain.Consultation, items []domain.PrescriptionItem) error {
	c.ID = int64(len(m.consultations) + 1)
	c.Prescription = &domain.Prescription{
		ID:             c.ID,
		ConsultationID: c.ID,
		QueueID:        c.QueueID,
		Status:         "RESERVED",
		Items:          items,
	}
	m.consultations[c.QueueID] = c

	// Reserve stock
	for _, it := range items {
		if !it.IsManual && it.DrugID != nil {
			if d, ok := m.drugs[*it.DrugID]; ok {
				d.ReservedStock += it.Quantity
				d.AvailableStock -= it.Quantity
			}
		}
	}

	if q, ok := m.queues[c.QueueID]; ok {
		q.QueueStatus = "PHARMACY_WAITING"
	}
	return nil
}

func (m *mockConsultationRepo) ListActiveDrugs(ctx context.Context, clinicUnit string) ([]*domain.CatalogDrug, error) {
	var list []*domain.CatalogDrug
	for _, d := range m.drugs {
		if clinicUnit != "" && clinicUnit != "ALL" && d.ClinicUnit != clinicUnit && d.ClinicUnit != "ALL" {
			continue
		}
		list = append(list, d)
	}
	return list, nil
}

func (m *mockConsultationRepo) GetDrugByID(ctx context.Context, id int64) (*domain.CatalogDrug, error) {
	if d, ok := m.drugs[id]; ok {
		return d, nil
	}
	return nil, domain.ErrDrugNotFound
}

func setupConsultationUsecase() (*usecase.ConsultationUsecase, *mockConsultationRepo) {
	repo := &mockConsultationRepo{
		queues:        make(map[int64]*domain.DoctorQueuePatient),
		drugs:         make(map[int64]*domain.CatalogDrug),
		consultations: make(map[int64]*domain.Consultation),
	}

	repo.drugs[1] = &domain.CatalogDrug{
		ID:             1,
		ClinicUnit:     "KIDINEST",
		Name:           "Paracetamol Drop 100mg/ml",
		SellPrice:      35000,
		PhysicalStock:  50,
		ReservedStock:  0,
		AvailableStock: 50,
		IsActive:       true,
	}

	repo.queues[100] = &domain.DoctorQueuePatient{
		QueueID:       100,
		QueueNumber:   "A-001",
		PatientName:   "Ananda Kenzo",
		ClinicUnit:    "KIDINEST",
		DoctorName:    "dr. Hendra Pratama, Sp.A",
		PaymentStatus: "PAID",
		QueueStatus:   "WAITING",
		CreatedAt:     time.Now(),
	}

	uc := usecase.NewConsultationUsecase(repo)
	return uc, repo
}

func TestCallAndStartExam(t *testing.T) {
	uc, repo := setupConsultationUsecase()
	ctx := context.Background()

	// 1. Call Patient
	err := uc.CallPatient(ctx, 100)
	if err != nil {
		t.Fatalf("unexpected error calling patient: %v", err)
	}
	if repo.queues[100].QueueStatus != "CALLING" {
		t.Errorf("expected status CALLING, got %s", repo.queues[100].QueueStatus)
	}

	// 2. Start Exam
	err = uc.StartExam(ctx, 100)
	if err != nil {
		t.Fatalf("unexpected error starting exam: %v", err)
	}
	if repo.queues[100].QueueStatus != "IN_EXAM" {
		t.Errorf("expected status IN_EXAM, got %s", repo.queues[100].QueueStatus)
	}
}

func TestSubmitConsultationValidation(t *testing.T) {
	uc, _ := setupConsultationUsecase()
	ctx := context.Background()

	// Missing doctor name
	_, err := uc.SubmitConsultation(ctx, domain.SaveConsultationInput{
		QueueID:    100,
		DoctorName: "",
		Diagnosis:  "ISPA Akut",
	})
	if err != domain.ErrDoctorNameRequired {
		t.Errorf("expected ErrDoctorNameRequired, got %v", err)
	}

	// Missing diagnosis
	_, err = uc.SubmitConsultation(ctx, domain.SaveConsultationInput{
		QueueID:    100,
		DoctorName: "dr. Hendra Pratama, Sp.A",
		Diagnosis:  "",
	})
	if err != domain.ErrDiagnosisRequired {
		t.Errorf("expected ErrDiagnosisRequired, got %v", err)
	}
}

func TestSubmitConsultationWithCatalogAndManualItems(t *testing.T) {
	uc, repo := setupConsultationUsecase()
	ctx := context.Background()

	drugID := int64(1)
	input := domain.SaveConsultationInput{
		QueueID:     100,
		DoctorName:  "dr. Hendra Pratama, Sp.A",
		Anamnesis:   "Demam 38.5C sejak semalam, batuk pilek ringan",
		Diagnosis:   "ISPA / Infeksi Saluran Pernapasan Akut",
		ICD10Code:   "J06.9",
		ActionNotes: "Edukasi hidrasi anak & kompres hangat",
		ActionFee:   75000,
		PrescriptionItems: []domain.RxItemInput{
			// Item 1: Catalog Drug
			{
				DrugID:            &drugID,
				IsManual:          false,
				Quantity:          2,
				DosageInstruction: "3x 0.6ml sehari sesudah makan",
			},
			// Item 2: Manual Racikan
			{
				ItemName:          "Puyer Batuk Pilek Racikan (Ambroxol + CTM)",
				IsManual:          true,
				Quantity:          10,
				UnitPrice:         3500,
				DosageInstruction: "3x 1 bungkus sesudah makan",
			},
		},
	}

	cons, err := uc.SubmitConsultation(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error submitting consultation: %v", err)
	}

	if cons.Prescription == nil || len(cons.Prescription.Items) != 2 {
		t.Fatalf("expected 2 prescription items, got %v", cons.Prescription)
	}

	// Verify catalog item calculations & stock reservation
	it1 := cons.Prescription.Items[0]
	if it1.ItemName != "Paracetamol Drop 100mg/ml" {
		t.Errorf("expected item name Paracetamol Drop, got %s", it1.ItemName)
	}
	if it1.Subtotal != 70000 {
		t.Errorf("expected subtotal 70000, got %d", it1.Subtotal)
	}
	if repo.drugs[1].ReservedStock != 2 {
		t.Errorf("expected reserved stock 2, got %d", repo.drugs[1].ReservedStock)
	}
	if repo.drugs[1].AvailableStock != 48 {
		t.Errorf("expected available stock 48, got %d", repo.drugs[1].AvailableStock)
	}

	// Verify manual item calculations
	it2 := cons.Prescription.Items[1]
	if !it2.IsManual {
		t.Errorf("expected is_manual true")
	}
	if it2.Subtotal != 35000 {
		t.Errorf("expected manual subtotal 35000, got %d", it2.Subtotal)
	}

	// Verify queue transition to PHARMACY_WAITING
	if repo.queues[100].QueueStatus != "PHARMACY_WAITING" {
		t.Errorf("expected queue status PHARMACY_WAITING, got %s", repo.queues[100].QueueStatus)
	}
}

func TestStockExceededError(t *testing.T) {
	uc, _ := setupConsultationUsecase()
	ctx := context.Background()

	drugID := int64(1)
	input := domain.SaveConsultationInput{
		QueueID:    100,
		DoctorName: "dr. Hendra Pratama, Sp.A",
		Diagnosis:  "ISPA",
		PrescriptionItems: []domain.RxItemInput{
			{
				DrugID:            &drugID,
				IsManual:          false,
				Quantity:          999, // Exceeds available 50
				DosageInstruction: "Overdose check",
			},
		},
	}

	_, err := uc.SubmitConsultation(ctx, input)
	if err == nil {
		t.Fatalf("expected error when stock is exceeded, got nil")
	}
}
