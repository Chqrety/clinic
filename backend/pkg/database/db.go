package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func InitDB(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create db directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite db: %w", err)
	}

	if err := migrateSchema(db); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	if err := seedInitialData(db); err != nil {
		log.Printf("[DB] Warning: seeding failed: %v", err)
	}

	return db, nil
}

func migrateSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS patients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nik TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		birth_date TEXT,
		gender TEXT,
		ihs_number TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS queues (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		patient_id INTEGER NOT NULL REFERENCES patients(id),
		queue_number TEXT NOT NULL,
		complaint TEXT,
		payment_method TEXT NOT NULL,
		payment_status TEXT NOT NULL,
		queue_status TEXT NOT NULL,
		paid_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS drugs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		kfa_code TEXT NOT NULL,
		unit TEXT NOT NULL,
		category TEXT NOT NULL,
		min_stock_alert INTEGER DEFAULT 10,
		sell_price INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS drug_batches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		drug_id INTEGER NOT NULL REFERENCES drugs(id),
		batch_number TEXT NOT NULL,
		expired_date TEXT NOT NULL,
		physical_stock INTEGER NOT NULL DEFAULT 0,
		reserved_stock INTEGER NOT NULL DEFAULT 0,
		buy_price INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS consultations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		queue_id INTEGER NOT NULL REFERENCES queues(id),
		doctor_name TEXT NOT NULL DEFAULT 'dr. Hendra Pratama, Sp.PD',
		anamnesis TEXT,
		diagnosis TEXT,
		icd10_code TEXT,
		status TEXT NOT NULL DEFAULT 'DRAFT',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS prescriptions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		consultation_id INTEGER NOT NULL REFERENCES consultations(id),
		status TEXT NOT NULL DEFAULT 'RESERVED',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS prescription_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		prescription_id INTEGER NOT NULL REFERENCES prescriptions(id),
		drug_id INTEGER NOT NULL REFERENCES drugs(id),
		batch_id INTEGER REFERENCES drug_batches(id),
		quantity INTEGER NOT NULL,
		dosage_instruction TEXT NOT NULL,
		unit_price INTEGER NOT NULL,
		subtotal INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS invoices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		queue_id INTEGER NOT NULL REFERENCES queues(id),
		invoice_number TEXT NOT NULL UNIQUE,
		registration_fee INTEGER DEFAULT 0,
		consultation_fee INTEGER DEFAULT 0,
		medication_fee INTEGER DEFAULT 0,
		total_amount INTEGER NOT NULL,
		payment_method TEXT NOT NULL,
		payment_status TEXT NOT NULL DEFAULT 'PAID',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS satusehat_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		resource_type TEXT NOT NULL,
		entity_id INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'READY',
		fhir_payload TEXT NOT NULL,
		synced_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(schema)
	return err
}

func seedInitialData(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM drugs").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // already seeded
	}

	log.Println("[DB] Seeding realistic clinic data...")

	// 1. Seed Drugs
	type seedDrug struct {
		name        string
		kfa         string
		unit        string
		cat         string
		minAlert    int
		sellPrice   int
		batch       string
		edMonths    int // months from today
		physStock   int
		buyPrice    int
	}

	now := time.Now()
	drugs := []seedDrug{
		{
			name:      "Paracetamol 500 mg Tablet",
			kfa:       "93001029",
			unit:      "Tablet",
			cat:       "Analgesik / Antipiretik",
			minAlert:  50,
			sellPrice: 1500,
			batch:     "BCH-2026-PCT-01",
			edMonths:  18, // aman
			physStock: 250,
			buyPrice:  800,
		},
		{
			name:      "Amoxicillin 500 mg Kapsul",
			kfa:       "93000412",
			unit:      "Kapsul",
			cat:       "Antibiotik",
			minAlert:  30,
			sellPrice: 3000,
			batch:     "BCH-2026-AMX-02",
			edMonths:  2, // mendekati expired (< 3 bulan)!
			physStock: 45,
			buyPrice:  1800,
		},
		{
			name:      "Cetirizine 10 mg Tablet",
			kfa:       "93003881",
			unit:      "Tablet",
			cat:       "Antihistamin",
			minAlert:  25,
			sellPrice: 2500,
			batch:     "BCH-2026-CTZ-01",
			edMonths:  12,
			physStock: 80,
			buyPrice:  1200,
		},
		{
			name:      "Metformin HCl 500 mg",
			kfa:       "93005912",
			unit:      "Tablet",
			cat:       "Antidiabetes",
			minAlert:  30,
			sellPrice: 2000,
			batch:     "BCH-2026-MET-03",
			edMonths:  24,
			physStock: 12, // stok menipis (< minAlert)!
			buyPrice:  1100,
		},
		{
			name:      "Omeprazole 20 mg Kapsul",
			kfa:       "93007123",
			unit:      "Kapsul",
			cat:       "Antasida / Lambung",
			minAlert:  20,
			sellPrice: 4000,
			batch:     "BCH-2026-OMP-01",
			edMonths:  14,
			physStock: 120,
			buyPrice:  2200,
		},
		{
			name:      "Vitamin C 500 mg",
			kfa:       "93009901",
			unit:      "Tablet",
			cat:       "Suplemen & Vitamin",
			minAlert:  40,
			sellPrice: 1000,
			batch:     "BCH-2026-VTC-05",
			edMonths:  20,
			physStock: 300,
			buyPrice:  500,
		},
	}

	for _, d := range drugs {
		res, err := db.Exec(`
			INSERT INTO drugs (name, kfa_code, unit, category, min_stock_alert, sell_price)
			VALUES (?, ?, ?, ?, ?, ?)
		`, d.name, d.kfa, d.unit, d.cat, d.minAlert, d.sellPrice)
		if err != nil {
			return err
		}
		drugID, _ := res.LastInsertId()

		edDate := now.AddDate(0, d.edMonths, 0).Format("2006-01-02")
		_, err = db.Exec(`
			INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, buy_price)
			VALUES (?, ?, ?, ?, 0, ?)
		`, drugID, d.batch, edDate, d.physStock, d.buyPrice)
		if err != nil {
			return err
		}
	}

	// 2. Seed Sample Patients & Demo Queues
	patients := []struct {
		nik     string
		name    string
		phone   string
		birth   string
		gender  string
		ihs     string
		qNum    string
		qComp   string
		qPayM   string
		qPayS   string
		qStat   string
	}{
		{
			nik:    "3171012345670001",
			name:   "Budi Santoso",
			phone:  "081234567890",
			birth:  "1988-04-12",
			gender: "Laki-laki",
			ihs:    "P01928374611",
			qNum:   "A-001",
			qComp:  "Demam tinggi dan batuk 3 hari",
			qPayM:  "QRIS",
			qPayS:  "PAID",
			qStat:  "CALLING",
		},
		{
			nik:    "3273019876540002",
			name:   "Siti Rahmawati",
			phone:  "085678901234",
			birth:  "1994-09-21",
			gender: "Perempuan",
			ihs:    "P01928374612",
			qNum:   "A-002",
			qComp:  "Nyeri lambung dan mual sejak kemarin",
			qPayM:  "QRIS",
			qPayS:  "PAID",
			qStat:  "WAITING",
		},
		{
			nik:    "3578014567890003",
			name:   "Ahmad Fauzi",
			phone:  "087890123456",
			birth:  "1976-11-05",
			gender: "Laki-laki",
			ihs:    "P01928374613",
			qNum:   "A-003",
			qComp:  "Pusing dan flu berat",
			qPayM:  "CASH",
			qPayS:  "PENDING", // Menunggu kasir terima tunai
			qStat:  "WAITING",
		},
	}

	for _, p := range patients {
		res, err := db.Exec(`
			INSERT INTO patients (nik, name, phone, birth_date, gender, ihs_number)
			VALUES (?, ?, ?, ?, ?, ?)
		`, p.nik, p.name, p.phone, p.birth, p.gender, p.ihs)
		if err != nil {
			return err
		}
		pID, _ := res.LastInsertId()

		var paidAt interface{} = nil
		if p.qPayS == "PAID" {
			paidAt = time.Now()
		}

		_, err = db.Exec(`
			INSERT INTO queues (patient_id, queue_number, complaint, payment_method, payment_status, queue_status, paid_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, pID, p.qNum, p.qComp, p.qPayM, p.qPayS, p.qStat, paidAt)
		if err != nil {
			return err
		}
	}

	log.Println("[DB] Seeding completed successfully.")
	return nil
}
