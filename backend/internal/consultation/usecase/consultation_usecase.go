package usecase

import (
	"context"
	"fmt"
	"strings"

	"clinic-api/internal/consultation/domain"
)

type ConsultationUsecase struct {
	repo domain.ConsultationRepository
}

func NewConsultationUsecase(repo domain.ConsultationRepository) *ConsultationUsecase {
	return &ConsultationUsecase{repo: repo}
}

func (uc *ConsultationUsecase) ListDoctorQueues(ctx context.Context, doctorName, clinicUnit string) ([]*domain.DoctorQueuePatient, error) {
	return uc.repo.ListDoctorQueues(ctx, doctorName, clinicUnit)
}

func (uc *ConsultationUsecase) CallPatient(ctx context.Context, queueID int64) error {
	p, err := uc.repo.GetQueuePatientByID(ctx, queueID)
	if err != nil {
		return err
	}
	if p.QueueStatus == "COMPLETED" || p.QueueStatus == "CANCELLED" {
		return fmt.Errorf("antrian sudah selesai atau dibatalkan")
	}
	return uc.repo.UpdateQueueStatus(ctx, queueID, "CALLING")
}

func (uc *ConsultationUsecase) StartExam(ctx context.Context, queueID int64) error {
	p, err := uc.repo.GetQueuePatientByID(ctx, queueID)
	if err != nil {
		return err
	}
	if p.QueueStatus == "COMPLETED" || p.QueueStatus == "CANCELLED" {
		return fmt.Errorf("antrian sudah selesai atau dibatalkan")
	}
	return uc.repo.UpdateQueueStatus(ctx, queueID, "IN_EXAM")
}

func (uc *ConsultationUsecase) GetConsultation(ctx context.Context, queueID int64) (*domain.Consultation, error) {
	return uc.repo.GetConsultationByQueueID(ctx, queueID)
}

func (uc *ConsultationUsecase) ListActiveDrugs(ctx context.Context, clinicUnit string) ([]*domain.CatalogDrug, error) {
	return uc.repo.ListActiveDrugs(ctx, clinicUnit)
}

func (uc *ConsultationUsecase) SubmitConsultation(ctx context.Context, in domain.SaveConsultationInput) (*domain.Consultation, error) {
	in.DoctorName = strings.TrimSpace(in.DoctorName)
	in.Diagnosis = strings.TrimSpace(in.Diagnosis)
	in.Anamnesis = strings.TrimSpace(in.Anamnesis)

	if in.QueueID <= 0 {
		return nil, domain.ErrQueueNotFound
	}
	if in.DoctorName == "" {
		return nil, domain.ErrDoctorNameRequired
	}
	if in.Diagnosis == "" {
		return nil, domain.ErrDiagnosisRequired
	}

	// Prepare prescription items
	var items []domain.PrescriptionItem
	for _, it := range in.PrescriptionItems {
		it.ItemName = strings.TrimSpace(it.ItemName)
		it.DosageInstruction = strings.TrimSpace(it.DosageInstruction)

		if it.Quantity <= 0 {
			return nil, domain.ErrInvalidQuantity
		}

		if it.IsManual {
			// Racikan manual atau obat khusus non-database
			if it.ItemName == "" {
				return nil, fmt.Errorf("nama obat racikan manual wajib diisi")
			}
			if it.UnitPrice < 0 {
				it.UnitPrice = 0
			}
			items = append(items, domain.PrescriptionItem{
				ItemName:          it.ItemName,
				IsManual:          true,
				Quantity:          it.Quantity,
				UnitPrice:         it.UnitPrice,
				Subtotal:          it.Quantity * it.UnitPrice,
				DosageInstruction: it.DosageInstruction,
			})
		} else {
			// Obat katalog
			if it.DrugID == nil || *it.DrugID <= 0 {
				return nil, domain.ErrDrugNotFound
			}
			drug, err := uc.repo.GetDrugByID(ctx, *it.DrugID)
			if err != nil {
				return nil, fmt.Errorf("gagal verifikasi obat katalog ID %d: %w", *it.DrugID, err)
			}
			if drug.AvailableStock < it.Quantity {
				return nil, fmt.Errorf("%w: %s (tersedia: %d, diminta: %d)", domain.ErrInsufficientStock, drug.Name, drug.AvailableStock, it.Quantity)
			}

			unitPrice := drug.SellPrice
			if it.UnitPrice > 0 {
				unitPrice = it.UnitPrice
			}

			items = append(items, domain.PrescriptionItem{
				DrugID:            it.DrugID,
				ItemName:          drug.Name,
				IsManual:          false,
				Quantity:          it.Quantity,
				UnitPrice:         unitPrice,
				Subtotal:          it.Quantity * unitPrice,
				DosageInstruction: it.DosageInstruction,
			})
		}
	}

	cons := &domain.Consultation{
		QueueID:     in.QueueID,
		DoctorID:    in.DoctorID,
		DoctorName:  in.DoctorName,
		Anamnesis:   in.Anamnesis,
		Diagnosis:   in.Diagnosis,
		ICD10Code:   in.ICD10Code,
		ActionNotes: in.ActionNotes,
		ActionFee:   in.ActionFee,
		Status:      "COMPLETED",
	}

	if err := uc.repo.SaveConsultationTx(ctx, cons, items); err != nil {
		return nil, fmt.Errorf("gagal memproses peresepan dokter: %w", err)
	}

	return cons, nil
}
