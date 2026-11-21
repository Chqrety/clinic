package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() Config {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "clinic"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "clinic_secret"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "clinic_db"
	}
	ssl := os.Getenv("DB_SSLMODE")
	if ssl == "" {
		ssl = "disable"
	}

	return Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pass,
		DBName:   name,
		SSLMode:  ssl,
	}
}

func InitPostgres(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database pgx: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping database postgresql (%s:%s): %w", cfg.Host, cfg.Port, err)
	}

	if err := MigrateStaffSchema(db); err != nil {
		return nil, fmt.Errorf("gagal migrasi skema staff: %w", err)
	}

	if err := MigrateQueueSchema(db); err != nil {
		return nil, fmt.Errorf("gagal migrasi skema queue: %w", err)
	}

	if err := MigrateCashierSchema(db); err != nil {
		return nil, fmt.Errorf("gagal migrasi skema cashier: %w", err)
	}

	if err := MigrateConsultationSchema(db); err != nil {
		return nil, fmt.Errorf("gagal migrasi skema consultation: %w", err)
	}

	if err := MigrateDoctorScheduleSchema(db); err != nil {
		return nil, fmt.Errorf("gagal migrasi skema jadwal dokter: %w", err)
	}

	if err := SeedInitialStaff(db); err != nil {
		log.Printf("[DB] Peringatan: Seeding data staf gagal: %v", err)
	}

	if err := MigratePayrollSchema(db); err != nil {
		log.Printf("[DB] Peringatan: Gagal migrasi skema payroll: %v", err)
	}

	if err := SeedInitialPayrollEmployees(db); err != nil {
		log.Printf("[DB] Peringatan: Seeding data payroll karyawan gagal: %v", err)
	}

	return db, nil
}

func MigrateStaffSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS staff_users (
		id SERIAL PRIMARY KEY,
		employee_number VARCHAR(50) UNIQUE NOT NULL,
		name VARCHAR(150) NOT NULL,
		role VARCHAR(30) NOT NULL CHECK (role IN ('SUPERADMIN_OWNER', 'DOCTOR', 'GRO_CASHIER', 'PHARMACIST', 'FINANCE')),
		clinic_unit VARCHAR(20) NOT NULL CHECK (clinic_unit IN ('KIDINEST', 'SSC', 'ALL')),
		password_hash VARCHAR(255) NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_staff_employee_number ON staff_users(employee_number);
	CREATE INDEX IF NOT EXISTS idx_staff_role ON staff_users(role);
	CREATE INDEX IF NOT EXISTS idx_staff_unit ON staff_users(clinic_unit);
	`
	_, err := db.Exec(schema)
	return err
}

func MigrateQueueSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS patients (
		id SERIAL PRIMARY KEY,
		name VARCHAR(150) NOT NULL,
		phone VARCHAR(30) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_patients_phone ON patients(phone);

	CREATE TABLE IF NOT EXISTS daily_queue_counters (
		counter_date DATE NOT NULL,
		clinic_unit VARCHAR(20) NOT NULL,
		category_prefix VARCHAR(5) NOT NULL,
		last_seq INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (counter_date, clinic_unit, category_prefix)
	);

	CREATE TABLE IF NOT EXISTS queues (
		id SERIAL PRIMARY KEY,
		patient_id INTEGER NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
		clinic_unit VARCHAR(20) NOT NULL CHECK (clinic_unit IN ('KIDINEST', 'SSC')),
		category VARCHAR(30) NOT NULL CHECK (category IN (
			'A_VAKSIN', 'B_SAKIT', 'C_TUMBANG', 'D_PAKET', 'E_WALKIN', 'SSC_TREATMENT'
		)),
		doctor_id INTEGER REFERENCES staff_users(id),
		doctor_name VARCHAR(150) NOT NULL,
		queue_number VARCHAR(20) NOT NULL,
		registration_fee INTEGER NOT NULL DEFAULT 0,
		payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN (
			'QRIS', 'BANK_TRANSFER', 'CASH', 'UNPAID_BOOKING', 'EDC'
		)),
		payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('PENDING', 'PAID', 'REFUNDED')),
		queue_status VARCHAR(30) NOT NULL CHECK (queue_status IN (
			'BOOKED', 'WAITING', 'CALLING', 'IN_EXAM', 'PHARMACY_WAITING', 'COMPLETED', 'CANCELLED'
		)),
		paid_at TIMESTAMP WITH TIME ZONE,
		called_at TIMESTAMP WITH TIME ZONE,
		completed_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_queues_unit_status ON queues(clinic_unit, queue_status);
	CREATE INDEX IF NOT EXISTS idx_queues_created_at ON queues(created_at);

	ALTER TABLE queues DROP CONSTRAINT IF EXISTS queues_payment_method_check;
	ALTER TABLE queues ADD CONSTRAINT queues_payment_method_check CHECK (payment_method IN (
		'QRIS', 'BANK_TRANSFER', 'CASH', 'UNPAID_BOOKING', 'EDC'
	));

	ALTER TABLE queues DROP CONSTRAINT IF EXISTS queues_queue_status_check;
	ALTER TABLE queues ADD CONSTRAINT queues_queue_status_check CHECK (queue_status IN (
		'BOOKED', 'WAITING', 'CALLING', 'IN_EXAM', 'PHARMACY_WAITING', 'COMPLETED', 'CANCELLED'
	));

	ALTER TABLE queues ADD COLUMN IF NOT EXISTS payment_reference VARCHAR(100);
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS payment_url TEXT;
	CREATE INDEX IF NOT EXISTS idx_queues_payment_ref ON queues(payment_reference);

	ALTER TABLE patients ADD COLUMN IF NOT EXISTS nik VARCHAR(20);
	ALTER TABLE patients ADD COLUMN IF NOT EXISTS ihs_number VARCHAR(50);
	ALTER TABLE patients ADD COLUMN IF NOT EXISTS birth_date DATE;
	ALTER TABLE patients ADD COLUMN IF NOT EXISTS gender VARCHAR(10);
	CREATE INDEX IF NOT EXISTS idx_patients_nik ON patients(nik);

	CREATE TABLE IF NOT EXISTS satusehat_logs (
		id SERIAL PRIMARY KEY,
		resource_type VARCHAR(50) NOT NULL,
		entity_id VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL,
		fhir_payload TEXT NOT NULL,
		synced_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_satusehat_logs_status ON satusehat_logs(status);
	`
	_, err := db.Exec(schema)
	return err
}

func MigrateCashierSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS invoices (
		id SERIAL PRIMARY KEY,
		queue_id INTEGER NOT NULL REFERENCES queues(id) ON DELETE RESTRICT,
		patient_id INTEGER NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
		invoice_number VARCHAR(50) UNIQUE NOT NULL,
		subtotal_registration INTEGER NOT NULL DEFAULT 0,
		subtotal_action INTEGER NOT NULL DEFAULT 0,
		subtotal_medication INTEGER NOT NULL DEFAULT 0,
		voucher_discount INTEGER NOT NULL DEFAULT 0,
		total_amount INTEGER NOT NULL DEFAULT 0,
		payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN ('CASH', 'BANK_TRANSFER', 'EDC', 'QRIS')),
		cash_received INTEGER NOT NULL DEFAULT 0,
		cash_change INTEGER NOT NULL DEFAULT 0,
		payment_reference VARCHAR(100),
		cashier_staff_id INTEGER REFERENCES staff_users(id),
		cashier_name VARCHAR(150) NOT NULL,
		status VARCHAR(20) NOT NULL CHECK (status IN ('PAID', 'REFUNDED')),
		refund_bank_name VARCHAR(50),
		refund_bank_account VARCHAR(50),
		refund_account_name VARCHAR(150),
		refund_reason TEXT,
		refunded_at TIMESTAMP WITH TIME ZONE,
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_invoices_queue_id ON invoices(queue_id);
	CREATE INDEX IF NOT EXISTS idx_invoices_created_at ON invoices(created_at);

	CREATE TABLE IF NOT EXISTS cashier_shift_reports (
		id SERIAL PRIMARY KEY,
		staff_id INTEGER NOT NULL REFERENCES staff_users(id),
		staff_name VARCHAR(150) NOT NULL,
		shift_date DATE NOT NULL,
		total_transactions INTEGER NOT NULL DEFAULT 0,
		total_cash INTEGER NOT NULL DEFAULT 0,
		total_transfer INTEGER NOT NULL DEFAULT 0,
		total_edc INTEGER NOT NULL DEFAULT 0,
		total_voucher INTEGER NOT NULL DEFAULT 0,
		total_refund INTEGER NOT NULL DEFAULT 0,
		net_total INTEGER NOT NULL DEFAULT 0,
		notes TEXT,
		sent_to_owner_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_shift_reports_date ON cashier_shift_reports(shift_date);
	`
	_, err := db.Exec(schema)
	return err
}

func SeedInitialStaff(db *sql.DB) error {
	// Default demo password for all staff: "klinik123"
	defaultHash, err := bcrypt.GenerateFromPassword([]byte("klinik123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password default: %w", err)
	}

	seeds := []struct {
		empNum string
		name   string
		role   string
		unit   string
	}{
		{"OWNER-01", "dr. Ratna Paramita, M.Kes", "SUPERADMIN_OWNER", "ALL"},
		{"DOC-KID-01", "dr. Patricia Vanessa A., Sp.A", "DOCTOR", "KIDINEST"},
		{"DOC-KID-02", "dr. Hendra Pratama, Sp.A", "DOCTOR", "KIDINEST"},
		{"DOC-SSC-01", "dr. Dian Paramita, Sp.D.V.E", "DOCTOR", "SSC"},
		{"DOC-SSC-02", "dr. Kevin Wijaya, Sp.D.V.E", "DOCTOR", "SSC"},
		{"GRO-01", "Siti Sarah, A.Md", "GRO_CASHIER", "ALL"},
		{"FARMASI-01", "Apt. Dimas Wicaksono, S.Farm", "PHARMACIST", "ALL"},
		{"FIN-01", "Budi Santoso, SE, Ak., CA", "FINANCE", "ALL"},
	}

	query := `
	INSERT INTO staff_users (employee_number, name, role, clinic_unit, password_hash, is_active)
	VALUES ($1, $2, $3, $4, $5, TRUE)
	ON CONFLICT (employee_number) DO UPDATE SET
		name = EXCLUDED.name,
		clinic_unit = EXCLUDED.clinic_unit;
	`

	for _, s := range seeds {
		if _, err := db.Exec(query, s.empNum, s.name, s.role, s.unit, string(defaultHash)); err != nil {
			return fmt.Errorf("gagal insert seed staf %s: %w", s.empNum, err)
		}
	}

	log.Printf("[DB] Berhasil menyelaraskan 7 akun staf demo (4 Dokter Spesialis Kidinest & SSC).")
	return nil
}

func MigrateConsultationSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS drugs (
		id SERIAL PRIMARY KEY,
		clinic_unit VARCHAR(20) NOT NULL CHECK (clinic_unit IN ('KIDINEST', 'SSC', 'ALL')),
		name VARCHAR(150) NOT NULL,
		kfa_code VARCHAR(50),
		unit VARCHAR(30) NOT NULL,
		category VARCHAR(50) NOT NULL,
		min_stock_alert INTEGER NOT NULL DEFAULT 10,
		buy_price INTEGER NOT NULL DEFAULT 0,
		sell_price INTEGER NOT NULL DEFAULT 0,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_drugs_unit ON drugs(clinic_unit);

	CREATE TABLE IF NOT EXISTS drug_batches (
		id SERIAL PRIMARY KEY,
		drug_id INTEGER NOT NULL REFERENCES drugs(id) ON DELETE CASCADE,
		batch_number VARCHAR(50) NOT NULL,
		expired_date DATE NOT NULL,
		physical_stock INTEGER NOT NULL DEFAULT 0,
		reserved_stock INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_batches_drug_id ON drug_batches(drug_id);

	CREATE TABLE IF NOT EXISTS consultations (
		id SERIAL PRIMARY KEY,
		queue_id INTEGER NOT NULL REFERENCES queues(id) ON DELETE RESTRICT,
		doctor_id INTEGER REFERENCES staff_users(id),
		doctor_name VARCHAR(150) NOT NULL,
		anamnesis TEXT,
		diagnosis TEXT,
		icd10_code VARCHAR(50),
		action_notes TEXT,
		action_fee INTEGER NOT NULL DEFAULT 0,
		status VARCHAR(30) NOT NULL DEFAULT 'COMPLETED',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_consultations_queue_id ON consultations(queue_id);

	CREATE TABLE IF NOT EXISTS prescriptions (
		id SERIAL PRIMARY KEY,
		consultation_id INTEGER NOT NULL REFERENCES consultations(id) ON DELETE CASCADE,
		queue_id INTEGER NOT NULL REFERENCES queues(id) ON DELETE RESTRICT,
		status VARCHAR(30) NOT NULL DEFAULT 'RESERVED',
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_prescriptions_queue_id ON prescriptions(queue_id);

	CREATE TABLE IF NOT EXISTS prescription_items (
		id SERIAL PRIMARY KEY,
		prescription_id INTEGER NOT NULL REFERENCES prescriptions(id) ON DELETE CASCADE,
		drug_id INTEGER REFERENCES drugs(id),
		batch_id INTEGER REFERENCES drug_batches(id) ON DELETE SET NULL,
		item_name VARCHAR(150) NOT NULL,
		is_manual BOOLEAN NOT NULL DEFAULT FALSE,
		quantity INTEGER NOT NULL DEFAULT 1,
		unit_price INTEGER NOT NULL DEFAULT 0,
		subtotal INTEGER NOT NULL DEFAULT 0,
		dosage_instruction VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_prescription_items_rx_id ON prescription_items(prescription_id);
	ALTER TABLE prescription_items ADD COLUMN IF NOT EXISTS batch_id INTEGER REFERENCES drug_batches(id) ON DELETE SET NULL;

	CREATE TABLE IF NOT EXISTS stock_opname_logs (
		id SERIAL PRIMARY KEY,
		drug_id INTEGER NOT NULL REFERENCES drugs(id) ON DELETE RESTRICT,
		batch_id INTEGER REFERENCES drug_batches(id) ON DELETE SET NULL,
		system_stock INTEGER NOT NULL,
		actual_stock INTEGER NOT NULL,
		discrepancy INTEGER NOT NULL,
		adjusted_by VARCHAR(150) NOT NULL,
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_stock_opname_drug_id ON stock_opname_logs(drug_id);
	CREATE INDEX IF NOT EXISTS idx_stock_opname_created_at ON stock_opname_logs(created_at);
	`
	_, err := db.Exec(schema)
	return err
}

func SeedInitialDrugs(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM drugs").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	seeds := []struct {
		unit      string
		name      string
		kfaCode   string
		unitPkg   string
		category  string
		buyPrice  int
		sellPrice int
		batchNum  string
		expDate   string
		stock     int
	}{
		{"KIDINEST", "Paracetamol Drop 100mg/ml", "KFA-93001", "Botol", "Antipiretik Anak", 22000, 35000, "BCH-PCM-01", "2027-12-31", 60},
		{"KIDINEST", "Amoxicillin Dry Sirup 125mg/5ml", "KFA-93002", "Botol", "Antibiotik Anak", 15000, 25000, "BCH-AMX-01", "2027-08-31", 45},
		{"KIDINEST", "Vaksin DTP-Hb-Hib Hexaxim", "KFA-93003", "Vial", "Vaksinasi Anak", 650000, 850000, "BCH-HEX-01", "2026-11-30", 25},
		{"KIDINEST", "Zinc Sirup 20mg/5ml", "KFA-93004", "Botol", "Suplemen Saluran Cerna", 18000, 28000, "BCH-ZNC-01", "2028-03-31", 50},
		{"KIDINEST", "Cetirizine Drop 10mg/ml", "KFA-93005", "Botol", "Antialergi Anak", 30000, 45000, "BCH-CTZ-01", "2027-09-30", 40},
		{"SSC", "Facial Cleanser Gentle Acne 100ml", "KFA-SSC-01", "Botol", "Cleanser Wajah", 45000, 75000, "BCH-CLN-01", "2027-10-31", 50},
		{"SSC", "Acne Day Cream Sunscreen SPF 30", "KFA-SSC-02", "Pot", "Sunscreen & Krim Pagi", 60000, 95000, "BCH-ADC-01", "2027-11-30", 45},
		{"SSC", "Serum Niacinamide 10% + Zinc 1%", "KFA-SSC-03", "Botol", "Serum Perawatan Kulit", 75000, 120000, "BCH-SER-01", "2027-07-31", 35},
		{"SSC", "Krim Malam Retinol 0.1% Micro", "KFA-SSC-04", "Pot", "Krim Malam Anti-Aging", 90000, 140000, "BCH-RET-01", "2027-06-30", 30},
		{"SSC", "Salep Gentamicin + Hidrokortison", "KFA-SSC-05", "Tube", "Salep Radang Kulit", 22000, 38000, "BCH-GNT-01", "2027-12-31", 40},
	}

	for _, s := range seeds {
		var drugID int64
		err := db.QueryRow(`
			INSERT INTO drugs (clinic_unit, name, kfa_code, unit, category, buy_price, sell_price, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE)
			RETURNING id
		`, s.unit, s.name, s.kfaCode, s.unitPkg, s.category, s.buyPrice, s.sellPrice).Scan(&drugID)
		if err != nil {
			return fmt.Errorf("gagal seed obat %s: %w", s.name, err)
		}

		_, err = db.Exec(`
			INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock)
			VALUES ($1, $2, $3, $4, 0)
		`, drugID, s.batchNum, s.expDate, s.stock)
		if err != nil {
			return fmt.Errorf("gagal seed batch obat %s: %w", s.name, err)
		}
	}

	log.Printf("[DB] Berhasil melakukan seed 10 master obat & batch awal untuk Kidinest & SSC.")
	return nil
}

func MigrateDoctorScheduleSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS doctor_schedules (
		id SERIAL PRIMARY KEY,
		doctor_id INTEGER NOT NULL REFERENCES staff_users(id) ON DELETE CASCADE,
		doctor_name VARCHAR(150) NOT NULL,
		clinic_unit VARCHAR(50) NOT NULL,
		day_of_week INT NOT NULL,
		shift_name VARCHAR(50) NOT NULL,
		start_time VARCHAR(10) NOT NULL,
		end_time VARCHAR(10) NOT NULL,
		slot_minutes INT NOT NULL DEFAULT 30,
		quota_per_slot INT NOT NULL DEFAULT 3,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_doc_schedules_doc ON doctor_schedules(doctor_id, day_of_week, is_active);

	CREATE TABLE IF NOT EXISTS doctor_schedule_overrides (
		id SERIAL PRIMARY KEY,
		doctor_id INTEGER NOT NULL REFERENCES staff_users(id) ON DELETE CASCADE,
		override_date DATE NOT NULL,
		is_cancelled BOOLEAN NOT NULL DEFAULT TRUE,
		reason VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_doc_overrides_date ON doctor_schedule_overrides(doctor_id, override_date);

	-- Perluasan tabel queues untuk booking tanggal masa depan & refund kasir
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS appointment_date DATE NOT NULL DEFAULT CURRENT_DATE;
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS appointment_time VARCHAR(30) NOT NULL DEFAULT '09:00 - 09:30';
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS patient_type VARCHAR(20) NOT NULL DEFAULT 'NEW';
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS reschedule_count INT NOT NULL DEFAULT 0;
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS refund_status VARCHAR(30) NOT NULL DEFAULT 'NONE';
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS refund_notes TEXT;
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS refund_bank VARCHAR(50);
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS refund_account_number VARCHAR(50);
	ALTER TABLE queues ADD COLUMN IF NOT EXISTS refund_account_name VARCHAR(150);

	CREATE INDEX IF NOT EXISTS idx_queues_appointment ON queues(appointment_date, doctor_id, queue_status);
	`
	_, err := db.Exec(schema)
	return err
}

func SeedInitialDoctorSchedules(db *sql.DB) error {
	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM doctor_schedules").Scan(&count)
	if count > 0 {
		return nil
	}

	// Ambil ID dokter dari staff_users
	docRows, err := db.Query("SELECT id, employee_number, name, clinic_unit FROM staff_users WHERE role = 'DOCTOR'")
	if err != nil {
		return err
	}
	defer docRows.Close()

	type docEntry struct {
		id     int64
		empNum string
		name   string
		unit   string
	}
	var docs []docEntry
	for docRows.Next() {
		var d docEntry
		if err := docRows.Scan(&d.id, &d.empNum, &d.name, &d.unit); err == nil {
			docs = append(docs, d)
		}
	}

	insertSched := `
	INSERT INTO doctor_schedules (doctor_id, doctor_name, clinic_unit, day_of_week, shift_name, start_time, end_time, slot_minutes, quota_per_slot, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, 30, $8, TRUE);
	`

	for _, d := range docs {
		switch d.empNum {
		case "DOC-KID-01": // dr. Patricia Vanessa A., Sp.A (Senin - Jumat Pagi 08:30 - 12:00)
			for day := 1; day <= 5; day++ {
				_, _ = db.Exec(insertSched, d.id, d.name, d.unit, day, "PAGI", "08:30", "12:00", 1)
			}
		case "DOC-KID-02": // dr. Hendra Pratama, Sp.A (Senin - Sabtu Sore 16:00 - 19:30)
			for day := 1; day <= 6; day++ {
				_, _ = db.Exec(insertSched, d.id, d.name, d.unit, day, "SORE", "16:00", "19:30", 1)
			}
		case "DOC-SSC-01": // dr. Dian Paramita, Sp.D.V.E (Senin - Jumat Pagi 09:00 - 13:00)
			for day := 1; day <= 5; day++ {
				_, _ = db.Exec(insertSched, d.id, d.name, d.unit, day, "PAGI", "09:00", "13:00", 1)
			}
		case "DOC-SSC-02": // dr. Kevin Wijaya, Sp.D.V.E (Senin - Sabtu Sore 15:00 - 19:00)
			for day := 1; day <= 6; day++ {
				_, _ = db.Exec(insertSched, d.id, d.name, d.unit, day, "SORE", "15:00", "19:00", 1)
			}
		}
	}

	// Update record yang sudah ter-seed agar quota_per_slot = 1
	_, _ = db.Exec("UPDATE doctor_schedules SET quota_per_slot = 1 WHERE quota_per_slot > 1")

	log.Printf("[DB] Berhasil melakukan seed jadwal rutin praktek 4 dokter spesialis.")
	return nil
}

func MigratePayrollSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS payroll_employees (
		id VARCHAR(50) PRIMARY KEY,
		nik VARCHAR(50) UNIQUE NOT NULL,
		name VARCHAR(150) NOT NULL,
		position VARCHAR(100) NOT NULL,
		department VARCHAR(100) NOT NULL,
		business_unit VARCHAR(20) NOT NULL CHECK (business_unit IN ('KIDINEST', 'SSC', 'ALL', 'MINIMO')),
		employment_type VARCHAR(20) NOT NULL CHECK (employment_type IN ('TETAP', 'KONTRAK', 'HONORER', 'MAGANG')),
		join_date DATE NOT NULL DEFAULT CURRENT_DATE,
		base_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
		position_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
		transport_meal_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
		bank_name VARCHAR(50) NOT NULL DEFAULT 'Bank Mandiri',
		bank_account_number VARCHAR(50) NOT NULL DEFAULT '',
		npwp VARCHAR(50) NOT NULL DEFAULT '',
		ptkp_status VARCHAR(10) NOT NULL DEFAULT 'TK/0',
		bpjs_kesehatan_no VARCHAR(50) NOT NULL DEFAULT '',
		bpjs_ketenagakerjaan_no VARCHAR(50) NOT NULL DEFAULT '',
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_payroll_emp_unit ON payroll_employees(business_unit);
	CREATE INDEX IF NOT EXISTS idx_payroll_emp_active ON payroll_employees(is_active);

	CREATE TABLE IF NOT EXISTS payroll_slips (
		id VARCHAR(50) PRIMARY KEY,
		slip_number VARCHAR(50) UNIQUE NOT NULL,
		employee_id VARCHAR(50) NOT NULL REFERENCES payroll_employees(id) ON DELETE RESTRICT,
		employee_nik VARCHAR(50) NOT NULL,
		employee_name VARCHAR(150) NOT NULL,
		position VARCHAR(100) NOT NULL,
		department VARCHAR(100) NOT NULL,
		business_unit VARCHAR(20) NOT NULL,
		period_month INT NOT NULL CHECK (period_month BETWEEN 1 AND 12),
		period_year INT NOT NULL CHECK (period_year >= 2020),
		period_name VARCHAR(50) NOT NULL,
		payment_date DATE NOT NULL DEFAULT CURRENT_DATE,
		working_days INT NOT NULL DEFAULT 22,
		present_days INT NOT NULL DEFAULT 22,
		sick_leave_days INT NOT NULL DEFAULT 0,
		paid_leave_days INT NOT NULL DEFAULT 0,
		absent_days INT NOT NULL DEFAULT 0,
		overtime_hours NUMERIC(5,2) NOT NULL DEFAULT 0,
		daily_meal_transport_rate NUMERIC(15,2) NOT NULL DEFAULT 30000,
		base_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
		position_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
		transport_meal_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
		shift_attendance_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
		incentive_bonus NUMERIC(15,2) NOT NULL DEFAULT 0,
		overtime_pay NUMERIC(15,2) NOT NULL DEFAULT 0,
		other_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
		total_gross_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
		bpjs_kesehatan NUMERIC(15,2) NOT NULL DEFAULT 0,
		bpjs_ketenagakerjaan NUMERIC(15,2) NOT NULL DEFAULT 0,
		pph21 NUMERIC(15,2) NOT NULL DEFAULT 0,
		attendance_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
		loan_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
		other_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
		total_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
		take_home_pay NUMERIC(15,2) NOT NULL DEFAULT 0,
		bank_name VARCHAR(50) NOT NULL DEFAULT 'Bank Mandiri',
		bank_account_number VARCHAR(50) NOT NULL DEFAULT '',
		status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'APPROVED', 'PAID')),
		approved_by VARCHAR(100),
		approved_at TIMESTAMP WITH TIME ZONE,
		paid_at TIMESTAMP WITH TIME ZONE,
		journal_id VARCHAR(50),
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_payroll_slips_period ON payroll_slips(period_year, period_month);
	CREATE INDEX IF NOT EXISTS idx_payroll_slips_emp ON payroll_slips(employee_id);
	CREATE INDEX IF NOT EXISTS idx_payroll_slips_unit ON payroll_slips(business_unit);
	CREATE INDEX IF NOT EXISTS idx_payroll_slips_status ON payroll_slips(status);

	ALTER TABLE payroll_slips ADD COLUMN IF NOT EXISTS working_days INT NOT NULL DEFAULT 22;
	ALTER TABLE payroll_slips ADD COLUMN IF NOT EXISTS present_days INT NOT NULL DEFAULT 22;
	ALTER TABLE payroll_slips ADD COLUMN IF NOT EXISTS sick_leave_days INT NOT NULL DEFAULT 0;
	ALTER TABLE payroll_slips ADD COLUMN IF NOT EXISTS paid_leave_days INT NOT NULL DEFAULT 0;
	ALTER TABLE payroll_slips ADD COLUMN IF NOT EXISTS absent_days INT NOT NULL DEFAULT 0;
	ALTER TABLE payroll_slips ADD COLUMN IF NOT EXISTS overtime_hours NUMERIC(5,2) NOT NULL DEFAULT 0;
	ALTER TABLE payroll_slips ADD COLUMN IF NOT EXISTS daily_meal_transport_rate NUMERIC(15,2) NOT NULL DEFAULT 30000;

	-- Relax & Expand business_unit check to allow dual company entities (Master Prompt 2 & 9 F)
	ALTER TABLE payroll_employees DROP CONSTRAINT IF EXISTS payroll_employees_business_unit_check;
	ALTER TABLE payroll_employees ADD CONSTRAINT payroll_employees_business_unit_check CHECK (business_unit IN ('KIDINEST', 'SSC', 'ALL', 'MINIMO', 'APOTEK', 'SMC', 'SMG'));
	`
	_, err := db.Exec(schema)
	return err
}

func SeedInitialPayrollEmployees(db *sql.DB) error {
	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM payroll_employees").Scan(&count)
	if count > 0 {
		return nil
	}

	seeds := []struct {
		id, nik, name, pos, dept, unit, empType, joinDate string
		base, posAllow, transAllow                         float64
		bank, acc, npwp, ptkp, bpjsKes, bpjsTk             string
	}{
		{
			"emp-001", "SMG-KID-2024-001", "Fitri Handayani", "Head Cashier & Kasir Shift GRO", "Front Office & Kasir", "KIDINEST", "TETAP", "2024-01-15",
			4800000, 750000, 650000, "Bank Mandiri", "136-00-9821345", "89.123.456.7-501.000", "TK/0", "0001239845671", "24098123456",
		},
		{
			"emp-002", "SMG-SSC-2024-002", "Nadia Safitri", "Kasir & Front Desk Resepsionis SSC", "Front Office & Kasir", "SSC", "TETAP", "2024-02-01",
			4500000, 500000, 650000, "Bank Mandiri", "136-00-5541290", "78.234.567.8-502.000", "TK/0", "0001892345612", "24098123899",
		},
		{
			"emp-003", "SMG-KID-2023-014", "Ns. Siti Rahmawati, S.Kep", "Perawat Senior Poli Spesialis Anak", "Keperawatan Medis", "KIDINEST", "TETAP", "2023-06-10",
			5200000, 900000, 700000, "Bank Mandiri", "136-00-7761823", "67.345.678.9-503.000", "K/1", "0002134567890", "23011234567",
		},
		{
			"emp-004", "SMG-SSC-2023-018", "Ns. Dwi Astuti, S.Kep", "Perawat Estetika & Asisten Laser Dermatologi", "Keperawatan Medis", "SSC", "TETAP", "2023-08-15",
			5100000, 850000, 700000, "Bank Mandiri", "136-00-4412987", "56.456.789.0-504.000", "TK/0", "0003245678901", "23018912345",
		},
		{
			"emp-005", "SMG-KID-2023-009", "apt. Rina Kusuma, S.Farm", "Kepala Instalasi Farmasi Kidinest", "Farmasi & Apotek", "KIDINEST", "TETAP", "2023-03-01",
			6200000, 1500000, 800000, "Bank BCA", "009-8812345", "45.567.890.1-505.000", "K/0", "0004356789012", "22091234567",
		},
		{
			"emp-006", "SMG-KID-2024-022", "Maya Indriyani, A.Md.Farm", "Asisten Apoteker & Staf Dispensing", "Farmasi & Apotek", "KIDINEST", "KONTRAK", "2024-04-01",
			4200000, 400000, 600000, "Bank Mandiri", "136-00-3321876", "34.678.901.2-506.000", "TK/0", "0005467890123", "24056789012",
		},
		{
			"emp-007", "SMG-SSC-2023-025", "Jessica Tan", "Senior Aesthetician & Beautician Facial", "Treatment & Estetika", "SSC", "TETAP", "2023-11-01",
			4400000, 600000, 650000, "Bank BCA", "009-5541299", "23.789.012.3-507.000", "TK/0", "0006578901234", "23112345678",
		},
		{
			"emp-008", "SMG-ALL-2022-005", "Bambang Sutrisno", "Teknisi Alkes & Pengelolaan Sarana Klinik", "Umum & Maintenance", "ALL", "TETAP", "2022-09-01",
			4300000, 500000, 650000, "Bank Mandiri", "136-00-1122339", "12.890.123.4-508.000", "K/2", "0007689012345", "22093456789",
		},
		{
			"emp-009", "SMG-ALL-2023-030", "Wahyu Setiawan, S.Kom", "Staf IT, Sistem POS & Integrasi SATUSEHAT", "Teknologi Informasi", "ALL", "TETAP", "2023-05-15",
			5500000, 900000, 700000, "Bank Mandiri", "136-00-9988112", "90.901.234.5-509.000", "TK/0", "0008790123456", "23056789012",
		},
		{
			"emp-010", "SMG-ALL-2022-002", "Ratna Paramita, S.E.", "Senior Accounting & Tax Specialist", "Keuangan & Akuntansi", "ALL", "TETAP", "2022-04-01",
			5800000, 1200000, 750000, "Bank Mandiri", "136-00-8877665", "81.012.345.6-510.000", "K/0", "0009801234567", "22045678901",
		},
		{
			"emp-011", "SMC-APT-2023-001", "apt. Nurul Hidayah, S.Farm", "Apoteker Penanggung Jawab Apotek (APA)", "Instalasi Farmasi & Apotek", "APOTEK", "TETAP", "2023-01-10",
			6500000, 1500000, 750000, "Bank BCA", "009-4412981", "44.123.789.0-511.000", "K/1", "0001928374651", "23019827364",
		},
		{
			"emp-012", "SMC-APT-2024-002", "Riko Saputra, A.Md.Farm", "Tenaga Teknis Kefarmasian & Logistik Medis", "Instalasi Farmasi & Apotek", "APOTEK", "TETAP", "2024-02-15",
			4500000, 500000, 650000, "Bank Mandiri", "136-00-2211990", "33.987.654.1-512.000", "TK/0", "0002819384756", "24029817263",
		},
	}

	insertEmp := `
	INSERT INTO payroll_employees (
		id, nik, name, position, department, business_unit, employment_type,
		join_date, base_salary, position_allowance, transport_meal_allowance,
		bank_name, bank_account_number, npwp, ptkp_status, bpjs_kesehatan_no,
		bpjs_ketenagakerjaan_no, is_active, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7,
		$8::date, $9, $10, $11,
		$12, $13, $14, $15, $16,
		$17, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	) ON CONFLICT (id) DO NOTHING;
	`

	for _, s := range seeds {
		_, err := db.Exec(insertEmp,
			s.id, s.nik, s.name, s.pos, s.dept, s.unit, s.empType,
			s.joinDate, s.base, s.posAllow, s.transAllow,
			s.bank, s.acc, s.npwp, s.ptkp, s.bpjsKes, s.bpjsTk,
		)
		if err != nil {
			log.Printf("[DB] Peringatan: gagal seed karyawan %s: %v", s.name, err)
		}
	}

	log.Printf("[DB] Berhasil melakukan seed 10 profil karyawan klinik Semarang Medical Group.")
	return nil
}
