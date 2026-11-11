package pharmacy

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"clinic-api/internal/domain"
)

type CreateDrugRequest struct {
	Name          string `json:"name"`
	KFACode       string `json:"kfa_code"`
	Unit          string `json:"unit"`
	Category      string `json:"category"`
	MinStockAlert int    `json:"min_stock_alert"`
	SellPrice     int    `json:"sell_price"`
}

type RestockRequest struct {
	BatchNumber string `json:"batch_number"`
	ExpiredDate string `json:"expired_date"` // YYYY-MM-DD
	Quantity    int    `json:"quantity"`
	BuyPrice    int    `json:"buy_price"`
}

type Service struct {
	repo Repository
}

// NewService creates a pharmacy Service using standard database/sql.DB.
func NewService(db *sql.DB) *Service {
	return &Service{repo: NewSQLRepository(db)}
}

// NewServiceWithRepo creates a pharmacy Service using a custom Repository.
func NewServiceWithRepo(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDrugs() ([]domain.Drug, error) {
	return s.repo.GetDrugs()
}

func (s *Service) GetDrugsPaginated(filter DrugFilter) ([]domain.Drug, int64, error) {
	return s.repo.GetDrugsPaginated(filter)
}

func (s *Service) GetDrugDetail(drugID int64) (*domain.Drug, error) {
	return s.repo.GetDrugDetail(drugID)
}

func (s *Service) CreateDrug(req CreateDrugRequest) (*domain.Drug, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Unit = strings.TrimSpace(req.Unit)
	req.Category = strings.TrimSpace(req.Category)
	if req.Name == "" || req.Unit == "" || req.Category == "" {
		return nil, errors.New("nama obat, satuan unit, dan kategori wajib diisi")
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

	newID, err := s.repo.CreateDrug(req.Name, req.KFACode, req.Unit, req.Category, req.MinStockAlert, req.SellPrice)
	if err != nil {
		return nil, fmt.Errorf("create drug: %w", err)
	}

	return s.GetDrugDetail(newID)
}

func (s *Service) Restock(drugID int64, req RestockRequest) error {
	if req.BatchNumber == "" || req.ExpiredDate == "" || req.Quantity <= 0 {
		return errors.New("batch, tanggal kadaluarsa, dan jumlah wajib valid")
	}

	batchID, err := s.repo.FindBatchByNumber(drugID, req.BatchNumber)
	if err == sql.ErrNoRows {
		_, err = s.repo.InsertBatch(drugID, req.BatchNumber, req.ExpiredDate, req.Quantity, req.BuyPrice)
		return err
	} else if err != nil {
		return fmt.Errorf("find batch: %w", err)
	}

	return s.repo.UpdateBatchStock(batchID, req.Quantity, req.ExpiredDate, req.BuyPrice)
}

