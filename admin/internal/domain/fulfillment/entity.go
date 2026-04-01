package fulfillment

import (
	"context"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ==================== Enums ====================

// ShipmentStatus 发货单状态
type ShipmentStatus int

const (
	ShipmentStatusPending ShipmentStatus = iota // 待发货
	ShipmentStatusShipped                       // 已发货
	ShipmentStatusInTransit                     // 运输中
	ShipmentStatusDelivered                     // 已送达
	ShipmentStatusFailed                        // 发货失败
	ShipmentStatusCancelled                     // 已取消
)

func (s ShipmentStatus) String() string {
	switch s {
	case ShipmentStatusPending:
		return "pending"
	case ShipmentStatusShipped:
		return "shipped"
	case ShipmentStatusInTransit:
		return "in_transit"
	case ShipmentStatusDelivered:
		return "delivered"
	case ShipmentStatusFailed:
		return "failed"
	case ShipmentStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

func (s ShipmentStatus) IsValid() bool {
	return s >= ShipmentStatusPending && s <= ShipmentStatusCancelled
}

// ParseShipmentStatus parses a string to ShipmentStatus
func ParseShipmentStatus(s string) ShipmentStatus {
	switch s {
	case "pending":
		return ShipmentStatusPending
	case "shipped":
		return ShipmentStatusShipped
	case "in_transit":
		return ShipmentStatusInTransit
	case "delivered":
		return ShipmentStatusDelivered
	case "failed":
		return ShipmentStatusFailed
	case "cancelled":
		return ShipmentStatusCancelled
	default:
		return ShipmentStatusPending
	}
}

// RefundStatus 退款状态
type RefundStatus int

const (
	RefundStatusPending RefundStatus = iota // 待处理
	RefundStatusApproved                    // 已批准
	RefundStatusRejected                    // 已拒绝
	RefundStatusCompleted                   // 已完成
	RefundStatusCancelled                   // 已取消
)

func (s RefundStatus) String() string {
	switch s {
	case RefundStatusPending:
		return "pending"
	case RefundStatusApproved:
		return "approved"
	case RefundStatusRejected:
		return "rejected"
	case RefundStatusCompleted:
		return "completed"
	case RefundStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

func (s RefundStatus) IsValid() bool {
	return s >= RefundStatusPending && s <= RefundStatusCancelled
}

// ParseRefundStatus parses a string to RefundStatus
func ParseRefundStatus(s string) RefundStatus {
	switch s {
	case "pending":
		return RefundStatusPending
	case "approved":
		return RefundStatusApproved
	case "rejected":
		return RefundStatusRejected
	case "completed":
		return RefundStatusCompleted
	case "cancelled":
		return RefundStatusCancelled
	default:
		return RefundStatusPending
	}
}

// RefundType 退款类型
type RefundType int

const (
	RefundTypeFull RefundType = iota + 1 // 全额退款
	RefundTypePartial                    // 部分退款 (Phase 2)
)

func (t RefundType) String() string {
	switch t {
	case RefundTypeFull:
		return "full_refund"
	case RefundTypePartial:
		return "partial_refund"
	default:
		return "unknown"
	}
}

// ParseRefundType parses a string to RefundType
func ParseRefundType(s string) RefundType {
	switch s {
	case "full_refund":
		return RefundTypeFull
	case "partial_refund":
		return RefundTypePartial
	default:
		return RefundTypeFull
	}
}

// FulfillmentStatus 订单履约状态
type FulfillmentStatus int

const (
	FulfillmentStatusPending FulfillmentStatus = iota // 待发货
	FulfillmentStatusPartialShipped                   // 部分发货
	FulfillmentStatusShipped                          // 已发货
	FulfillmentStatusDelivered                        // 已送达
)

func (s FulfillmentStatus) String() string {
	switch s {
	case FulfillmentStatusPending:
		return "pending"
	case FulfillmentStatusPartialShipped:
		return "partial_shipped"
	case FulfillmentStatusShipped:
		return "shipped"
	case FulfillmentStatusDelivered:
		return "delivered"
	default:
		return "unknown"
	}
}

func (s FulfillmentStatus) IsValid() bool {
	return s >= FulfillmentStatusPending && s <= FulfillmentStatusDelivered
}

// ParseFulfillmentStatus parses a string to FulfillmentStatus
func ParseFulfillmentStatus(s string) FulfillmentStatus {
	switch s {
	case "pending":
		return FulfillmentStatusPending
	case "partial_shipped":
		return FulfillmentStatusPartialShipped
	case "shipped":
		return FulfillmentStatusShipped
	case "delivered":
		return FulfillmentStatusDelivered
	default:
		return FulfillmentStatusPending
	}
}

// ==================== Carrier (物流公司) ====================

// Carrier 物流公司
type Carrier struct {
	application.Model
	Code       string
	Name       string
	TrackingURL string
	IsActive   bool
	Sort       int
}

func (c *Carrier) TableName() string {
	return "carriers"
}

func (c *Carrier) GetTrackingURL(trackingNo string) string {
	if c.TrackingURL == "" || trackingNo == "" {
		return ""
	}
	return fmt.Sprintf(c.TrackingURL, trackingNo)
}

// PredefinedCarrierCodes 预定义物流公司代码
const (
	CarrierCodeSF   = "SF"   // 顺丰
	CarrierCodeYT   = "YT"   // 圆通
	CarrierCodeZT   = "ZT"   // 中通
	CarrierCodeST   = "ST"   // 申通
	CarrierCodeYD   = "YD"   // 韵达
	CarrierCodeEMS  = "EMS"  // EMS
	CarrierCodeJD   = "JD"   // 京东物流
	CarrierCodeOther = "OTHER" // 其他
)

// DefaultCarriers 默认物流公司列表
var DefaultCarriers = []Carrier{
	{Code: CarrierCodeSF, Name: "顺丰速运", TrackingURL: "https://www.sf-express.com/track?id=%s", IsActive: true, Sort: 1},
	{Code: CarrierCodeYT, Name: "圆通速递", TrackingURL: "https://www.yto.net.cn/query.html?id=%s", IsActive: true, Sort: 2},
	{Code: CarrierCodeZT, Name: "中通快递", TrackingURL: "https://www.zto.com/track?id=%s", IsActive: true, Sort: 3},
	{Code: CarrierCodeST, Name: "申通快递", TrackingURL: "https://www.sto.cn/track?id=%s", IsActive: true, Sort: 4},
	{Code: CarrierCodeYD, Name: "韵达快递", TrackingURL: "https://www.yundaex.com/track?id=%s", IsActive: true, Sort: 5},
	{Code: CarrierCodeEMS, Name: "EMS", TrackingURL: "https://www.ems.com.cn/track?id=%s", IsActive: true, Sort: 6},
	{Code: CarrierCodeJD, Name: "京东物流", TrackingURL: "https://www.jdl.com/track?id=%s", IsActive: true, Sort: 7},
	{Code: CarrierCodeOther, Name: "其他", TrackingURL: "", IsActive: true, Sort: 99},
}

// ==================== RefundReason (退款原因) ====================

// RefundReason 退款原因
type RefundReason struct {
	application.Model
	Code      string
	Name      string
	Sort      int
	IsActive  bool
}

func (r *RefundReason) TableName() string {
	return "refund_reasons"
}

// Predefined refund reason codes
const (
	RefundReasonDefective      = "DEFECTIVE"       // 产品有缺陷
	RefundReasonWrongItem      = "WRONG_ITEM"      // 发错货
	RefundReasonNotAsDescribed = "NOT_AS_DESCRIBED" // 与描述不符
	RefundReasonDamaged        = "DAMAGED"         // 运输损坏
	RefundReasonNotNeeded      = "NO_LONGER_NEEDED" // 不再需要
	RefundReasonLateDelivery   = "LATE_DELIVERY"   // 配送太慢
	RefundReasonOther          = "OTHER"           // 其他
)

// DefaultRefundReasons 默认退款原因列表
var DefaultRefundReasons = []RefundReason{
	{Code: RefundReasonDefective, Name: "产品有缺陷", Sort: 1, IsActive: true},
	{Code: RefundReasonWrongItem, Name: "发错货", Sort: 2, IsActive: true},
	{Code: RefundReasonNotAsDescribed, Name: "与描述不符", Sort: 3, IsActive: true},
	{Code: RefundReasonDamaged, Name: "运输损坏", Sort: 4, IsActive: true},
	{Code: RefundReasonNotNeeded, Name: "不再需要", Sort: 5, IsActive: true},
	{Code: RefundReasonLateDelivery, Name: "配送太慢", Sort: 6, IsActive: true},
	{Code: RefundReasonOther, Name: "其他", Sort: 99, IsActive: true},
}

// ==================== Shipment (发货单) ====================

// Shipment 发货单
type Shipment struct {
	application.Model
	TenantID         shared.TenantID `gorm:"column:tenant_id;not null;index"`
	OrderID          int64           `gorm:"column:order_id;not null;index"`
	ShipmentNo       string          `gorm:"column:shipment_no;not null;uniqueIndex:uk_shipment_no"`
	Status           ShipmentStatus  `gorm:"column:status;not null;default:0;index"`
	Carrier          string          `gorm:"column:carrier;not null;default:''"`
	CarrierCode      string          `gorm:"column:carrier_code;not null;default:''"`
	TrackingNo       string          `gorm:"column:tracking_no;not null;default:'';index"`
	ShippingCost     decimal.Decimal `gorm:"column:cost_amount;type:decimal(19,4);not null;default:0"`   // 运费成本
	ShippingCurrency string          `gorm:"column:cost_currency;not null;default:'CNY'"`
	Weight           decimal.Decimal `gorm:"column:weight;type:decimal(10,3);not null;default:0"` // 重量（kg）
	ShippedAt        *time.Time     `gorm:"column:shipped_at"`
	DeliveredAt      *time.Time     `gorm:"column:delivered_at"`
	CancelledAt      *time.Time     `gorm:"column:cancelled_at"`
	CancelledBy      int64          `gorm:"column:cancelled_by;not null;default:0"`
	CancelledReason  string         `gorm:"column:cancelled_reason;not null;default:''"`
	Remark           string          `gorm:"column:remark;not null;default:''"`
	Items            []ShipmentItem  `gorm:"foreignKey:ShipmentID"`
}

func (s *Shipment) TableName() string {
	return "shipments"
}

// NewShipment 创建发货单
func NewShipment(tenantID shared.TenantID, orderID int64, items []ShipmentItem, createdBy int64) *Shipment {
	return &Shipment{
		TenantID:         tenantID,
		OrderID:          orderID,
		Status:           ShipmentStatusPending,
		ShippingCurrency: "CNY",
		Items:            items,
	}
}

// Ship 发货
func (s *Shipment) Ship(carrier, carrierCode, trackingNo string, updatedBy int64) error {
	if s.Status != ShipmentStatusPending {
		return code.ErrShipmentAlreadyShipped
	}
	if carrier == "" {
		return code.ErrShipmentCarrierRequired
	}
	if trackingNo == "" {
		return code.ErrShipmentTrackingRequired
	}

	s.Carrier = carrier
	s.CarrierCode = carrierCode
	s.TrackingNo = trackingNo
	now := time.Now().UTC()
	s.ShippedAt = &now
	s.Status = ShipmentStatusShipped
	s.UpdatedAt = now
	return nil
}

// MarkInTransit 标记为运输中
func (s *Shipment) MarkInTransit(updatedBy int64) error {
	if s.Status != ShipmentStatusShipped {
		return code.ErrShipmentInvalidStatusTransition
	}
	s.Status = ShipmentStatusInTransit
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// MarkDelivered 标记为已送达
func (s *Shipment) MarkDelivered(updatedBy int64) error {
	if s.Status != ShipmentStatusShipped && s.Status != ShipmentStatusInTransit {
		return code.ErrShipmentInvalidStatusTransition
	}
	now := time.Now().UTC()
	s.DeliveredAt = &now
	s.Status = ShipmentStatusDelivered
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// MarkFailed 标记为发货失败
func (s *Shipment) MarkFailed(reason string, updatedBy int64) error {
	if s.Status == ShipmentStatusDelivered || s.Status == ShipmentStatusCancelled {
		return code.ErrShipmentInvalidStatusTransition
	}
	s.Status = ShipmentStatusFailed
	s.Remark = reason
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// Cancel 取消发货
// Deprecated: Use CancelShipment() instead which properly sets CancelledAt, CancelledBy, and CancelledReason
func (s *Shipment) Cancel(reason string, updatedBy int64) error {
	if s.Status == ShipmentStatusDelivered {
		return code.ErrShipmentCannotCancelDelivered
	}
	s.Status = ShipmentStatusCancelled
	s.Remark = reason
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// CancelShipment 取消发货单（带原因）
func (s *Shipment) CancelShipment(reason string, cancelledBy int64) error {
	if s.Status == ShipmentStatusDelivered {
		return code.ErrShipmentCannotCancelDelivered
	}
	if s.Status == ShipmentStatusCancelled {
		return code.ErrShipmentAlreadyCancelled
	}
	now := time.Now().UTC()
	s.Status = ShipmentStatusCancelled
	s.CancelledAt = &now
	s.CancelledBy = cancelledBy
	s.CancelledReason = reason
	s.UpdatedAt = now
	return nil
}

// SetShippingCost 设置运费
func (s *Shipment) SetShippingCost(cost decimal.Decimal, currency string) {
	s.ShippingCost = cost
	if currency != "" {
		s.ShippingCurrency = currency
	}
}

// SetWeight 设置重量
func (s *Shipment) SetWeight(weight decimal.Decimal) {
	s.Weight = weight
}

// ShipmentSequenceMod is the modulus for generating shipment sequence numbers
// This ensures the sequence number fits within 6 digits (000000-999999)
const ShipmentSequenceMod = 1000000

// GenerateShipmentNo 生成发货单号
func GenerateShipmentNo(tenantID shared.TenantID, sequence int) string {
	// Format: SHP{YYYYMMDD}{sequence:06d}
	dateStr := time.Now().UTC().Format("20060102")
	// Ensure sequence fits in 6 digits
	seq := sequence % ShipmentSequenceMod
	return fmt.Sprintf("SHP%s%06d", dateStr, seq)
}

// IsShipped 是否已发货
func (s *Shipment) IsShipped() bool {
	return s.Status == ShipmentStatusShipped ||
		s.Status == ShipmentStatusInTransit ||
		s.Status == ShipmentStatusDelivered
}

// IsDelivered 是否已送达
func (s *Shipment) IsDelivered() bool {
	return s.Status == ShipmentStatusDelivered
}

// CanShip 是否可以发货
func (s *Shipment) CanShip() bool {
	return s.Status == ShipmentStatusPending
}

// ==================== ShipmentItem (发货单明细) ====================

// ShipmentItem 发货单明细
type ShipmentItem struct {
	Model        application.Model
	TenantID     shared.TenantID `gorm:"column:tenant_id;not null;index"`
	ShipmentID   int64           `gorm:"column:shipment_id;not null;index"`
	OrderItemID  int64           `gorm:"column:order_item_id;not null;index"`
	ProductID    int64           `gorm:"column:product_id;not null;index"`
	SKUID        int64           `gorm:"column:sku_id;not null;index"`
	ProductName  string          `gorm:"column:product_name;not null;default:''"` // 商品名称快照
	SKUName      string          `gorm:"column:sku_name;not null;default:''"`     // SKU名称快照
	Image        string          `gorm:"column:image;not null;default:''"`        // 商品图片快照
	Quantity     int             `gorm:"column:quantity;not null;default:1"`
}

func (si *ShipmentItem) TableName() string {
	return "shipment_items"
}

// NewShipmentItem 创建发货单明细
func NewShipmentItem(tenantID shared.TenantID, orderItemID, productID, skuID int64,
	productName, skuName, image string, quantity int) *ShipmentItem {
	return &ShipmentItem{
		Model: application.Model{
			CreatedAt: time.Now().UTC(),
		},
		TenantID:    tenantID,
		OrderItemID: orderItemID,
		ProductID:   productID,
		SKUID:       skuID,
		ProductName: productName,
		SKUName:     skuName,
		Image:       image,
		Quantity:    quantity,
	}
}

// ==================== Refund (退款单) ====================

// Refund 退款单
type Refund struct {
	application.Model
	TenantID     shared.TenantID `gorm:"column:tenant_id;not null;index"`
	OrderID      int64         `gorm:"column:order_id;not null;index"`
	RefundNo     string         `gorm:"column:refund_no;not null;uniqueIndex:uk_refund_no"`
	UserID       int64          `gorm:"column:user_id;not null;index"`
	Type         RefundType     `gorm:"column:type;not null;default:1"`
	Status       RefundStatus   `gorm:"column:status;not null;default:0;index"`
	ReasonType   string         `gorm:"column:reason_type;not null;default:''"`
	Reason       string         `gorm:"column:reason;not null;default:''"`
	Description  string         `gorm:"column:description"`
	Images       string         `gorm:"column:images;type:json"` // JSON array of image URLs
	Amount       decimal.Decimal `gorm:"column:amount;type:decimal(19,4);not null;default:0"`       // 退款金额
	Currency     string         `gorm:"column:currency;not null;default:'CNY'"`
	RejectReason string         `gorm:"column:reject_reason;not null;default:''"`
	ApprovedAt   *time.Time   `gorm:"column:approved_at"`
	ApprovedBy   int64          `gorm:"column:approved_by;not null;default:0"`
	CompletedAt  *time.Time   `gorm:"column:completed_at"`
}

func (r *Refund) TableName() string {
	return "refunds"
}

// NewRefund 创建退款申请
func NewRefund(tenantID shared.TenantID, orderID int64, userID int64,
	reasonType, reason, description string, images []string, amount decimal.Decimal, currency string) *Refund {
	return &Refund{
		TenantID:    tenantID,
		OrderID:     orderID,
		UserID:      userID,
		Type:        RefundTypeFull,
		Status:      RefundStatusPending,
		ReasonType:  reasonType,
		Reason:      reason,
		Description: description,
		Amount:      amount,
		Currency:    currency,
	}
}

// Approve 批准退款
func (r *Refund) Approve(approvedBy int64) error {
	if r.Status != RefundStatusPending {
		return code.ErrRefundInvalidStatus
	}
	now := time.Now().UTC()
	r.Status = RefundStatusApproved
	r.ApprovedAt = &now
	r.ApprovedBy = approvedBy
	r.UpdatedAt = now
	return nil
}

// Reject 拒绝退款
func (r *Refund) Reject(rejectReason string, updatedBy int64) error {
	if r.Status != RefundStatusPending {
		return code.ErrRefundInvalidStatus
	}
	r.Status = RefundStatusRejected
	r.RejectReason = rejectReason
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// Complete 完成退款
func (r *Refund) Complete(updatedBy int64) error {
	if r.Status != RefundStatusApproved {
		return code.ErrRefundInvalidStatus
	}
	now := time.Now().UTC()
	r.Status = RefundStatusCompleted
	r.CompletedAt = &now
	r.UpdatedAt = now
	return nil
}

// Cancel 取消退款申请
func (r *Refund) Cancel(updatedBy int64) error {
	if r.Status != RefundStatusPending {
		return code.ErrRefundCannotCancel
	}
	r.Status = RefundStatusCancelled
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// GenerateRefundNo 生成退款单号
func GenerateRefundNo(tenantID shared.TenantID, sequence int) string {
	// Format: REF{YYYYMMDD}{sequence:06d}
	dateStr := time.Now().UTC().Format("20060102")
	return fmt.Sprintf("REF%s%06d", dateStr, sequence)
}

// IsPending 是否待处理
func (r *Refund) IsPending() bool {
	return r.Status == RefundStatusPending
}

// IsApproved 是否已批准
func (r *Refund) IsApproved() bool {
	return r.Status == RefundStatusApproved
}

// IsCompleted 是否已完成
func (r *Refund) IsCompleted() bool {
	return r.Status == RefundStatusCompleted
}

// IsRejected 是否已拒绝
func (r *Refund) IsRejected() bool {
	return r.Status == RefundStatusRejected
}

// CanCancel 是否可以取消
func (r *Refund) CanCancel() bool {
	return r.Status == RefundStatusPending
}

// ==================== Query Types ====================

// ShipmentQuery 发货单查询参数
type ShipmentQuery struct {
	shared.PageQuery
	OrderID    int64
	Status     ShipmentStatus
	CarrierCode string
	TrackingNo string
	StartTime  time.Time
	EndTime    time.Time
}

// RefundQuery 退款查询参数
type RefundQuery struct {
	shared.PageQuery
	RefundNo   string
	OrderID    int64
	UserID     int64
	Status     RefundStatus
	ReasonType string
	StartTime  time.Time
	EndTime    time.Time
}

// FulfillmentSummary 履约统计摘要
type FulfillmentSummary struct {
	PendingShipment  int64 `json:"pending_shipment"`  // 待发货订单数
	PartialShipped   int64 `json:"partial_shipped"`   // 部分发货订单数
	Shipped          int64 `json:"shipped"`           // 已发货订单数
	Delivered        int64 `json:"delivered"`         // 已送达订单数
	PendingRefund    int64 `json:"pending_refund"`    // 待处理退款数
}

// RefundStatistics 退款统计
type RefundStatistics struct {
	TotalRefunds      int64           `json:"total_refunds"`        // 总退款数
	TotalRefundAmount decimal.Decimal `json:"total_refund_amount"`  // 总退款金额
	RefundRate        float64         `json:"refund_rate"`          // 退款率
	TopReasons        []RefundReasonCount `json:"top_reasons"`     // 热门退款原因
	TopProducts       []ProductRefundCount `json:"top_products"`  // 高退款率商品
}

// RefundReasonCount 退款原因统计
type RefundReasonCount struct {
	ReasonType string `json:"reason_type"`
	ReasonName string `json:"reason_name"`
	Count      int64  `json:"count"`
	Percentage float64 `json:"percentage"`
}

// ProductRefundCount 商品退款统计
type ProductRefundCount struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	RefundCount int64   `json:"refund_count"`
	RefundRate  float64 `json:"refund_rate"`
}

// ==================== Repository Interfaces ====================

// ShipmentRepository 发货单仓储接口
type ShipmentRepository interface {
	Create(ctx context.Context, db *gorm.DB, shipment *Shipment) error
	Update(ctx context.Context, db *gorm.DB, shipment *Shipment) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Shipment, error)
	FindByShipmentNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, shipmentNo string) (*Shipment, error)
	FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) ([]*Shipment, error)
	FindByTrackingNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, trackingNo string) (*Shipment, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query ShipmentQuery) ([]*Shipment, int64, error)
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	CountByStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, status ShipmentStatus) (int64, error)
}

// ShipmentItemRepository 发货单明细仓储接口
type ShipmentItemRepository interface {
	BatchCreate(ctx context.Context, db *gorm.DB, items []ShipmentItem) error
	FindByShipmentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, shipmentID int64) ([]ShipmentItem, error)
	FindByShipmentIDs(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, shipmentIDs []int64) (map[int64][]ShipmentItem, error)
	FindByOrderItemID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderItemID int64) ([]ShipmentItem, error)
	DeleteByShipmentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, shipmentID int64) error
}

// RefundRepository 退款单仓储接口
type RefundRepository interface {
	Create(ctx context.Context, db *gorm.DB, refund *Refund) error
	Update(ctx context.Context, db *gorm.DB, refund *Refund) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Refund, error)
	FindByRefundNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, refundNo string) (*Refund, error)
	FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) ([]*Refund, error)
	FindPendingByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) (*Refund, error)
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query RefundQuery) ([]*Refund, int64, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query RefundQuery) ([]*Refund, int64, error)
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	CountByStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, status RefundStatus) (int64, error)
}

// CarrierRepository 物流公司仓储接口
type CarrierRepository interface {
	Create(ctx context.Context, db *gorm.DB, carrier *Carrier) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*Carrier, error)
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*Carrier, error)
	FindByCodes(ctx context.Context, db *gorm.DB, codes []string) (map[string]*Carrier, error)
	FindAll(ctx context.Context, db *gorm.DB) ([]*Carrier, error)
	FindActive(ctx context.Context, db *gorm.DB) ([]*Carrier, error)
	Update(ctx context.Context, db *gorm.DB, carrier *Carrier) error
}

// RefundReasonRepository 退款原因仓储接口
type RefundReasonRepository interface {
	Create(ctx context.Context, db *gorm.DB, reason *RefundReason) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*RefundReason, error)
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*RefundReason, error)
	FindAll(ctx context.Context, db *gorm.DB) ([]*RefundReason, error)
	FindActive(ctx context.Context, db *gorm.DB) ([]*RefundReason, error)
	Update(ctx context.Context, db *gorm.DB, reason *RefundReason) error
}