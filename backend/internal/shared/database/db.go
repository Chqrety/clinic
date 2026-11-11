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

	if err := SeedInitialStaff(db); err != nil {
		log.Printf("[DB] Peringatan: Seeding data staf gagal: %v", err)
	}

	if err := SeedInitialDrugs(db); err != nil {
		log.Printf("[DB] Peringatan: Seeding data obat gagal: %v", err)
	}

	return db, nil
}

func MigrateStaffSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS staff_users (
		id SERIAL PRIMARY KEY,
		employee_number VARCHAR(50) UNIQUE NOT NULL,
		name VARCHAR(150) NOT NULL,
		role VARCHAR(30) NOT NULL CHECK (role IN ('SUPERADMIN_OWNER', 'DOCTOR', 'GRO_CASHIER', 'PHARMACIST')),
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
			'QRIS', 'BANK_TRANSFER', 'CASH', 'UNPAID_BOOKING'
		)),
		payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('PENDING', 'PAID', 'REFUNDED')),
		queue_status VARCHAR(30) NOT NULL CHECK (queue_status IN (
			'WAITING', 'CALLING', 'IN_EXAM', 'PHARMACY_WAITING', 'COMPLETED', 'CANCELLED'
		)),
		paid_at TIMESTAMP WITH TIME ZONE,
		called_at TIMESTAMP WITH TIME ZONE,
		completed_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_queues_unit_status ON queues(clinic_unit, queue_status);
	CREATE INDEX IF NOT EXISTS idx_queues_created_at ON queues(created_at);
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
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM staff_users").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

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
		{"DOC-KID-01", "dr. Hendra Pratama, Sp.A", "DOCTOR", "KIDINEST"},
		{"DOC-SSC-01", "dr. Nadia Kusuma, Sp.KK", "DOCTOR", "SSC"},
		{"GRO-01", "Siti Sarah, A.Md", "GRO_CASHIER", "ALL"},
		{"FARMASI-01", "Apt. Dimas Wicaksono, S.Farm", "PHARMACIST", "ALL"},
	}

	query := `
	INSERT INTO staff_users (employee_number, name, role, clinic_unit, password_hash, is_active)
	VALUES ($1, $2, $3, $4, $5, TRUE)
	ON CONFLICT (employee_number) DO NOTHING;
	`

	for _, s := range seeds {
		if _, err := db.Exec(query, s.empNum, s.name, s.role, s.unit, string(defaultHash)); err != nil {
			return fmt.Errorf("gagal insert seed staf %s: %w", s.empNum, err)
		}
	}

	log.Printf("[DB] Berhasil melakukan seed 5 akun staf demo awal.")
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
		item_name VARCHAR(150) NOT NULL,
		is_manual BOOLEAN NOT NULL DEFAULT FALSE,
		quantity INTEGER NOT NULL DEFAULT 1,
		unit_price INTEGER NOT NULL DEFAULT 0,
		subtotal INTEGER NOT NULL DEFAULT 0,
		dosage_instruction VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_prescription_items_rx_id ON prescription_items(prescription_id);

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
