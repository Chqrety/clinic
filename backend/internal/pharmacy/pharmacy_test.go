package pharmacy_test

import (
	"testing"

	"clinic-api/internal/billing"
	"clinic-api/internal/consultation"
	"clinic-api/internal/pharmacy"
	"clinic-api/internal/queue"
	"clinic-api/pkg/database"
)

func TestCorePharmacyAndStockLifecycle(t *testing.T) {
	// 1. Initialize in-memory SQLite DB
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init in-memory db: %v", err)
	}
	defer db.Close()

	pharmSvc := pharmacy.NewService(db)
	consSvc := consultation.NewService(db)
	queueSvc := queue.NewService(db)

	// 2. Verify initial seed drugs
	drugs, err := pharmSvc.GetDrugs()
	if err != nil || len(drugs) == 0 {
		t.Fatalf("expected seeded drugs, got %d, err: %v", len(drugs), err)
	}

	pct := drugs[0] // Paracetamol
	initialAvailable := pct.AvailableStock
	if initialAvailable <= 0 {
		t.Fatalf("expected paracetamol available stock > 0, got %d", initialAvailable)
	}

	// 3. Register a test queue & patient
	qItem, err := queueSvc.CreateQueue(queue.CreateQueueRequest{
		NIK:           "3171999988880001",
		Name:          "Test Pasien",
		Phone:         "08111222333",
		Complaint:     "Demam",
		PaymentMethod: "QRIS",
	})
	if err != nil {
		t.Fatalf("failed to create queue: %v", err)
	}

	// 4. Doctor creates consultation with prescription (Reserve 20 tablets)
	requestQty := 20
	_, err = consSvc.SaveConsultation(consultation.SaveConsultationRequest{
		QueueID:    qItem.ID,
		DoctorName: "dr. Hendra Pratama, Sp.PD",
		Diagnosis:  "Febris / Demam",
		ICD10Code:  "R50.9",
		Prescription: []consultation.PrescriptionItemInput{
			{
				DrugID:            pct.ID,
				Quantity:          requestQty,
				DosageInstruction: "3x1 sehari sesudah makan",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to save consultation & prescription: %v", err)
	}

	// 5. Verify Reserved Stock increased, Physical Stock unchanged, Available decreased!
	updatedPct, err := pharmSvc.GetDrugDetail(pct.ID)
	if err != nil {
		t.Fatalf("failed to get drug detail: %v", err)
	}

	if updatedPct.ReservedStock != requestQty {
		t.Errorf("expected reserved_stock = %d, got %d", requestQty, updatedPct.ReservedStock)
	}
	if updatedPct.PhysicalStock != pct.PhysicalStock {
		t.Errorf("expected physical_stock to remain %d, got %d", pct.PhysicalStock, updatedPct.PhysicalStock)
	}
	expectedAvailable := initialAvailable - requestQty
	if updatedPct.AvailableStock != expectedAvailable {
		t.Errorf("expected available_stock = %d, got %d", expectedAvailable, updatedPct.AvailableStock)
	}

	// 6. Pharmacist dispenses medicine
	prescriptions, err := consSvc.GetPrescriptions()
	if err != nil || len(prescriptions) == 0 {
		t.Fatalf("expected prescription in pharmacy, err: %v", err)
	}

	activeRxID := prescriptions[0].ID
	if err := consSvc.DispensePrescription(activeRxID); err != nil {
		t.Fatalf("dispense prescription failed: %v", err)
	}

	// 7. Verify post-dispense: Physical Stock decreased, Reserved Stock cleared!
	finalPct, err := pharmSvc.GetDrugDetail(pct.ID)
	if err != nil {
		t.Fatalf("failed to get drug detail: %v", err)
	}
	if finalPct.ReservedStock != 0 {
		t.Errorf("expected reserved_stock to be 0 after dispense, got %d", finalPct.ReservedStock)
	}
	if finalPct.PhysicalStock != (pct.PhysicalStock - requestQty) {
		t.Errorf("expected physical_stock = %d, got %d", pct.PhysicalStock-requestQty, finalPct.PhysicalStock)
	}

	// 8. Verify Billing: Invoice automatically created and matches 1 consolidated invoice rule
	billingSvc := billing.NewService(db)
	invList, err := billingSvc.GetInvoices()
	if err != nil {
		t.Fatalf("failed to get invoices: %v", err)
	}
	if len(invList) == 0 {
		t.Fatalf("expected at least 1 invoice after pharmacy dispense, got 0")
	}

	inv := invList[0]
	if inv.QueueID != qItem.ID {
		t.Errorf("expected invoice queueID %d, got %d", qItem.ID, inv.QueueID)
	}
	expectedMedFee := pct.SellPrice * requestQty
	if inv.MedicationFee != expectedMedFee {
		t.Errorf("expected medication fee %d, got %d", expectedMedFee, inv.MedicationFee)
	}
	expectedTotal := 35000 + 50000 + expectedMedFee
	if inv.TotalAmount != expectedTotal {
		t.Errorf("expected total invoice %d, got %d", expectedTotal, inv.TotalAmount)
	}
}
