package fulfillment

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"gorm.io/gorm"
)

// OrderValidationInfo represents order information needed for shipment validation
type OrderValidationInfo struct {
	OrderID           string
	Status            string        // Order status: pending_payment, paid, shipped, etc.
	FulfillmentStatus int8          // Fulfillment status: pending, partial_shipped, shipped, delivered
	IsPaid            bool          // Whether order is paid
	Items             []OrderItemInfo
}

// OrderItemInfo represents order item information for shipment validation
type OrderItemInfo struct {
	OrderItemID int64
	ProductID   int64
	SKUID       int64
	Quantity    int
	ShippedQty  int // Already shipped quantity
}

// OrderValidator validates order before shipment creation
// This interface should be implemented by the Order service client
type OrderValidator interface {
	// GetOrderForShipment retrieves order info for shipment validation
	GetOrderForShipment(ctx context.Context, tenantID shared.TenantID, orderID string) (*OrderValidationInfo, error)
	// ValidateShipmentItems validates that shipment items match order items and quantities
	ValidateShipmentItems(order *OrderValidationInfo, items []CreateShipmentItemRequest) error
}

// NopOrderValidator is a no-op validator for use when order service is not available
// In production, this should be replaced with actual order service integration
type NopOrderValidator struct{}

func (n *NopOrderValidator) GetOrderForShipment(ctx context.Context, tenantID shared.TenantID, orderID string) (*OrderValidationInfo, error) {
	// Return nil to skip validation when order service is not integrated
	return nil, nil
}

func (n *NopOrderValidator) ValidateShipmentItems(order *OrderValidationInfo, items []CreateShipmentItemRequest) error {
	return nil
}

// OrderFulfillmentItem 订单履约明细
type OrderFulfillmentItem struct {
	OrderItemID int64
	ProductID   int64
	SKUID       int64
	ProductName string
	SKUName     string
	Image       string
	Quantity    int
	ShippedQty  int
	PendingQty  int
	UnitPrice   int64
	Currency    string
}

// OrderFulfillmentDetail 订单履约详情
type OrderFulfillmentDetail struct {
	OrderID           string
	OrderNo           string
	Status            string
	FulfillmentStatus int8
	FulfillmentText   string
	RefundStatus      int8
	RefundText        string
	TotalAmount       int64
	Currency          string
	UserID            int64
	UserName          string
	UserPhone         string
	ShippingAddress   string
	Items             []*OrderFulfillmentItem
	Shipments         []*ShipmentResponse
	Refund            *RefundResponse
	PaidAt            *time.Time
	ShippedAt         *time.Time
	DeliveredAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// QueryOrderRequest 查询订单请求
type QueryOrderRequest struct {
	Page              int
	PageSize          int
	OrderNo           string
	UserID            int64
	UserName          string
	Status            string
	FulfillmentStatus int8
	RefundStatus      int8
	StartTime         time.Time
	EndTime           time.Time
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	List     []*OrderFulfillmentDetail `json:"list"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
}

// RefundResponse 退款响应
type RefundResponse struct {
	ID             int64
	RefundNo       string
	OrderID        string
	UserID         int64
	UserName       string
	Type           int
	TypeText       string
	Status         int
	StatusText     string
	ReasonType     string
	Reason         string
	Description    string
	Images         []string
	Amount         int64
	Currency       string
	RejectReason   string
	ApprovedAt     *time.Time
	ApprovedBy     int64
	ApprovedByName string
	CompletedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	OrderNo        string
	OrderAmount    int64
}

// FulfillmentSummary 履约摘要
type FulfillmentSummary struct {
	PendingShipment int64
	PartialShipped  int64
	Shipped         int64
	Delivered       int64
	PendingRefund   int64
	Refunding       int64
	TotalOrders     int64
	TodayOrders     int64
	TodayGMV        int64
}

// OrderFulfillmentApp 订单履约应用服务接口
type OrderFulfillmentApp interface {
	ListOrders(ctx context.Context, tenantID shared.TenantID, req QueryOrderRequest) (*OrderListResponse, error)
	GetOrderFulfillment(ctx context.Context, tenantID shared.TenantID, orderID string) (*OrderFulfillmentDetail, error)
	ShipOrder(ctx context.Context, tenantID shared.TenantID, userID int64, orderID string, req ShipOrderRequest) (*ShipmentResponse, error)
	GetFulfillmentSummary(ctx context.Context, tenantID shared.TenantID) (*FulfillmentSummary, error)
	// New methods
	UpdateOrderRemark(ctx context.Context, tenantID shared.TenantID, orderID int64, remark string) error
	AdjustOrderPrice(ctx context.Context, tenantID shared.TenantID, userID int64, orderID int64, adjustAmount int64, reason string) (*fulfillment.AdjustPriceResponse, error)
	ExportOrders(ctx context.Context, tenantID shared.TenantID, req QueryOrderRequest) ([]*fulfillment.OrderExportRow, int64, error)
}

// ShipOrderRequest 发货请求
type ShipOrderRequest struct {
	CarrierCode  string
	CarrierName  string
	TrackingNo   string
	ShippingCost int64
	Currency     string
	Weight       float64
	Remark       string
	Items        []CreateShipmentItemRequest
}

type orderFulfillmentApp struct {
	db               *gorm.DB
	shipmentRepo     fulfillment.ShipmentRepository
	shipmentItemRepo fulfillment.ShipmentItemRepository
	carrierRepo      fulfillment.CarrierRepository
	refundRepo       fulfillment.RefundRepository
	orderRepo        fulfillment.OrderRepository
	orderItemRepo    fulfillment.OrderItemRepository
	idGen            snowflake.Snowflake
	orderValidator   OrderValidator
}

// NewOrderFulfillmentApp creates a new order fulfillment application service
func NewOrderFulfillmentApp(
	db *gorm.DB,
	shipmentRepo fulfillment.ShipmentRepository,
	shipmentItemRepo fulfillment.ShipmentItemRepository,
	carrierRepo fulfillment.CarrierRepository,
	refundRepo fulfillment.RefundRepository,
	orderRepo fulfillment.OrderRepository,
	orderItemRepo fulfillment.OrderItemRepository,
	idGen snowflake.Snowflake,
	orderValidator OrderValidator,
) OrderFulfillmentApp {
	return &orderFulfillmentApp{
		db:               db,
		shipmentRepo:     shipmentRepo,
		shipmentItemRepo: shipmentItemRepo,
		carrierRepo:      carrierRepo,
		refundRepo:       refundRepo,
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,
		idGen:            idGen,
		orderValidator:   orderValidator,
	}
}

func (a *orderFulfillmentApp) ListOrders(ctx context.Context, tenantID shared.TenantID, req QueryOrderRequest) (*OrderListResponse, error) {
	// This is a placeholder implementation
	// In a real implementation, this would query the order service/database
	// For now, we return an empty list
	return &OrderListResponse{
		List:     []*OrderFulfillmentDetail{},
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (a *orderFulfillmentApp) GetOrderFulfillment(ctx context.Context, tenantID shared.TenantID, orderID string) (*OrderFulfillmentDetail, error) {
	// Get shipments for the order
	shipments, err := a.shipmentRepo.FindByOrderID(ctx, a.db, tenantID, orderID)
	if err != nil {
		return nil, err
	}

	// Get refunds for the order
	refunds, err := a.refundRepo.FindByOrderID(ctx, a.db, tenantID, orderID)
	if err != nil {
		return nil, err
	}

	// Build response
	detail := &OrderFulfillmentDetail{
		OrderID: orderID,
	}

	// Convert shipments
	detail.Shipments = make([]*ShipmentResponse, len(shipments))
	for i, s := range shipments {
		items, _ := a.shipmentItemRepo.FindByShipmentID(ctx, a.db, tenantID, s.ID)
		s.Items = items
		carrier, _ := a.carrierRepo.FindByCode(ctx, a.db, s.CarrierCode)
		detail.Shipments[i] = toShipmentResponse(s, carrier)
	}

	// Convert refunds
	if len(refunds) > 0 {
		// Get the latest refund
		latestRefund := refunds[0]
		for _, r := range refunds {
			if r.CreatedAt.After(latestRefund.CreatedAt) {
				latestRefund = r
			}
		}
		detail.Refund = toRefundResponse(latestRefund)
	}

	return detail, nil
}

func (a *orderFulfillmentApp) ShipOrder(ctx context.Context, tenantID shared.TenantID, userID int64, orderID string, req ShipOrderRequest) (*ShipmentResponse, error) {
	// Validate order exists and can be shipped
	if a.orderValidator != nil {
		orderInfo, err := a.orderValidator.GetOrderForShipment(ctx, tenantID, orderID)
		if err != nil {
			return nil, err
		}

		// Validate order is paid (BR-001: Cannot ship unpaid orders)
		if !orderInfo.IsPaid {
			return nil, code.ErrRefundOrderNotPaid
		}

		// Validate order is not already fully shipped
		if orderInfo.FulfillmentStatus == 2 || orderInfo.FulfillmentStatus == 3 { // shipped or delivered
			return nil, code.ErrShipmentAlreadyShipped
		}

		// Validate shipment items (BR-003: Cannot ship more items than ordered)
		if len(req.Items) > 0 {
			if err := a.orderValidator.ValidateShipmentItems(orderInfo, req.Items); err != nil {
				return nil, code.ErrShipmentItemQuantityExceeded
			}
		}
	}

	// Validate carrier
	carrier, err := a.carrierRepo.FindByCode(ctx, a.db, req.CarrierCode)
	if err != nil {
		if errors.Is(err, code.ErrCarrierNotFound) {
			if req.CarrierName == "" {
				return nil, code.ErrShipmentCarrierRequired
			}
			carrier = &fulfillment.Carrier{
				Code: req.CarrierCode,
				Name: req.CarrierName,
			}
		} else {
			return nil, err
		}
	}

	// Validate tracking number
	if req.TrackingNo == "" {
		return nil, code.ErrShipmentTrackingRequired
	}

	var result *fulfillment.Shipment

	err = a.db.Transaction(func(tx *gorm.DB) error {
		// Generate shipment ID and number
		shipmentID, err := a.idGen.NextID(ctx)
		if err != nil {
			return err
		}

		shipmentNo := fulfillment.GenerateShipmentNo(tenantID, int(shipmentID%1000000))

		// Create shipment entity
		shipment := &fulfillment.Shipment{
			TenantID:         tenantID,
			OrderID:          orderID,
			ShipmentNo:       shipmentNo,
			Status:           fulfillment.ShipmentStatusShipped,
			Carrier:          carrier.Name,
			CarrierCode:      carrier.Code,
			TrackingNo:       req.TrackingNo,
			ShippingCost:     req.ShippingCost,
			ShippingCurrency: req.Currency,
			Weight:           req.Weight,
			Remark:           req.Remark,
		}

		// Set shipped at
		now := time.Now().UTC()
		shipment.ShippedAt = &now

		// Create shipment
		if err := a.shipmentRepo.Create(ctx, tx, shipment); err != nil {
			return err
		}

		// Create shipment items
		if len(req.Items) > 0 {
			items := make([]fulfillment.ShipmentItem, len(req.Items))
			for i, itemReq := range req.Items {
				itemID, err := a.idGen.NextID(ctx)
				if err != nil {
					return err
				}
				items[i] = fulfillment.ShipmentItem{
					ID:          itemID,
					TenantID:    tenantID,
					ShipmentID:  shipmentID,
					OrderItemID: itemReq.OrderItemID,
					ProductID:   itemReq.ProductID,
					SKUID:       itemReq.SKUID,
					ProductName: itemReq.ProductName,
					SKUName:     itemReq.SKUName,
					Image:       itemReq.Image,
					Quantity:    itemReq.Quantity,
					CreatedAt:   time.Now().UTC(),
				}
			}
			if err := a.shipmentItemRepo.BatchCreate(ctx, tx, items); err != nil {
				return err
			}
			shipment.Items = items
		}

		result = shipment
		return nil
	})

	if err != nil {
		return nil, err
	}

	return toShipmentResponse(result, carrier), nil
}

func (a *orderFulfillmentApp) GetFulfillmentSummary(ctx context.Context, tenantID shared.TenantID) (*FulfillmentSummary, error) {
	// Count shipments by status
	pendingCount, _ := a.shipmentRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.ShipmentStatusPending)
	shippedCount, _ := a.shipmentRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.ShipmentStatusShipped)
	inTransitCount, _ := a.shipmentRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.ShipmentStatusInTransit)
	deliveredCount, _ := a.shipmentRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.ShipmentStatusDelivered)

	// Count refunds by status
	pendingRefundCount, _ := a.refundRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.RefundStatusPending)
	approvedRefundCount, _ := a.refundRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.RefundStatusApproved)

	// Get today's stats
	todayOrders, _ := a.orderRepo.CountTodayOrders(ctx, a.db, tenantID)
	todayGMV, _ := a.orderRepo.SumTodayGMV(ctx, a.db, tenantID)

	return &FulfillmentSummary{
		PendingShipment: pendingCount,
		PartialShipped:  0, // Would need order-level data
		Shipped:         shippedCount + inTransitCount,
		Delivered:       deliveredCount,
		PendingRefund:   pendingRefundCount,
		Refunding:       approvedRefundCount,
		TotalOrders:     pendingCount + shippedCount + inTransitCount + deliveredCount,
		TodayOrders:     todayOrders,
		TodayGMV:        todayGMV,
	}, nil
}

// UpdateOrderRemark updates the merchant remark for an order
func (a *orderFulfillmentApp) UpdateOrderRemark(ctx context.Context, tenantID shared.TenantID, orderID int64, remark string) error {
	return a.orderRepo.UpdateRemark(ctx, a.db, tenantID, orderID, remark)
}

// AdjustOrderPrice adjusts the price of an order
func (a *orderFulfillmentApp) AdjustOrderPrice(ctx context.Context, tenantID shared.TenantID, userID int64, orderID int64, adjustAmount int64, reason string) (*fulfillment.AdjustPriceResponse, error) {
	// Get order
	order, err := a.orderRepo.FindByID(ctx, a.db, tenantID, orderID)
	if err != nil {
		return nil, err
	}

	// Adjust price (uses optimistic lock internally)
	err = order.AdjustPrice(adjustAmount, reason, userID)
	if err != nil {
		return nil, err
	}

	// Update with version check
	err = a.orderRepo.UpdateWithVersion(ctx, a.db, order)
	if err != nil {
		return nil, err
	}

	return &fulfillment.AdjustPriceResponse{
		OrderID:        order.ID,
		OriginalAmount: order.OriginalAmount,
		AdjustAmount:   order.AdjustAmount,
		NewPayAmount:   order.PayAmount,
		AdjustReason:   order.AdjustReason,
		AdjustedAt:     order.AdjustedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ExportOrders exports orders to CSV
func (a *orderFulfillmentApp) ExportOrders(ctx context.Context, tenantID shared.TenantID, req QueryOrderRequest) ([]*fulfillment.OrderExportRow, int64, error) {
	// Build query
	query := fulfillment.OrderQuery{
		PageQuery: shared.PageQuery{
			Page:     1,
			PageSize: 10001, // Check if exceeds limit
		},
		OrderNo:           req.OrderNo,
		UserID:            req.UserID,
		Status:            fulfillment.OrderStatus(req.Status),
		FulfillmentStatus: fulfillment.OrderFulfillmentStatus(req.FulfillmentStatus),
		RefundStatus:      fulfillment.OrderRefundStatus(req.RefundStatus),
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
	}

	// Get orders for export
	orders, err := a.orderRepo.FindForExport(ctx, a.db, tenantID, query)
	if err != nil {
		return nil, 0, err
	}

	if len(orders) > 10000 {
		return nil, 0, code.ErrOrderExportLimitExceed
	}

	// Convert to export rows
	rows := make([]*fulfillment.OrderExportRow, len(orders))
	for i, o := range orders {
		rows[i] = &fulfillment.OrderExportRow{
			OrderNo:          o.OrderNo,
			Status:           o.Status.Text(),
			FulfillmentStatus: o.FulfillmentStatus.Text(),
			RefundStatus:     o.RefundStatus.Text(),
			TotalAmount:      o.TotalAmount,
			DiscountAmount:   o.DiscountAmount,
			ShippingFee:      o.ShippingFee,
			PayAmount:        o.PayAmount,
			ReceiverName:     o.ReceiverName,
			ReceiverPhone:    o.ReceiverPhone,
			ReceiverAddress:  o.ReceiverAddress,
			PaymentMethod:    o.PaymentMethod,
			CreatedAt:        o.Audit.CreatedAt,
			PaidAt:           o.PaidAt,
		}
	}

	return rows, int64(len(rows)), nil
}

// toRefundResponse converts refund entity to response
func toRefundResponse(r *fulfillment.Refund) *RefundResponse {
	// Parse images JSON
	images := []string{}
	// In a real implementation, you would parse the JSON string

	return &RefundResponse{
		ID:           r.ID,
		RefundNo:     r.RefundNo,
		OrderID:      r.OrderID,
		UserID:       r.UserID,
		Type:         int(r.Type),
		TypeText:     r.Type.String(),
		Status:       int(r.Status),
		StatusText:   r.Status.String(),
		ReasonType:   r.ReasonType,
		Reason:       r.Reason,
		Description:  r.Description,
		Images:       images,
		Amount:       r.Amount,
		Currency:     r.Currency,
		RejectReason: r.RejectReason,
		ApprovedAt:   r.ApprovedAt,
		ApprovedBy:   r.ApprovedBy,
		CompletedAt:  r.CompletedAt,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}