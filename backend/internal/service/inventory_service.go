package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"go.uber.org/zap"
)

type InventoryService interface {
	ListItems(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.InventoryItemResponse, int64, error)
	GetItem(ctx context.Context, id string) (*dto.InventoryItemResponse, error)
	CreateItem(ctx context.Context, schoolID string, req dto.CreateInventoryItemRequest) (*dto.InventoryItemResponse, error)
	RecordMovement(ctx context.Context, itemID, createdBy string, req dto.StockMovementRequest) error

	ListAssets(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.AssetResponse, int64, error)
	GetAsset(ctx context.Context, id string) (*dto.AssetResponse, error)
	CreateAsset(ctx context.Context, schoolID string, req dto.CreateAssetRequest) (*dto.AssetResponse, error)
}

type inventoryService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewInventoryService(db *sqlx.DB, logger *zap.Logger) InventoryService {
	return &inventoryService{db: db, logger: logger}
}

func (s *inventoryService) ListItems(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.InventoryItemResponse, int64, error) {
	filter.Defaults()

	var items []domain.InventoryItem
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM inventory_items WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM inventory_items WHERE school_id=$1 ORDER BY name LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list items", err)
	}

	result := make([]dto.InventoryItemResponse, len(items))
	for i, item := range items {
		result[i] = dto.InventoryItemResponse{
			ID:           item.ID,
			SchoolID:     item.SchoolID,
			Code:         item.Code,
			Name:         item.Name,
			Category:     item.Category,
			Unit:         item.Unit,
			StockIn:      item.StockIn,
			StockOut:     item.StockOut,
			StockCurrent: item.StockIn - item.StockOut,
			StockMin:     item.StockMin,
			Location:     item.Location,
			CreatedAt:    item.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *inventoryService) GetItem(ctx context.Context, id string) (*dto.InventoryItemResponse, error) {
	var item domain.InventoryItem
	if err := s.db.GetContext(ctx, &item, `SELECT * FROM inventory_items WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("inventory item", id)
	}
	return &dto.InventoryItemResponse{
		ID:           item.ID,
		SchoolID:     item.SchoolID,
		Code:         item.Code,
		Name:         item.Name,
		Category:     item.Category,
		Unit:         item.Unit,
		StockIn:      item.StockIn,
		StockOut:     item.StockOut,
		StockCurrent: item.StockIn - item.StockOut,
		StockMin:     item.StockMin,
		Location:     item.Location,
		CreatedAt:    item.CreatedAt,
	}, nil
}

func (s *inventoryService) CreateItem(ctx context.Context, schoolID string, req dto.CreateInventoryItemRequest) (*dto.InventoryItemResponse, error) {
	item := &domain.InventoryItem{
		ID:       uuid.New().String(),
		SchoolID: schoolID,
		Code:     req.Code,
		Name:     req.Name,
		Category: req.Category,
		Unit:     req.Unit,
		StockMin: req.StockMin,
		Location: req.Location,
	}

	query := `INSERT INTO inventory_items (id, school_id, code, name, category, unit, stock_min, location) VALUES (:id, :school_id, :code, :name, :category, :unit, :stock_min, :location)`
	if _, err := s.db.NamedExecContext(ctx, query, item); err != nil {
		return nil, domain.NewInternalError("failed to create item", err)
	}

	return s.GetItem(ctx, item.ID)
}

func (s *inventoryService) RecordMovement(ctx context.Context, itemID, createdBy string, req dto.StockMovementRequest) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	movement := &domain.StockMovement{
		ID:            uuid.New().String(),
		ItemID:        itemID,
		Type:          req.Type,
		Quantity:      req.Quantity,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Notes:         req.Notes,
		CreatedBy:     createdBy,
	}

	movQuery := `INSERT INTO stock_movements (id, item_id, type, quantity, reference_type, reference_id, notes, created_by) VALUES (:id, :item_id, :type, :quantity, :reference_type, :reference_id, :notes, :created_by)`
	if _, err := tx.NamedExecContext(ctx, movQuery, movement); err != nil {
		return domain.NewInternalError("failed to record movement", err)
	}

	if req.Type == "in" {
		tx.ExecContext(ctx, `UPDATE inventory_items SET stock_in=stock_in+$1, updated_at=NOW() WHERE id=$2`, req.Quantity, itemID)
	} else {
		tx.ExecContext(ctx, `UPDATE inventory_items SET stock_out=stock_out+$1, updated_at=NOW() WHERE id=$2`, req.Quantity, itemID)
	}

	return tx.Commit()
}

func (s *inventoryService) ListAssets(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.AssetResponse, int64, error) {
	filter.Defaults()

	var items []domain.Asset
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM assets WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM assets WHERE school_id=$1 ORDER BY name LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list assets", err)
	}

	result := make([]dto.AssetResponse, len(items))
	for i, a := range items {
		result[i] = dto.AssetResponse{
			ID:               a.ID,
			SchoolID:         a.SchoolID,
			Code:             a.Code,
			Name:             a.Name,
			Category:         a.Category,
			PurchaseDate:     a.PurchaseDate,
			PurchasePrice:    a.PurchasePrice,
			CurrentValue:     a.CurrentValue,
			Location:         a.Location,
			Condition:        a.Condition,
			Status:           a.Status,
			DepreciationRate: a.DepreciationRate,
			CreatedAt:        a.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *inventoryService) GetAsset(ctx context.Context, id string) (*dto.AssetResponse, error) {
	var a domain.Asset
	if err := s.db.GetContext(ctx, &a, `SELECT * FROM assets WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("asset", id)
	}
	return &dto.AssetResponse{
		ID:               a.ID,
		SchoolID:         a.SchoolID,
		Code:             a.Code,
		Name:             a.Name,
		Category:         a.Category,
		PurchaseDate:     a.PurchaseDate,
		PurchasePrice:    a.PurchasePrice,
		CurrentValue:     a.CurrentValue,
		Location:         a.Location,
		Condition:        a.Condition,
		Status:           a.Status,
		DepreciationRate: a.DepreciationRate,
		CreatedAt:        a.CreatedAt,
	}, nil
}

func (s *inventoryService) CreateAsset(ctx context.Context, schoolID string, req dto.CreateAssetRequest) (*dto.AssetResponse, error) {
	a := &domain.Asset{
		ID:               uuid.New().String(),
		SchoolID:         schoolID,
		Code:             req.Code,
		Name:             req.Name,
		Category:         req.Category,
		PurchaseDate:     req.PurchaseDate,
		PurchasePrice:    req.PurchasePrice,
		CurrentValue:     req.PurchasePrice,
		Location:         req.Location,
		DepreciationRate: req.DepreciationRate,
		ResponsibleID:    req.ResponsibleID,
		Status:           "active",
		Condition:        "good",
	}

	query := `INSERT INTO assets (id, school_id, code, name, category, purchase_date, purchase_price, current_value, location, depreciation_rate, responsible_id, status, condition) VALUES (:id, :school_id, :code, :name, :category, :purchase_date, :purchase_price, :current_value, :location, :depreciation_rate, :responsible_id, :status, :condition)`
	if _, err := s.db.NamedExecContext(ctx, query, a); err != nil {
		return nil, domain.NewInternalError("failed to create asset", err)
	}

	return s.GetAsset(ctx, a.ID)
}
