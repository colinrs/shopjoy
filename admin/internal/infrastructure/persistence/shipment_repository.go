package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type shipmentRepo struct{}

func NewShipmentRepository() fulfillment.ShipmentRepository {
	return &shipmentRepo{}
}

// shipmentModel represents the database model for Shipment
type shipmentModel struct {
	ID               int64           `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID         int64           `gorm:"column:tenant_id;not null;index"`
	OrderID          int64           `gorm:"column:order_id;not null;index"`
	ShipmentNo       string          `gorm:"column:shipment_no;size:32;not null;uniqueIndex:uk_shipment_no"`
	Status           int             `gorm:"column:status;not null;default:0;index"`
	Carrier          string          `gorm:"column:carrier;size:50;not null;default:''"`
	CarrierCode      string          `gorm:"column:carrier_code;size:20;not null;default:''"`
	TrackingNo       string          `gorm:"column:tracking_no;size:100;not null;default:'';index"`
	ShippingCost     decimal.Decimal `gorm:"column:shipping_cost;type:decimal(19,4);not null;default:0"`
	ShippingCurrency string          `gorm:"column:shipping_currency;size:10;not null;default:'CNY'"`
	Weight           decimal.Decimal `gorm:"column:weight;type:decimal(10,3);not null;default:0"`
	ShippedAt        *time.Time      `gorm:"column:shipped_at"`
	DeliveredAt      *time.Time      `gorm:"column:delivered_at"`
	Remark           string          `gorm:"column:remark;size:500;not null;default:''"`
	DeletedAt        *time.Time      `gorm:"column:deleted_at;index"`
	CreatedAt        time.Time       `gorm:"column:created_at;not null"`
	UpdatedAt        time.Time       `gorm:"column:updated_at;not null"`
}

func (shipmentModel) TableName() string {
	return "shipments"
}

func (m *shipmentModel) toEntity() *fulfillment.Shipment {
	var deletedAt gorm.DeletedAt
	if m.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *m.DeletedAt, Valid: true}
	}
	return &fulfillment.Shipment{
		Model: application.Model{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: deletedAt,
		},
		TenantID:         shared.TenantID(m.TenantID),
		OrderID:          m.OrderID,
		ShipmentNo:       m.ShipmentNo,
		Status:           fulfillment.ShipmentStatus(m.Status),
		Carrier:          m.Carrier,
		CarrierCode:      m.CarrierCode,
		TrackingNo:       m.TrackingNo,
		ShippingCost:     m.ShippingCost,
		ShippingCurrency: m.ShippingCurrency,
		Weight:           m.Weight,
		ShippedAt:        m.ShippedAt,
		DeliveredAt:      m.DeliveredAt,
		Remark:           m.Remark,
	}
}

func fromShipmentEntity(s *fulfillment.Shipment) *shipmentModel {
	var deletedAt *time.Time
	if s.Model.DeletedAt.Valid {
		deletedAt = &s.Model.DeletedAt.Time
	}
	return &shipmentModel{
		ID:               s.Model.ID,
		TenantID:         s.TenantID.Int64(),
		OrderID:          s.OrderID,
		ShipmentNo:       s.ShipmentNo,
		Status:           int(s.Status),
		Carrier:          s.Carrier,
		CarrierCode:      s.CarrierCode,
		TrackingNo:       s.TrackingNo,
		ShippingCost:     s.ShippingCost,
		ShippingCurrency: s.ShippingCurrency,
		Weight:           s.Weight,
		ShippedAt:        s.ShippedAt,
		DeliveredAt:      s.DeliveredAt,
		Remark:           s.Remark,
		DeletedAt:        deletedAt,
		CreatedAt:        s.Model.CreatedAt,
		UpdatedAt:        s.Model.UpdatedAt,
	}
}

// Create inserts a new shipment
func (r *shipmentRepo) Create(ctx context.Context, db *gorm.DB, s *fulfillment.Shipment) error {
	model := fromShipmentEntity(s)
	return db.WithContext(ctx).Create(model).Error
}

// Update updates an existing shipment
func (r *shipmentRepo) Update(ctx context.Context, db *gorm.DB, s *fulfillment.Shipment) error {
	model := fromShipmentEntity(s)
	return db.WithContext(ctx).
		Model(&shipmentModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", s.ID, s.TenantID.Int64()).
		Updates(map[string]interface{}{
			"status":            model.Status,
			"carrier":           model.Carrier,
			"carrier_code":      model.CarrierCode,
			"tracking_no":       model.TrackingNo,
			"shipping_cost":     model.ShippingCost,
			"shipping_currency": model.ShippingCurrency,
			"weight":            model.Weight,
			"shipped_at":        model.ShippedAt,
			"delivered_at":      model.DeliveredAt,
			"remark":            model.Remark,
			"updated_at":        model.UpdatedAt,
		}).Error
}

// FindByID finds a shipment by ID
func (r *shipmentRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*fulfillment.Shipment, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model shipmentModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrShipmentNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByShipmentNo finds a shipment by shipment number
func (r *shipmentRepo) FindByShipmentNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, shipmentNo string) (*fulfillment.Shipment, error) {
	query := db.WithContext(ctx).Model(&shipmentModel{}).Where("shipment_no = ? AND deleted_at IS NULL", shipmentNo)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model shipmentModel
	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrShipmentNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByOrderID finds all shipments for an order
func (r *shipmentRepo) FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) ([]*fulfillment.Shipment, error) {
	query := db.WithContext(ctx).Model(&shipmentModel{}).Where("order_id = ? AND deleted_at IS NULL", orderID)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var models []shipmentModel
	err := query.Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	shipments := make([]*fulfillment.Shipment, len(models))
	for i, m := range models {
		shipments[i] = m.toEntity()
	}
	return shipments, nil
}

// FindByTrackingNo finds a shipment by tracking number
func (r *shipmentRepo) FindByTrackingNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, trackingNo string) (*fulfillment.Shipment, error) {
	query := db.WithContext(ctx).Model(&shipmentModel{}).Where("tracking_no = ? AND deleted_at IS NULL", trackingNo)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model shipmentModel
	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrShipmentNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindList finds shipments with pagination and filters
func (r *shipmentRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query fulfillment.ShipmentQuery) ([]*fulfillment.Shipment, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&shipmentModel{}).Where("deleted_at IS NULL")

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.OrderID != 0 {
		dbQuery = dbQuery.Where("order_id = ?", query.OrderID)
	}
	if query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.CarrierCode != "" {
		dbQuery = dbQuery.Where("carrier_code = ?", query.CarrierCode)
	}
	if query.TrackingNo != "" {
		dbQuery = dbQuery.Where("tracking_no LIKE ?", escapeLikePattern(query.TrackingNo))
	}
	if !query.StartTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime)
	}
	if !query.EndTime.IsZero() {
		dbQuery = dbQuery.Where("created_at < ?", query.EndTime)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []shipmentModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	shipments := make([]*fulfillment.Shipment, len(models))
	for i, m := range models {
		shipments[i] = m.toEntity()
	}
	return shipments, total, nil
}

// Delete soft deletes a shipment
func (r *shipmentRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&shipmentModel{}).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	result := query.Update("deleted_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrShipmentNotFound
	}
	return nil
}

// CountByStatus counts shipments by status
func (r *shipmentRepo) CountByStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, status fulfillment.ShipmentStatus) (int64, error) {
	query := db.WithContext(ctx).Model(&shipmentModel{}).Where("status = ? AND deleted_at IS NULL", status)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// escapeLikePattern escapes special LIKE pattern characters and wraps with wildcards
func escapeLikePattern(pattern string) string {
	// Escape special LIKE characters: %, _, \
	escaped := strings.ReplaceAll(pattern, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "%", "\\%")
	escaped = strings.ReplaceAll(escaped, "_", "\\_")
	return fmt.Sprintf("%%%s%%", escaped)
}
