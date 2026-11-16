package inventory

import (
	"context"
	"time"

	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type MovementType string

const (
	MovementInPurchase  MovementType = "PURCHASE_RECEIPT"
	MovementOutSale     MovementType = "SALE_DISPENSE"
	MovementOutWaste    MovementType = "EXPIRED_WASTE"
	MovementAdjustment  MovementType = "STOCK_OPNAME_ADJ"
)

type StockBatch struct {
	ID             uuid.UUID       `json:"id"`
	BranchID       uuid.UUID       `json:"branch_id"`
	CostCenterID   uuid.UUID       `json:"cost_center_id"`
	ProductID      uuid.UUID       `json:"product_id"`
	ProductCode    string          `json:"product_code"`
	ProductName    string          `json:"product_name"`
	BatchNumber    string          `json:"batch_number"`
	ExpiredDate    time.Time       `json:"expired_date"`
	InitialQty     int             `json:"initial_qty"`
	CurrentQty     int             `json:"current_qty"`
	BuyPrice       decimal.Decimal `json:"buy_price"` // HPP FIFO
	SellPrice      decimal.Decimal `json:"sell_price"`
	InventoryAcc   string          `json:"inventory_acc"` // e.g. '1-1310', '1-1311', '1-1330'
	COGSAccount    string          `json:"cogs_account"`  // e.g. '5-1110', '5-1120', '5-2110'
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type StockMovement struct {
	ID             uuid.UUID       `json:"id"`
	BatchID        uuid.UUID       `json:"batch_id"`
	MovementType   MovementType    `json:"movement_type"`
	Quantity       int             `json:"quantity"`
	UnitPrice      decimal.Decimal `json:"unit_price"`
	TotalHPP       decimal.Decimal `json:"total_hpp"`
	ReferenceType  string          `json:"reference_type"` // 'INVOICE', 'PURCHASE_ORDER'
	ReferenceID    string          `json:"reference_id"`
	JournalID      *uuid.UUID      `json:"journal_id,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

type StockValuation struct {
	AccountCode   string          `json:"account_code"`
	AccountName   string          `json:"account_name"`
	TotalQty      int             `json:"total_qty"`
	TotalValue    decimal.Decimal `json:"total_value"`
}

type Repository interface {
	GetAvailableBatchesFIFO(ctx context.Context, branchID uuid.UUID, productID uuid.UUID, requiredQty int) ([]StockBatch, error)
	DeductBatchStock(ctx context.Context, batchID uuid.UUID, qty int) error
	RecordMovement(ctx context.Context, movement *StockMovement) error
	GetStockValuation(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) ([]StockValuation, error)
}
