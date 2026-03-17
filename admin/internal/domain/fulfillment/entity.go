package fulfillment

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

var (
	ErrShipmentNotFound = errors.New("shipment not found")
	ErrInvalidTracking  = errors.New("invalid tracking number")
	ErrAlreadyShipped   = errors.New("order already shipped")
)

type ShipmentStatus int

const (
	ShipmentStatusPending ShipmentStatus = iota
	ShipmentStatusShipped
	ShipmentStatusInTransit
	ShipmentStatusDelivered
	ShipmentStatusFailed
)

type Shipment struct {
	ID          int64
	TenantID    shared.TenantID
	OrderID     string
	Status      ShipmentStatus
	Carrier     string
	TrackingNo  string
	Items       []ShipmentItem
	Weight      float64
	Cost        shared.Money
	ShippedAt   *time.Time
	DeliveredAt *time.Time
	Audit       shared.AuditInfo
}

func (s *Shipment) TableName() string {
	return "shipments"
}

func (s *Shipment) Ship(carrier, trackingNo string) error {
	if s.Status != ShipmentStatusPending {
		return ErrAlreadyShipped
	}
	if carrier == "" || trackingNo == "" {
		return ErrInvalidTracking
	}
	s.Carrier = carrier
	s.TrackingNo = trackingNo
	now := time.Now().UTC()
	s.ShippedAt = &now
	s.Status = ShipmentStatusShipped
	return nil
}

func (s *Shipment) Deliver() {
	now := time.Now().UTC()
	s.DeliveredAt = &now
	s.Status = ShipmentStatusDelivered
}

type ShipmentItem struct {
	ID          int64
	ShipmentID  int64
	OrderItemID int64
	ProductID   int64
	SKUId       int64
	Quantity    int
}

func (si *ShipmentItem) TableName() string {
	return "shipment_items"
}

type RefundStatus int

const (
	RefundStatusPending RefundStatus = iota
	RefundStatusApproved
	RefundStatusRejected
	RefundStatusCompleted
)

type Refund struct {
	ID          int64
	TenantID    shared.TenantID
	OrderID     string
	UserID      int64
	Status      RefundStatus
	Reason      string
	Description string
	Images      []string
	Amount      shared.Money
	ApprovedAt  *time.Time
	CompletedAt *time.Time
	Audit       shared.AuditInfo
}

func (r *Refund) TableName() string {
	return "refunds"
}

func (r *Refund) Approve() {
	now := time.Now().UTC()
	r.ApprovedAt = &now
	r.Status = RefundStatusApproved
}

func (r *Refund) Reject() {
	r.Status = RefundStatusRejected
}

func (r *Refund) Complete() {
	now := time.Now().UTC()
	r.CompletedAt = &now
	r.Status = RefundStatusCompleted
}

type ShipmentRepository interface {
	Create(ctx context.Context, db *gorm.DB, shipment *Shipment) error
	Update(ctx context.Context, db *gorm.DB, shipment *Shipment) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Shipment, error)
	FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID string) (*Shipment, error)
	FindByTrackingNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, trackingNo string) (*Shipment, error)
}

type RefundRepository interface {
	Create(ctx context.Context, db *gorm.DB, refund *Refund) error
	Update(ctx context.Context, db *gorm.DB, refund *Refund) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Refund, error)
	FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID string) ([]*Refund, error)
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query Query) ([]*Refund, int64, error)
}

type Query struct {
	shared.PageQuery
	Status RefundStatus
}
