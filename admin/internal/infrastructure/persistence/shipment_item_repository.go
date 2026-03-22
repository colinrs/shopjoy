package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type shipmentItemRepo struct{}

func NewShipmentItemRepository() fulfillment.ShipmentItemRepository {
	return &shipmentItemRepo{}
}

// shipmentItemModel represents the database model for ShipmentItem
type shipmentItemModel struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID    int64  `gorm:"column:tenant_id;not null;index"`
	ShipmentID  int64  `gorm:"column:shipment_id;not null;index"`
	OrderItemID int64  `gorm:"column:order_item_id;not null;index"`
	ProductID   int64  `gorm:"column:product_id;not null;index"`
	SKUID       int64  `gorm:"column:sku_id;not null;index"`
	ProductName string `gorm:"column:product_name;size:255;not null;default:''"`
	SKUName     string `gorm:"column:sku_name;size:255;not null;default:''"`
	Image       string `gorm:"column:image;size:500;not null;default:''"`
	Quantity    int    `gorm:"column:quantity;not null;default:1"`
	CreatedAt   int64  `gorm:"column:created_at;not null"`
}

func (shipmentItemModel) TableName() string {
	return "shipment_items"
}

func (m *shipmentItemModel) toEntity() fulfillment.ShipmentItem {
	return fulfillment.ShipmentItem{
		ID:          m.ID,
		TenantID:    shared.TenantID(m.TenantID),
		ShipmentID:  m.ShipmentID,
		OrderItemID: m.OrderItemID,
		ProductID:   m.ProductID,
		SKUID:       m.SKUID,
		ProductName: m.ProductName,
		SKUName:     m.SKUName,
		Image:       m.Image,
		Quantity:    m.Quantity,
		CreatedAt:   time.Unix(m.CreatedAt, 0).UTC(),
	}
}

func fromShipmentItemEntity(item fulfillment.ShipmentItem) *shipmentItemModel {
	return &shipmentItemModel{
		ID:          item.ID,
		TenantID:    item.TenantID.Int64(),
		ShipmentID:  item.ShipmentID,
		OrderItemID: item.OrderItemID,
		ProductID:   item.ProductID,
		SKUID:       item.SKUID,
		ProductName: item.ProductName,
		SKUName:     item.SKUName,
		Image:       item.Image,
		Quantity:    item.Quantity,
		CreatedAt:   item.CreatedAt.Unix(),
	}
}

// BatchCreate inserts multiple shipment items
func (r *shipmentItemRepo) BatchCreate(ctx context.Context, db *gorm.DB, items []fulfillment.ShipmentItem) error {
	if len(items) == 0 {
		return nil
	}
	models := make([]*shipmentItemModel, len(items))
	for i, item := range items {
		models[i] = fromShipmentItemEntity(item)
	}
	return db.WithContext(ctx).Create(&models).Error
}

// FindByShipmentID finds all items for a shipment with tenant isolation
func (r *shipmentItemRepo) FindByShipmentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, shipmentID int64) ([]fulfillment.ShipmentItem, error) {
	query := db.WithContext(ctx).Where("shipment_id = ?", shipmentID)
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var models []shipmentItemModel
	err := query.Order("id ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := make([]fulfillment.ShipmentItem, len(models))
	for i, m := range models {
		items[i] = m.toEntity()
	}
	return items, nil
}

// FindByOrderItemID finds all shipment items for an order item with tenant isolation
func (r *shipmentItemRepo) FindByOrderItemID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderItemID int64) ([]fulfillment.ShipmentItem, error) {
	query := db.WithContext(ctx).Where("order_item_id = ?", orderItemID)
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var models []shipmentItemModel
	err := query.Order("id ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := make([]fulfillment.ShipmentItem, len(models))
	for i, m := range models {
		items[i] = m.toEntity()
	}
	return items, nil
}

// DeleteByShipmentID deletes all items for a shipment with tenant isolation
func (r *shipmentItemRepo) DeleteByShipmentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, shipmentID int64) error {
	query := db.WithContext(ctx).Where("shipment_id = ?", shipmentID)
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	return query.Delete(&shipmentItemModel{}).Error
}