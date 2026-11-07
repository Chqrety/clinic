package pharmacy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"clinic-api/internal/domain"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func RegisterRoutes(r chi.Router, s *Service) {
	r.Get("/drugs", s.HandleGetDrugs)
	r.Post("/drugs", s.HandleCreateDrug)
	r.Get("/drugs/{id}", s.HandleGetDrugDetail)
	r.Post("/drugs/{id}/restock", s.HandleRestock)
}

func (s *Service) HandleGetDrugs(w http.ResponseWriter, r *http.Request) {
	drugs, err := s.GetDrugs()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: drugs})
}

func (s *Service) HandleGetDrugDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid drug ID")
		return
	}

	drug, err := s.GetDrugDetail(id)
	if err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: drug})
}

type CreateDrugRequest struct {
	Name          string `json:"name"`
	KFACode       string `json:"kfa_code"`
	Unit          string `json:"unit"`
	Category      string `json:"category"`
	MinStockAlert int    `json:"min_stock_alert"`
	SellPrice     int    `json:"sell_price"`
}

func (s *Service) HandleCreateDrug(w http.ResponseWriter, r *http.Request) {
	var req CreateDrugRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Unit = strings.TrimSpace(req.Unit)
	req.Category = strings.TrimSpace(req.Category)
	if req.Name == "" || req.Unit == "" || req.Category == "" {
		sendError(w, http.StatusBadRequest, "Nama obat, satuan unit, dan kategori wajib diisi")
		return
	}
	if req.MinStockAlert <= 0 {
		req.MinStockAlert = 10
	}
	if req.SellPrice < 0 {
		req.SellPrice = 0
	}
	if req.KFACode == "" {
		req.KFACode = fmt.Sprintf("KFA-%d", time.Now().Unix()%1000000)
	}

	res, err := s.db.Exec(`
		INSERT INTO drugs (name, kfa_code, unit, category, min_stock_alert, sell_price)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.Name, req.KFACode, req.Unit, req.Category, req.MinStockAlert, req.SellPrice)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Gagal menyimpan jenis obat baru: "+err.Error())
		return
	}

	newID, _ := res.LastInsertId()
	created, err := s.GetDrugDetail(newID)
	if err != nil {
		sendJSON(w, http.StatusCreated, domain.ApiResponse{
			Success: true,
			Message: "Jenis obat baru berhasil ditambahkan",
			Data:    map[string]any{"id": newID},
		})
		return
	}

	sendJSON(w, http.StatusCreated, domain.ApiResponse{
		Success: true,
		Message: "Jenis obat baru berhasil ditambahkan",
		Data:    created,
	})
}

type RestockRequest struct {
	BatchNumber string `json:"batch_number"`
	ExpiredDate string `json:"expired_date"` // YYYY-MM-DD
	Quantity    int    `json:"quantity"`
	BuyPrice    int    `json:"buy_price"`
}

func (s *Service) HandleRestock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid drug ID")
		return
	}

	var req RestockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.BatchNumber == "" || req.ExpiredDate == "" || req.Quantity <= 0 {
		sendError(w, http.StatusBadRequest, "Batch, tanggal kadaluarsa, dan jumlah wajib valid")
		return
	}

	if err := s.Restock(id, req); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Message: "Restock obat berhasil dicatat ke gudang apotek",
	})
}

func (s *Service) GetDrugs() ([]domain.Drug, error) {
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
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		d.AvailableStock = d.PhysicalStock - d.ReservedStock
		d.IsLowStock = d.AvailableStock <= d.MinStockAlert

		// Check expiring batches (< 3 bulan)
		var expCount int
		_ = s.db.QueryRow(`
			SELECT COUNT(*) FROM drug_batches
			WHERE drug_id = ? AND physical_stock > 0 AND expired_date <= ?
		`, d.ID, threeMonthsLater).Scan(&expCount)
		d.HasExpiring = expCount > 0

		drugs = append(drugs, d)
	}
	return drugs, nil
}

func (s *Service) GetDrugDetail(drugID int64) (*domain.Drug, error) {
	var d domain.Drug
	err := s.db.QueryRow(`
		SELECT id, name, kfa_code, unit, category, min_stock_alert, sell_price, created_at
		FROM drugs WHERE id = ?
	`, drugID).Scan(&d.ID, &d.Name, &d.KFACode, &d.Unit, &d.Category, &d.MinStockAlert, &d.SellPrice, &d.CreatedAt)
	if err != nil {
		return nil, err
	}

	threeMonthsLater := time.Now().AddDate(0, 3, 0).Format("2006-01-02")
	batchRows, err := s.db.Query(`
		SELECT id, drug_id, batch_number, expired_date, physical_stock, reserved_stock, buy_price, created_at
		FROM drug_batches
		WHERE drug_id = ?
		ORDER BY expired_date ASC
	`, drugID)
	if err != nil {
		return nil, err
	}
	defer batchRows.Close()

	for batchRows.Next() {
		var b domain.DrugBatch
		if err := batchRows.Scan(
			&b.ID, &b.DrugID, &b.BatchNumber, &b.ExpiredDate, &b.PhysicalStock, &b.ReservedStock, &b.BuyPrice, &b.CreatedAt,
		); err != nil {
			return nil, err
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

func (s *Service) Restock(drugID int64, req RestockRequest) error {
	// If batch exists, update physical_stock, else insert
	var batchID int64
	err := s.db.QueryRow(`
		SELECT id FROM drug_batches WHERE drug_id = ? AND batch_number = ?
	`, drugID, req.BatchNumber).Scan(&batchID)

	if err == sql.ErrNoRows {
		_, err = s.db.Exec(`
			INSERT INTO drug_batches (drug_id, batch_number, expired_date, physical_stock, reserved_stock, buy_price)
			VALUES (?, ?, ?, ?, 0, ?)
		`, drugID, req.BatchNumber, req.ExpiredDate, req.Quantity, req.BuyPrice)
		return err
	} else if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		UPDATE drug_batches
		SET physical_stock = physical_stock + ?, expired_date = ?, buy_price = ?
		WHERE id = ?
	`, req.Quantity, req.ExpiredDate, req.BuyPrice, batchID)
	return err
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, domain.ApiResponse{Success: false, Error: msg})
}
