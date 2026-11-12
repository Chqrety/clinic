package dev

import (
	"database/sql"
	"fmt"
	"time"
)

// WipeTransactions clears all transactional data leaving master staff & drugs intact.
func WipeTransactions(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi database: %w", err)
	}
	defer tx.Rollback()

	queries := []string{
		"DELETE FROM stock_opname_logs",
		"DELETE FROM prescription_items",
		"DELETE FROM prescriptions",
		"DELETE FROM consultations",
		"DELETE FROM cashier_shift_reports",
		"DELETE FROM invoices",
		"DELETE FROM queues",
		"DELETE FROM daily_queue_counters",
		"DELETE FROM patients",
		"UPDATE drug_batches SET reserved_stock = 0",
	}

	for _, q := range queries {
		if _, err := tx.Exec(q); err != nil {
			return fmt.Errorf("gagal truncate tabel (%s): %w", q, err)
		}
	}

	return tx.Commit()
}

// SeedInitialDemo seeds ~8 initial clean queues for standard demonstration.
func SeedInitialDemo(db *sql.DB) error {
	if err := WipeTransactions(db); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	type InitialItem struct {
		name       string
		phone      string
		unit       string
		cat        string
		docName    string
		qNum       string
		regFee     int
		payMethod  string
		payStatus  string
		qStatus    string
	}

	initials := []InitialItem{
		{"Kenzo Alvaro", "081234567801", "KIDINEST", "A_VAKSIN", "dr. Hendra Pratama, Sp.A", "A-001", 35000, "QRIS", "PAID", "CALLING"},
		{"Aisyah Putri", "081234567802", "KIDINEST", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-001", 35000, "CASH", "PAID", "WAITING"},
		{"Rafathar Malik", "081234567803", "KIDINEST", "C_TUMBANG", "dr. Hendra Pratama, Sp.A", "C-001", 50000, "BANK_TRANSFER", "PAID", "WAITING"},
		{"Clarissa Aurelia", "081234567804", "KIDINEST", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-002", 35000, "QRIS", "PAID", "PHARMACY_WAITING"},
		{"Bima Sakti Nugraha", "081234567805", "KIDINEST", "A_VAKSIN", "dr. Hendra Pratama, Sp.A", "A-002", 35000, "CASH", "PAID", "COMPLETED"},
		{"Nadia Saphira", "081234567806", "SSC", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-001", 50000, "QRIS", "PAID", "IN_EXAM"},
		{"Jessica Mila", "081234567807", "SSC", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-002", 50000, "BANK_TRANSFER", "PAID", "WAITING"},
		{"Kevin Sanjaya", "081234567808", "SSC", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-003", 50000, "CASH", "PAID", "COMPLETED"},
	}

	for idx, it := range initials {
		var patID int
		err := tx.QueryRow(`
			INSERT INTO patients (name, phone, created_at, updated_at)
			VALUES ($1, $2, $3, $3) RETURNING id
		`, it.name, it.phone, now.Add(-time.Duration(60-idx*5)*time.Minute)).Scan(&patID)
		if err != nil {
			return err
		}

		var qID int
		err = tx.QueryRow(`
			INSERT INTO queues (
				patient_id, clinic_unit, category, doctor_name, queue_number,
				registration_fee, payment_method, payment_status, queue_status,
				paid_at, called_at, completed_at, created_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9,
				$10, $11, $12, $13
			) RETURNING id
		`, patID, it.unit, it.cat, it.docName, it.qNum,
			it.regFee, it.payMethod, it.payStatus, it.qStatus,
			now.Add(-time.Duration(45-idx*5)*time.Minute),
			now.Add(-time.Duration(30-idx*5)*time.Minute),
			now.Add(-time.Duration(15-idx*5)*time.Minute),
			now.Add(-time.Duration(60-idx*5)*time.Minute),
		).Scan(&qID)
		if err != nil {
			return err
		}

		if it.qStatus == "COMPLETED" {
			invNum := fmt.Sprintf("INV-%s-%04d", now.Format("20060102"), idx+1)
			_, err = tx.Exec(`
				INSERT INTO invoices (
					queue_id, patient_id, invoice_number,
					subtotal_registration, subtotal_action, subtotal_medication,
					voucher_discount, total_amount, payment_method,
					cash_received, cash_change, cashier_name, status, created_at, updated_at
				) VALUES (
					$1, $2, $3,
					$4, 150000, 75000,
					0, $4 + 150000 + 75000, $5,
					300000, 300000 - ($4 + 150000 + 75000), 'Siti Sarah, A.Md', 'PAID', $6, $6
				)
			`, qID, patID, invNum, it.regFee, it.payMethod, now.Add(-time.Duration(15-idx*5)*time.Minute))
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// SeedRushHour seeds 50 realistic queues across Kidinest (30) & SSC (20).
func SeedRushHour(db *sql.DB) error {
	if err := WipeTransactions(db); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// 30 Pediatric Patients (Kidinest)
	kidinestPatients := []struct {
		name     string
		phone    string
		cat      string
		doc      string
		qNum     string
		status   string
		payMeth  string
		fee      int
		diag     string
		action   string
		actFee   int
		rxFee    int
	}{
		// COMPLETED (14)
		{"Kenzo Alvaro", "08121111001", "A_VAKSIN", "dr. Hendra Pratama, Sp.A", "A-001", "COMPLETED", "QRIS", 35000, "Imunisasi DTP-3 Sesuai Jadwal", "Injeksi Vaksin DTP Hexaxim paha kiri", 100000, 850000},
		{"Aisyah Putri", "08121111002", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-001", "COMPLETED", "CASH", 35000, "ISPA Ringan / Common Cold", "Pemeriksaan THT & Nebulisasi NaCl", 75000, 60000},
		{"Bima Sakti", "08121111003", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-002", "COMPLETED", "CASH", 35000, "Faringitis Akut Febris", "Pemeriksaan fisik & edukasi hidrasi", 50000, 70000},
		{"Naura Zahrani", "08121111004", "C_TUMBANG", "dr. Hendra Pratama, Sp.A", "C-001", "COMPLETED", "EDC", 50000, "Evaluasi Speech Delay & Motorik Kasar", "Skrining Denver II & Konseling Stimulasi", 250000, 0},
		{"Arsya Malik", "08121111005", "A_VAKSIN", "dr. Patricia, Sp.A", "A-002", "COMPLETED", "BANK_TRANSFER", 35000, "Vaksin Rotavirus Dosis 2", "Pemberian Vaksin Rotavirus Oral", 100000, 350000},
		{"Clarissa Aurelia", "08121111006", "B_SAKIT", "dr. Patricia, Sp.A", "B-003", "COMPLETED", "QRIS", 35000, "Bronkitis Akut Pediatri", "Nebulisasi Salbutamol & Fisioterapi Dada", 120000, 95000},
		{"Raffi Athariz", "08121111007", "D_PAKET", "dr. Patricia, Sp.A", "D-001", "COMPLETED", "EDC", 75000, "Paket Skrining Alergi & Nutrisi Emas", "Tes Alergi Kulit & Cek Darah Rutin", 350000, 120000},
		{"Zahra Amira", "08121111008", "A_VAKSIN", "dr. Hendra Pratama, Sp.A", "A-003", "COMPLETED", "QRIS", 35000, "Vaksin PCV 13 Dosis 1", "Injeksi Intramuskular Paha Kanan", 100000, 750000},
		{"Dhafin Alfarizi", "08121111009", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-004", "COMPLETED", "CASH", 35000, "Gastroenteritis Akut Dehidrasi Ringan", "Edukasi Oralit & Terapi Zinc 10 Hari", 60000, 63000},
		{"Mikayla Bilqis", "08121111010", "C_TUMBANG", "dr. Hendra Pratama, Sp.A", "C-002", "COMPLETED", "QRIS", 50000, "Gagal Tumbuh (Failure to Thrive) & Picky Eater", "Pemeriksaan Antropometri & Rekomendasi Pangan FMS", 200000, 150000},
		{"Fathir Arkananta", "08121111011", "A_VAKSIN", "dr. Patricia, Sp.A", "A-004", "COMPLETED", "BANK_TRANSFER", 35000, "Vaksin Influenza Tahunan Pediatri", "Injeksi Deltoid Subkutan", 100000, 320000},
		{"Kayla Shaqueena", "08121111012", "B_SAKIT", "dr. Patricia, Sp.A", "B-005", "COMPLETED", "CASH", 35000, "Otitis Media Akut Dextra", "Otoskopi & Pembersihan Liang Telinga Luar", 85000, 75000},
		{"Rayyanza Malik", "08121111013", "D_PAKET", "dr. Patricia, Sp.A", "D-002", "COMPLETED", "EDC", 75000, "Paket Imunisasi Dasar Lengkap 9 Bulan", "Injeksi Campak Rubella & Vitamin A", 250000, 200000},
		{"Alana Zivanna", "08121111014", "E_WALKIN", "dr. Hendra Pratama, Sp.A", "E-001", "COMPLETED", "QRIS", 35000, "Alergi Makanan Akut (Urtikaria)", "Injeksi Deksametason & Edukasi Eliminasi Diet", 75000, 45000},

		// PHARMACY_WAITING (5)
		{"Gibran Rakabuming", "08121111015", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-006", "PHARMACY_WAITING", "QRIS", 35000, "Tonsilofaringitis Akut", "Eksudat Tonsil T2-T2, Resep Amoxicillin & Paracetamol", 65000, 80000},
		{"Jan Ethes Srinarendra", "08121111016", "A_VAKSIN", "dr. Hendra Pratama, Sp.A", "A-005", "PHARMACY_WAITING", "CASH", 35000, "Vaksin MMR & Varicella", "Observasi Pasca Injeksi, Resep Penurun Panas Profilaksis", 100000, 40000},
		{"Sedah Mirah Nasution", "08121111017", "B_SAKIT", "dr. Patricia, Sp.A", "B-007", "PHARMACY_WAITING", "BANK_TRANSFER", 35000, "Dermatitis Atopik Pediatri Eksaserbasi", "Resep Salep Racikan Hidrokortison & Cetirizine Drop", 75000, 110000},
		{"Panembahan Al Nahyan", "08121111018", "C_TUMBANG", "dr. Patricia, Sp.A", "C-003", "PHARMACY_WAITING", "EDC", 50000, "Perawakan Pendek (Stunting) & Anemia Defisiensi Besi", "Resep Zat Besi Tetes Maltofer & Multivitamin Zinc", 150000, 145000},
		{"La Lembah Manah", "08121111019", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-008", "PHARMACY_WAITING", "QRIS", 35000, "Asma Bronkial Eksaserbasi Ringan", "Nebulisasi Berotec & Resep Sirup Salbutamol", 95000, 65000},

		// IN_EXAM (2)
		{"Arsy Addara", "08121111020", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-009", "IN_EXAM", "CASH", 35000, "", "", 0, 0},
		{"Kiano Tiger Wong", "08121111021", "A_VAKSIN", "dr. Patricia, Sp.A", "A-006", "IN_EXAM", "QRIS", 35000, "", "", 0, 0},

		// CALLING (1)
		{"Kenzo Eldrago", "08121111022", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-010", "CALLING", "QRIS", 35000, "", "", 0, 0},

		// WAITING (8)
		{"Ameena Hanna", "08121111023", "B_SAKIT", "dr. Hendra Pratama, Sp.A", "B-011", "WAITING", "CASH", 35000, "", "", 0, 0},
		{"Azzaura Nur", "08121111024", "A_VAKSIN", "dr. Patricia, Sp.A", "A-007", "WAITING", "QRIS", 35000, "", "", 0, 0},
		{"Adzam Adriansyah", "08121111025", "B_SAKIT", "dr. Patricia, Sp.A", "B-012", "WAITING", "BANK_TRANSFER", 35000, "", "", 0, 0},
		{"Shaquille Kaili", "08121111026", "C_TUMBANG", "dr. Hendra Pratama, Sp.A", "C-004", "WAITING", "EDC", 50000, "", "", 0, 0},
		{"Gemma Nur", "08121111027", "A_VAKSIN", "dr. Hendra Pratama, Sp.A", "A-008", "WAITING", "QRIS", 35000, "", "", 0, 0},
		{"Guevara Rayyan", "08121111028", "D_PAKET", "dr. Patricia, Sp.A", "D-003", "WAITING", "EDC", 75000, "", "", 0, 0},
		{"Zunaira Alessia", "08121111029", "C_TUMBANG", "dr. Patricia, Sp.A", "C-005", "WAITING", "CASH", 50000, "", "", 0, 0},
		{"Kalandra Zayan", "08121111030", "E_WALKIN", "dr. Hendra Pratama, Sp.A", "E-002", "WAITING", "QRIS", 35000, "", "", 0, 0},
	}

	// 20 Dermatology / Esthetics Patients (Semarang Skin Center)
	sscPatients := []struct {
		name     string
		phone    string
		cat      string
		doc      string
		qNum     string
		status   string
		payMeth  string
		fee      int
		diag     string
		action   string
		actFee   int
		rxFee    int
	}{
		// COMPLETED (8)
		{"Nadia Saphira", "08122222001", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-001", "COMPLETED", "QRIS", 50000, "Acne Vulgaris Papulopustular Grade III", "Ekstraksi Komedo, IPL Rejuvenation & Masker Peel-Off", 450000, 215000},
		{"Jessica Mila", "08122222002", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-002", "COMPLETED", "EDC", 50000, "Melasma Centrofascial & Hiperpigmentasi Pasca Inflamasi", "Laser Nd:YAG Q-Switched Toning Wajah", 850000, 260000},
		{"Kevin Sanjaya", "08122222003", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-003", "COMPLETED", "BANK_TRANSFER", 50000, "Veruka Vulgaris Digiti II Manus Dextra", "Kauterisasi Elektrokauter Bedah Minor", 350000, 38000},
		{"Laura Basuki", "08122222004", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-004", "COMPLETED", "EDC", 50000, "Photoaging & Kerut Halus Periorbita", "Skinbooster Hyaluronic Acid & Vitamin Infusion", 1250000, 140000},
		{"Nicholas Saputra", "08122222005", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-005", "COMPLETED", "QRIS", 50000, "Dermatitis Seboroik Regio Scalp & T-Zone", "Kortikosteroid Topikal & Shampoo Ketoconazole 2%", 150000, 113000},
		{"Dian Sastrowardoyo", "08122222006", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-006", "COMPLETED", "EDC", 50000, "Aging Skin & Loss of Elasticity", "HIFU Full Face & Double Chin Lifting", 1500000, 335000},
		{"Chicco Jerikho", "08122222007", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-007", "COMPLETED", "CASH", 50000, "Folikulitis Barbae Akut", "Eksisi Pustul & Kompres Antiseptik Terbuka", 200000, 95000},
		{"Pevita Pearce", "08122222008", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-008", "COMPLETED", "QRIS", 50000, "Scar Atrofi Acne (Bopeng Ice-pick)", "Microneedling RF & Serum PDRN Salmon DNA", 950000, 215000},

		// PHARMACY_WAITING (3)
		{"Maudy Ayunda", "08122222009", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-009", "PHARMACY_WAITING", "QRIS", 50000, "Acne Comedonica & Seborrhea Berat", "Facial Medical Deep Cleansing & Resep Paket Krim SSC", 350000, 290000},
		{"Reza Rahadian", "08122222010", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-010", "PHARMACY_WAITING", "EDC", 50000, "Tinea Cruris & Korporis Kronik", "Pemeriksaan Lampu Wood & Resep Antifungal Oral/Salep", 150000, 120000},
		{"Chelsea Islan", "08122222011", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-011", "PHARMACY_WAITING", "BANK_TRANSFER", 50000, "Rosacea Eritematotelangiektatik", "Laser V-Beam Vascular Therapy & Resep Krim Soothing", 750000, 175000},

		// IN_EXAM (2)
		{"Adinia Wirasti", "08122222012", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-012", "IN_EXAM", "QRIS", 50000, "", "", 0, 0},
		{"Iko Uwais", "08122222013", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-013", "IN_EXAM", "CASH", 50000, "", "", 0, 0},

		// CALLING (1)
		{"Tara Basro", "08122222014", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-014", "CALLING", "QRIS", 50000, "", "", 0, 0},

		// WAITING (6)
		{"Joe Taslim", "08122222015", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-015", "WAITING", "EDC", 50000, "", "", 0, 0},
		{"Velove Vexia", "08122222016", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-016", "WAITING", "QRIS", 50000, "", "", 0, 0},
		{"Hamish Daud", "08122222017", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-017", "WAITING", "BANK_TRANSFER", 50000, "", "", 0, 0},
		{"Tatjana Saphira", "08122222018", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-018", "WAITING", "CASH", 50000, "", "", 0, 0},
		{"Rio Dewanto", "08122222019", "SSC_TREATMENT", "dr. Kevin Wijaya, Sp.D.V.E", "SSC-019", "WAITING", "EDC", 50000, "", "", 0, 0},
		{"Marsha Timothy", "08122222020", "SSC_TREATMENT", "dr. Dian Paramita, Sp.D.V.E", "SSC-020", "WAITING", "QRIS", 50000, "", "", 0, 0},
	}

	// Insert Kidinest
	for i, p := range kidinestPatients {
		if err := insertQueueItem(tx, p.name, p.phone, "KIDINEST", p.cat, p.doc, p.qNum, p.status, p.payMeth, p.fee, p.diag, p.action, p.actFee, p.rxFee, i+1, now); err != nil {
			return fmt.Errorf("gagal insert rush kidinest (%s): %w", p.name, err)
		}
	}

	// Insert SSC
	for i, p := range sscPatients {
		if err := insertQueueItem(tx, p.name, p.phone, "SSC", p.cat, p.doc, p.qNum, p.status, p.payMeth, p.fee, p.diag, p.action, p.actFee, p.rxFee, 31+i, now); err != nil {
			return fmt.Errorf("gagal insert rush ssc (%s): %w", p.name, err)
		}
	}

	return tx.Commit()
}

func insertQueueItem(
	tx *sql.Tx,
	name, phone, unit, cat, doc, qNum, status, payMeth string,
	regFee int,
	diag, action string,
	actFee, rxFee int,
	seq int,
	baseTime time.Time,
) error {
	createdAt := baseTime.Add(-time.Duration(180-seq*3) * time.Minute)
	paidAt := createdAt.Add(3 * time.Minute)
	calledAt := createdAt.Add(15 * time.Minute)
	completedAt := createdAt.Add(35 * time.Minute)

	var patID int
	err := tx.QueryRow(`
		INSERT INTO patients (name, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $3) RETURNING id
	`, name, phone, createdAt).Scan(&patID)
	if err != nil {
		return err
	}

	var qID int
	err = tx.QueryRow(`
		INSERT INTO queues (
			patient_id, clinic_unit, category, doctor_name, queue_number,
			registration_fee, payment_method, payment_status, queue_status,
			paid_at, called_at, completed_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, 'PAID', $8,
			$9, $10, $11, $12
		) RETURNING id
	`, patID, unit, cat, doc, qNum,
		regFee, payMeth, status,
		paidAt, calledAt, completedAt, createdAt,
	).Scan(&qID)
	if err != nil {
		return err
	}

	// If COMPLETED or PHARMACY_WAITING, create Consultation
	if status == "COMPLETED" || status == "PHARMACY_WAITING" {
		var consID int
		err = tx.QueryRow(`
			INSERT INTO consultations (
				queue_id, doctor_name, anamnesis, diagnosis, icd10_code,
				action_notes, action_fee, status, created_at, updated_at
			) VALUES (
				$1, $2, 'Keluhan utama pasien terdata di triase awal', $3, 'Z00.0',
				$4, $5, 'COMPLETED', $6, $6
			) RETURNING id
		`, qID, doc, diag, action, actFee, calledAt).Scan(&consID)
		if err != nil {
			return err
		}

		// Create Prescription
		var rxID int
		rxStatus := "DISPENSED"
		if status == "PHARMACY_WAITING" {
			rxStatus = "RESERVED"
		}
		err = tx.QueryRow(`
			INSERT INTO prescriptions (
				consultation_id, queue_id, status, notes, created_at, updated_at
			) VALUES (
				$1, $2, $3, 'Paket resep klinis terintegrasi', $4, $4
			) RETURNING id
		`, consID, qID, rxStatus, calledAt.Add(5*time.Minute)).Scan(&rxID)
		if err != nil {
			return err
		}

		// Insert 1 Prescription item
		_, err = tx.Exec(`
			INSERT INTO prescription_items (
				prescription_id, item_name, is_manual, quantity, unit_price, subtotal, dosage_instruction, created_at
			) VALUES (
				$1, 'Obat Terapi Standar Klinis', FALSE, 1, $2, $2, '3x1 sesudah makan', $3
			)
		`, rxID, rxFee, calledAt.Add(6*time.Minute))
		if err != nil {
			return err
		}
	}

	// If COMPLETED, generate settled Invoice
	if status == "COMPLETED" {
		invNum := fmt.Sprintf("INV-%s-%04d", baseTime.Format("20060102"), seq)
		totalAmount := regFee + actFee + rxFee
		cashierName := "Siti Sarah, A.Md"
		if unit == "SSC" {
			cashierName = "Amanda L., S.E"
		}

		_, err = tx.Exec(`
			INSERT INTO invoices (
				queue_id, patient_id, invoice_number,
				subtotal_registration, subtotal_action, subtotal_medication,
				voucher_discount, total_amount, payment_method,
				cash_received, cash_change, cashier_name, status, created_at, updated_at
			) VALUES (
				$1, $2, $3,
				$4, $5, $6,
				0, $7, $8,
				$7, 0, $9, 'PAID', $10, $10
			)
		`, qID, patID, invNum, regFee, actFee, rxFee, totalAmount, payMeth, cashierName, completedAt)
		if err != nil {
			return err
		}
	}

	return nil
}
