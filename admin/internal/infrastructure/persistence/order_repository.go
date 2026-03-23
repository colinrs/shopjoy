package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type orderRepo struct{}

func NewOrderRepository() fulfillment.OrderRepository {
	return &orderRepo{}
}

type orderModel struct {
	ID               int64     `gorm:"column:id;primaryKey"`
	TenantID         int64     `gorm:"column:tenant_id;not null;index"`
	OrderNo          string    `gorm:"column:order_no;not null;uniqueIndex:uk_order_no"`
	UserID           int64     `gorm:"column:user_id;not null;index"`
	Status           string    `gorm:"column:status;not null;default:'pending_payment';index"`
	FulfillmentStatus int8     `gorm:"column:fulfillment_status;not null;default:0;index"`
	RefundStatus     int8      `gorm:"column:refund_status;not null;default:0;index"`
	TotalAmount      int64     `gorm:"column:total_amount;not null"`
	DiscountAmount   int64     `gorm:"column:discount_amount;not null"`
	ShippingFee      int64     `gorm:"column:shipping_fee;not null"`
	PayAmount        int64     `gorm:"column:pay_amount;not null"`
	Currency         string    `gorm:"column:currency;not null;default:'CNY'"`
	MerchantRemark   string    `gorm:"column:merchant_remark;not null;default:''"`
	Remark           string    `gorm:"column:remark;not null;default:''"`
	OriginalAmount   int64     `gorm:"column:original_amount;not null;default:0"`
	AdjustAmount     int64     `gorm:"column:adjust_amount;not null;default:0"`
	AdjustReason     string    `gorm:"column:adjust_reason;not null;default:''"`
	AdjustedBy       int64     `gorm:"column:adjusted_by;not null;default:0"`
	AdjustedAt       *int64    `gorm:"column:adjusted_at"`
	Version          int       `gorm:"column:version;not null;default:1"`
	PaymentMethod    string    `gorm:"column:payment_method;not null;default:''"`
	Source           string    `gorm:"column:source;not null;default:''"`
	ReceiverName     string    `gorm:"column:receiver_name;not null"`
	ReceiverPhone    string    `gorm:"column:receiver_phone;not null"`
	ReceiverAddress  string    `gorm:"column:receiver_address;not null"`
	PaidAt           *int64    `gorm:"column:paid_at"`
	ShippedAt        *int64    `gorm:"column:shipped_at"`
	DeliveredAt      *int64    `gorm:"column:delivered_at"`
	CancelledAt      *int64    `gorm:"column:cancelled_at"`
	CreatedBy        int64     `gorm:"column:created_by;not null"`
	UpdatedBy        int64     `gorm:"column:updated_by;not null"`
	DeletedAt        *int64    `gorm:"column:deleted_at;index"`
	CreatedAt        int64     `gorm:"column:created_at;not null"`
	UpdatedAt        int64     `gorm:"column:updated_at;not null"`
}

func (orderModel) TableName() string {
	return "orders"
}

func (m *orderModel) toEntity() *fulfillment.Order {
	var paidAt, shippedAt, deliveredAt, cancelledAt, adjustedAt *time.Time
	if m.PaidAt != nil {
		t := time.Unix(*m.PaidAt, 0).UTC()
		paidAt = &t
	}
	if m.ShippedAt != nil {
		t := time.Unix(*m.ShippedAt, 0).UTC()
		shippedAt = &t
	}
	if m.DeliveredAt != nil {
		t := time.Unix(*m.DeliveredAt, 0).UTC()
		deliveredAt = &t
	}
	if m.CancelledAt != nil {
		t := time.Unix(*m.CancelledAt, 0).UTC()
		cancelledAt = &t
	}
	if m.AdjustedAt != nil {
		t := time.Unix(*m.AdjustedAt, 0).UTC()
		adjustedAt = &t
	}

	var deletedAt gorm.DeletedAt
	if m.DeletedAt != nil {
		t := time.Unix(*m.DeletedAt, 0).UTC()
		deletedAt = gorm.DeletedAt{Time: t, Valid: true}
	}

	return &fulfillment.Order{
		ID:                m.ID,
		TenantID:          shared.TenantID(m.TenantID),
		OrderNo:           m.OrderNo,
		UserID:            m.UserID,
		Status:            fulfillment.OrderStatus(m.Status),
		FulfillmentStatus: fulfillment.OrderFulfillmentStatus(m.FulfillmentStatus),
		RefundStatus:      fulfillment.OrderRefundStatus(m.RefundStatus),
		TotalAmount:       m.TotalAmount,
		DiscountAmount:    m.DiscountAmount,
		ShippingFee:       m.ShippingFee,
		PayAmount:         m.PayAmount,
		Currency:          m.Currency,
		MerchantRemark:    m.MerchantRemark,
		Remark:            m.Remark,
		OriginalAmount:    m.OriginalAmount,
		AdjustAmount:      m.AdjustAmount,
		AdjustReason:      m.AdjustReason,
		AdjustedBy:        m.AdjustedBy,
		AdjustedAt:        adjustedAt,
		Version:           m.Version,
		PaymentMethod:     m.PaymentMethod,
		Source:            m.Source,
		ReceiverName:      m.ReceiverName,
		ReceiverPhone:     m.ReceiverPhone,
		ReceiverAddress:   m.ReceiverAddress,
		PaidAt:            paidAt,
		ShippedAt:         shippedAt,
		DeliveredAt:       deliveredAt,
		CancelledAt:       cancelledAt,
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0).UTC(),
			UpdatedAt: time.Unix(m.UpdatedAt, 0).UTC(),
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
		},
		DeletedAt: deletedAt,
	}
}

func fromOrderEntity(o *fulfillment.Order) *orderModel {
	var paidAt, shippedAt, deliveredAt, cancelledAt, adjustedAt, deletedAt *int64
	if o.PaidAt != nil {
		ts := o.PaidAt.Unix()
		paidAt = &ts
	}
	if o.ShippedAt != nil {
		ts := o.ShippedAt.Unix()
		shippedAt = &ts
	}
	if o.DeliveredAt != nil {
		ts := o.DeliveredAt.Unix()
		deliveredAt = &ts
	}
	if o.CancelledAt != nil {
		ts := o.CancelledAt.Unix()
		cancelledAt = &ts
	}
	if o.AdjustedAt != nil {
		ts := o.AdjustedAt.Unix()
		adjustedAt = &ts
	}
	if o.DeletedAt.Valid {
		ts := o.DeletedAt.Time.Unix()
		deletedAt = &ts
	}

	return &orderModel{
		ID:                o.ID,
		TenantID:          o.TenantID.Int64(),
		OrderNo:           o.OrderNo,
		UserID:            o.UserID,
		Status:            o.Status.String(),
		FulfillmentStatus: int8(o.FulfillmentStatus),
		RefundStatus:      int8(o.RefundStatus),
		TotalAmount:       o.TotalAmount,
		DiscountAmount:    o.DiscountAmount,
		ShippingFee:       o.ShippingFee,
		PayAmount:         o.PayAmount,
		Currency:          o.Currency,
		MerchantRemark:    o.MerchantRemark,
		Remark:            o.Remark,
		OriginalAmount:    o.OriginalAmount,
		AdjustAmount:      o.AdjustAmount,
		AdjustReason:      o.AdjustReason,
		AdjustedBy:        o.AdjustedBy,
		AdjustedAt:        adjustedAt,
		Version:           o.Version,
		PaymentMethod:     o.PaymentMethod,
		Source:            o.Source,
		ReceiverName:      o.ReceiverName,
		ReceiverPhone:     o.ReceiverPhone,
		ReceiverAddress:   o.ReceiverAddress,
		PaidAt:            paidAt,
		ShippedAt:         shippedAt,
		DeliveredAt:       deliveredAt,
		CancelledAt:       cancelledAt,
		CreatedBy:         o.Audit.CreatedBy,
		UpdatedBy:         o.Audit.UpdatedBy,
		DeletedAt:         deletedAt,
		CreatedAt:         o.Audit.CreatedAt.Unix(),
		UpdatedAt:         o.Audit.UpdatedAt.Unix(),
	}
}

// FindByID 根据ID查询订单
func (r *orderRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*fulfillment.Order, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model orderModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrOrderNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByOrderNo 根据订单号查询订单
func (r *orderRepo) FindByOrderNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderNo string) (*fulfillment.Order, error) {
	query := db.WithContext(ctx).Model(&orderModel{}).Where("order_no = ? AND deleted_at IS NULL", orderNo)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model orderModel
	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrOrderNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindList 分页查询订单列表
func (r *orderRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query fulfillment.OrderQuery) ([]*fulfillment.Order, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&orderModel{}).Where("deleted_at IS NULL")

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.OrderNo != "" {
		dbQuery = dbQuery.Where("order_no LIKE ?", escapeLikePattern(query.OrderNo))
	}
	if query.UserID != 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.FulfillmentStatus != 0 {
		dbQuery = dbQuery.Where("fulfillment_status = ?", query.FulfillmentStatus)
	}
	if query.RefundStatus != 0 {
		dbQuery = dbQuery.Where("refund_status = ?", query.RefundStatus)
	}
	if !query.StartTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime.Unix())
	}
	if !query.EndTime.IsZero() {
		dbQuery = dbQuery.Where("created_at < ?", query.EndTime.Unix())
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []orderModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	orders := make([]*fulfillment.Order, len(models))
	for i, m := range models {
		orders[i] = m.toEntity()
	}
	return orders, total, nil
}

// UpdateWithVersion 带乐观锁的更新
func (r *orderRepo) UpdateWithVersion(ctx context.Context, db *gorm.DB, order *fulfillment.Order) error {
	model := fromOrderEntity(order)

	// Use optimistic lock: update only if version matches
	result := db.WithContext(ctx).Model(&orderModel{}).
		Where("id = ? AND tenant_id = ? AND version = ? AND deleted_at IS NULL",
			order.ID, order.TenantID.Int64(), order.Version).
		Updates(map[string]interface{}{
			"pay_amount":         model.PayAmount,
			"adjust_amount":      model.AdjustAmount,
			"adjust_reason":      model.AdjustReason,
			"adjusted_by":        model.AdjustedBy,
			"adjusted_at":        model.AdjustedAt,
			"merchant_remark":    model.MerchantRemark,
			"version":            gorm.Expr("version + 1"),
			"updated_by":         model.UpdatedBy,
			"updated_at":         model.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// Version conflict or record not found
		return code.ErrOrderVersionConflict
	}

	// Update version in entity
	order.Version++
	return nil
}

// UpdateRemark 更新商家备注
func (r *orderRepo) UpdateRemark(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64, remark string) error {
	result := db.WithContext(ctx).Model(&orderModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", orderID, tenantID.Int64()).
		Update("merchant_remark", remark)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return code.ErrOrderNotFound
	}

	return nil
}

// CountTodayOrders 统计今日订单数
func (r *orderRepo) CountTodayOrders(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, error) {
	// Get start and end of today in UTC
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := db.WithContext(ctx).Model(&orderModel{}).
		Where("created_at >= ? AND created_at < ? AND deleted_at IS NULL", startOfDay.Unix(), endOfDay.Unix())
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// SumTodayGMV 统计今日GMV（已支付订单的总金额）
func (r *orderRepo) SumTodayGMV(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, error) {
	// Get start and end of today in UTC
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := db.WithContext(ctx).Model(&orderModel{}).
		Where("created_at >= ? AND created_at < ? AND status IN ? AND deleted_at IS NULL",
			startOfDay.Unix(), endOfDay.Unix(),
			[]string{"paid", "shipped", "delivered"})
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var result struct {
		Total int64
	}
	err := query.Select("COALESCE(SUM(pay_amount), 0) as total").Scan(&result).Error
	return result.Total, err
}

// FindForExport 导出订单（最多10000条）
func (r *orderRepo) FindForExport(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query fulfillment.OrderQuery) ([]*fulfillment.Order, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&orderModel{}).Where("deleted_at IS NULL")

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.OrderNo != "" {
		dbQuery = dbQuery.Where("order_no LIKE ?", escapeLikePattern(query.OrderNo))
	}
	if query.UserID != 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.FulfillmentStatus != 0 {
		dbQuery = dbQuery.Where("fulfillment_status = ?", query.FulfillmentStatus)
	}
	if query.RefundStatus != 0 {
		dbQuery = dbQuery.Where("refund_status = ?", query.RefundStatus)
	}
	if !query.StartTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime.Unix())
	}
	if !query.EndTime.IsZero() {
		dbQuery = dbQuery.Where("created_at < ?", query.EndTime.Unix())
	}

	var models []orderModel
	err := dbQuery.Order("created_at DESC").
		Limit(10000).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	orders := make([]*fulfillment.Order, len(models))
	for i, m := range models {
		orders[i] = m.toEntity()
	}
	return orders, nil
}

// ==================== OrderItem Repository ====================

type orderItemRepo struct{}

func NewOrderItemRepository() fulfillment.OrderItemRepository {
	return &orderItemRepo{}
}

type orderItemModel struct {
	ID          int64  `gorm:"column:id;primaryKey"`
	TenantID    int64  `gorm:"column:tenant_id;not null;index"`
	OrderID     int64  `gorm:"column:order_id;not null;index"`
	ProductID   int64  `gorm:"column:product_id;not null;index"`
	SKUID       int64  `gorm:"column:sku_id;not null;index"`
	ProductName string `gorm:"column:product_name;not null"`
	SKUName     string `gorm:"column:sku_name;not null"`
	Image       string `gorm:"column:image"`
	Quantity    int    `gorm:"column:quantity;not null"`
	UnitPrice   int64  `gorm:"column:unit_price;not null"`
	TotalPrice  int64  `gorm:"column:total_price;not null"`
	Currency    string `gorm:"column:currency;not null;default:'CNY'"`
	CreatedAt   int64  `gorm:"column:created_at;not null"`
}

func (orderItemModel) TableName() string {
	return "order_items"
}

func (m *orderItemModel) toEntity() fulfillment.OrderItem {
	return fulfillment.OrderItem{
		ID:          m.ID,
		TenantID:    shared.TenantID(m.TenantID),
		OrderID:     m.OrderID,
		ProductID:   m.ProductID,
		SKUID:       m.SKUID,
		ProductName: m.ProductName,
		SKUName:     m.SKUName,
		Image:       m.Image,
		Quantity:    m.Quantity,
		UnitPrice:   m.UnitPrice,
		TotalPrice:  m.TotalPrice,
		Currency:    m.Currency,
		CreatedAt:   time.Unix(m.CreatedAt, 0).UTC(),
	}
}

// FindByOrderID 根据订单ID查询明细
func (r *orderItemRepo) FindByOrderID(ctx context.Context, db *gorm.DB, orderID int64) ([]fulfillment.OrderItem, error) {
	var models []orderItemModel
	err := db.WithContext(ctx).Model(&orderItemModel{}).
		Where("order_id = ?", orderID).
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := make([]fulfillment.OrderItem, len(models))
	for i, m := range models {
		items[i] = m.toEntity()
	}
	return items, nil
}

// FindByOrderIDs 批量查询订单明细
func (r *orderItemRepo) FindByOrderIDs(ctx context.Context, db *gorm.DB, orderIDs []int64) (map[int64][]fulfillment.OrderItem, error) {
	if len(orderIDs) == 0 {
		return make(map[int64][]fulfillment.OrderItem), nil
	}

	var models []orderItemModel
	err := db.WithContext(ctx).Model(&orderItemModel{}).
		Where("order_id IN ?", orderIDs).
		Order("order_id, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make(map[int64][]fulfillment.OrderItem)
	for _, m := range models {
		result[m.OrderID] = append(result[m.OrderID], m.toEntity())
	}
	return result, nil
}