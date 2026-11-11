package pharmacy

import (
	"database/sql"
	"fmt"
	"time"

	"clinic-api/internal/domain"
)

// Repository defines data access operations for pharmacy drugs and batch inventory.
type Repository interface {
	GetDrugs() ([]domain.Drug, error)
	GetDrugDetail(drugID int64) (*domain.Drug, error)
	CreateDrug(name, kfaCode, unit, category string, minStockAlert, sellPrice int) (int64, error)
	FindBatchByNumber(drugID int64, batchNumber string) (int64, error)
	InsertBatch(drugID int64, batchNumber, expiredDate string, quantity, buyPrice int) (int64, error)
	UpdateBatchStock(batchID int64, quantity int, expiredDate string, buyPrice int) error
}

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) GetDrugs() ([]domain.Drug, error) {
	query := `
		SELECT d.id, d.name, d.kfa_code, d.unit, d.category, d.min_stock_alert, d.sell_price,
		       coalesce(sum(b.physical_stock), 0) as physical_stock,
		       coalesce(sum(b.reserved_stock), 0) as reserved_stock,
		       d.created_at
		FROM drugs d
		LEFT JOIN drug_batches b ON d.id = b.drug_id
		GROUP BY d.id
		ORDER BY d.name ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query drugs: %w", err)
	}
	defer rows.Close()

	var drugs []domain.Drug
	threeMonthsLater := time.Now().AddDate(0, 3, 0).Format("2006-01-02")

	for rows.Next() {
		var d domain.Drug
		if err := rows.Scan(
			&d.ID, &d.Name, &d.KFACode, &d.Unit, &d.Category, &d.MinStockAlert, &d.SellPrice,
			&d.PhysicalStock, &d.ReservedStock, &d.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan drug row: %w", err)
		}
		d.AvailableStock = d.PhysicalStock - d.ReservedStock
		d.IsLowStock = d.AvailableStock <= d.MinStockAlert

		// Check expiring batches (< 3 months)
		var expCount int
		_ = r.db.QueryRow(`
			SELECT COUNT(*) FROM drug_batches
			WHERE drug_id = ? AND physical_stock > 0 AND expired_date <= ?
		`, d.ID, threeMonthsLater).Scan(&expCount)
		d.HasExpiring = expCount > 0

		drugs = append(drugs, d)
	}
	return drugs, nil
}

func (r *sqlRepository) GetDrugDetail(drugID int64) (*domain.Drug, error) {
	var d domain.Drug
	err := r.db.QueryRow(`
		SELECT id, name, kfa_code, unit, category, min_stock_alert, sell_price, created_at
		FROM drugs WHERE id = ?
	`, drugID).Scan(&d.ID, &d.Name, &d.KFACode, &d.Unit, &d.Category, &d.MinStockAlert, &d.SellPrice, &d.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get drug header %d: %w", drugID, err)
	}

	threeMonthsLater := time.Now().AddDate(0, 3, 0).Format("2006-01-02")
	batchRows, err := r.db.Query(`
		SELECT id, drug_id, batch_number, expired_date, physical_stock, reserved_stock, buy_price, created_at
		FROM drug_batches
		WHERE drug_id = ?
		ORDER BY expired_date ASC
	`, drugID)
	if err != nil {
		return nil, fmt.Errorf("query batches for drug %d: %w", drugID, err)
	}
	defer batchRows.Close()

	for batchRows.Next() {
		var b domain.DrugBatch
		if err := batchRows.Scan(
			&b.ID, &b.DrugID, &b.BatchNumber, &b.ExpiredDate, &b.PhysicalStock, &b.ReservedStock, &b.BuyPrice, &b.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan drug batch: %w", err)
		}
		b.AvailableStock = b.PhysicalStock - b.ReservedStock
		b.IsExpiringSoon = b.ExpiredDate <= threeMonthsLater && b.PhysicalStock > 0

		d.PhysicalStock += b.PhysicalStock
		d.ReservedStock += b.ReservedStock
		if b.IsExpiringSoon {
			d.HasExpiring = true
		}
		d.Batches = append(d.Batches, b)
	}

	d.AvailableStock = d.PhysicalStock - d.ReservedStock
	d.IsLowStock = d.AvailableStock <= d.MinStockAlert
	return &d, nil
}

func (r *sqlRepository) CreateDrug(name, kfaCode, unit, category string, minStockAlert, sellPrice int) (int64, error) {
	query := `
		INSERT INTO drugs (name, kfa_code, unit, category, min_stock_alert, sell_price)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.Exec(query, name, kfaCode, unit, category, minStockAlert, sellPrice)
	if err != nil {
		return 0, fmt.Errorf("insert drug: %w", err)
	}
	return res.LastInsertId()
}

func (r *sqlRepository) FindBatchByNumber(drugID int64, batchNumber string) (int64, error) {
	query := `SELECT id FROM drug_batches WHERE drug_id = ? AND batch_number = ?`
	var batchID int64
	err := r.db.QueryRow(query, drugID, batchNumber).Scan(&batchID)
	return batchID, err
}

func (r *sqlRepository) InsertBatch(drugID int64, batchNumber, expiredDate string, quantity, buyPrice int) (int64, error) {
	query := `
		INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, buy_price)
		VALUES (?, ?, ?, ?, 0, ?)
	`
	res, err := r.db.Exec(query, drugID, batchNumber, expiredDate, quantity, buyPrice)
	if err != nil {
		return 0, fmt.Errorf("insert drug batch: %w", err)
	}
	return res.LastInsertId()
}

func (r *sqlRepository) UpdateBatchStock(batchID int64, quantity int, expiredDate string, buyPrice int) error {
	query := `
		UPDATE drug_batches
		SET physical_stock = physical_stock + ?, expired_date = ?, buy_price = ?
		WHERE id = ?
	`
	res, err := r.db.Exec(query, quantity, expiredDate, buyPrice, batchID)
	if err != nil {
		return fmt.Errorf("update batch stock: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("batch not found or not modified")
	}
	return nil
}

