package fulfillment

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// OrderValidationInfo represents order information needed for shipment validation
type OrderValidationInfo struct {
	OrderID           int64
	Status            string // Order status: pending_payment, paid, shipped, etc.
	FulfillmentStatus int8   // Fulfillment status: pending, partial_shipped, shipped, delivered
	IsPaid            bool   // Whether order is paid
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
	GetOrderForShipment(ctx context.Context,  orderID int64) (*OrderValidationInfo, error)
	// ValidateShipmentItems validates that shipment items match order items and quantities
	ValidateShipmentItems(order *OrderValidationInfo, items []CreateShipmentItemRequest) error
}

// DefaultOrderValidator is the default validator when order service is not available
// It returns an error indicating order validation is required but not available
type DefaultOrderValidator struct{}

// GetOrderForShipment returns error because order service is not integrated
func (v *DefaultOrderValidator) GetOrderForShipment(ctx context.Context, orderID int64) (*OrderValidationInfo, error) {
	return nil, code.ErrOrderServiceUnavailable
}

// ValidateShipmentItems returns error
func (v *DefaultOrderValidator) ValidateShipmentItems(order *OrderValidationInfo, items []CreateShipmentItemRequest) error {
	return code.ErrOrderServiceUnavailable
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
	UnitPrice   decimal.Decimal
	Currency    string
}

// OrderShippingAddress 订单收货地址
type OrderShippingAddress struct {
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	Address       string
	FullAddress   string
}

// OrderFulfillmentDetail 订单履约详情
type OrderFulfillmentDetail struct {
	OrderID           int64
	OrderNo           string
	Status            string
	FulfillmentStatus int8
	FulfillmentText   string
	RefundStatus      int8
	RefundText        string
	TotalAmount       decimal.Decimal
	Currency          string
	UserID            int64
	UserName          string
	UserPhone         string
	ShippingAddress   *OrderShippingAddress
	Items             []*OrderFulfillmentItem
	Shipments         []*ShipmentResponse
	Refund            *RefundResponse
	Remark            string
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
	OrderID        int64
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
	Amount         decimal.Decimal
	Currency       string
	RejectReason   string
	ApprovedAt     *time.Time
	ApprovedBy     int64
	ApprovedByName string
	CompletedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	OrderNo        string
	OrderAmount    decimal.Decimal
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
	TodayGMV        decimal.Decimal
}

// OrderFulfillmentApp 订单履约应用服务接口
type OrderFulfillmentApp interface {
	ListOrders(ctx context.Context,  req QueryOrderRequest) (*OrderListResponse, error)
	GetOrderFulfillment(ctx context.Context,  orderID int64) (*OrderFulfillmentDetail, error)
	ShipOrder(ctx context.Context,  userID int64, orderID int64, req ShipOrderRequest) (*ShipmentResponse, error)
	GetFulfillmentSummary(ctx context.Context) (*FulfillmentSummary, error)
	// New methods
	UpdateOrderRemark(ctx context.Context,  orderID int64, remark string) error
	AdjustOrderPrice(ctx context.Context,  userID int64, orderID int64, adjustAmount decimal.Decimal, reason string) (*fulfillment.AdjustPriceResponse, error)
	ExportOrders(ctx context.Context,  req QueryOrderRequest) ([]*fulfillment.OrderExportRow, int64, error)
}

// ShipOrderRequest 发货请求
type ShipOrderRequest struct {
	CarrierCode  string
	CarrierName  string
	TrackingNo   string
	ShippingCost decimal.Decimal
	Currency     string
	Weight       decimal.Decimal
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

func (a *orderFulfillmentApp) ListOrders(ctx context.Context,  req QueryOrderRequest) (*OrderListResponse, error) {
	orders, total, err := a.orderRepo.FindList(ctx, a.db,  buildOrderQuery(req))
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return &OrderListResponse{
			List:     []*OrderFulfillmentDetail{},
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, nil
	}

	orderIDs := collectOrderIDs(orders)
	itemsMap, err := a.orderItemRepo.FindByOrderIDs(ctx, a.db, orderIDs)
	if err != nil {
		return nil, err
	}

	shipmentMap := a.loadShipmentsByOrder(ctx,  orderIDs)
	shipmentItemsMap, carrierMap := a.preloadShipmentContexts(ctx,  shipmentMap)
	refundMap := a.loadRefundsByOrder(ctx,  orderIDs)

	list := make([]*OrderFulfillmentDetail, len(orders))
	for i, o := range orders {
		detail := toOrderFulfillmentDetail(o)
		detail.Items = toOrderFulfillmentItems(itemsMap[o.ID])
		detail.Shipments = buildShipmentResponses(shipmentMap[o.ID], shipmentItemsMap, carrierMap, o.OrderNo)
		detail.Refund = pickLatestRefund(refundMap[o.ID])
		list[i] = detail
	}

	return &OrderListResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// loadShipmentsByOrder loads shipments for each order. Returns map keyed by
// orderID. Errors for individual orders are swallowed to avoid failing the
// whole list. NOTE: This is an N+1 query pattern; fixing it requires adding
// ShipmentRepository.FindByOrderIDs (separate change, out of scope).
func (a *orderFulfillmentApp) loadShipmentsByOrder(ctx context.Context, orderIDs []int64) map[int64][]*fulfillment.Shipment {
	shipmentMap := make(map[int64][]*fulfillment.Shipment, len(orderIDs))
	for _, orderID := range orderIDs {
		shipments, err := a.shipmentRepo.FindByOrderID(ctx, a.db,  orderID)
		if err != nil {
			continue
		}
		if len(shipments) > 0 {
			shipmentMap[orderID] = shipments
		}
	}
	return shipmentMap
}

// preloadShipmentContexts batches the loading of shipment items (keyed by
// shipment ID) and carriers (keyed by carrier code) across all loaded shipments.
// Returns empty maps when there are no shipments.
func (a *orderFulfillmentApp) preloadShipmentContexts(
	ctx context.Context,
	shipmentMap map[int64][]*fulfillment.Shipment,
) (map[int64][]fulfillment.ShipmentItem, map[string]*fulfillment.Carrier) {
	shipmentItemsMap := make(map[int64][]fulfillment.ShipmentItem)
	carrierMap := make(map[string]*fulfillment.Carrier)

	if len(shipmentMap) == 0 {
		return shipmentItemsMap, carrierMap
	}

	shipmentIDs := make([]int64, 0)
	carrierCodes := make([]string, 0)
	for _, shipments := range shipmentMap {
		for _, s := range shipments {
			shipmentIDs = append(shipmentIDs, s.Model.ID)
			carrierCodes = append(carrierCodes, s.CarrierCode)
		}
	}
	shipmentItemsMap, _ = a.shipmentItemRepo.FindByShipmentIDs(ctx, a.db,  shipmentIDs)
	carrierMap, _ = a.carrierRepo.FindByCodes(ctx, a.db, carrierCodes)
	return shipmentItemsMap, carrierMap
}

// loadRefundsByOrder loads refunds per order. Returns map keyed by orderID;
// errors for individual orders are swallowed. Same N+1 caveat as loadShipmentsByOrder.
func (a *orderFulfillmentApp) loadRefundsByOrder(ctx context.Context, orderIDs []int64) map[int64][]*fulfillment.Refund {
	refundMap := make(map[int64][]*fulfillment.Refund, len(orderIDs))
	for _, orderID := range orderIDs {
		refunds, err := a.refundRepo.FindByOrderID(ctx, a.db,  orderID)
		if err != nil {
			continue
		}
		if len(refunds) > 0 {
			refundMap[orderID] = refunds
		}
	}
	return refundMap
}

func (a *orderFulfillmentApp) GetOrderFulfillment(ctx context.Context,  orderID int64) (*OrderFulfillmentDetail, error) {
	order, err := a.orderRepo.FindByID(ctx, a.db,  orderID)
	if err != nil {
		return nil, err
	}

	itemsMap, err := a.orderItemRepo.FindByOrderIDs(ctx, a.db, []int64{orderID})
	if err != nil {
		return nil, err
	}

	detail := toOrderFulfillmentDetail(order)
	detail.Items = toOrderFulfillmentItems(itemsMap[order.ID])

	if err := a.loadAndAttachShipments(ctx,  detail); err != nil {
		return nil, err
	}
	if err := a.loadAndAttachRefund(ctx,  detail); err != nil {
		return nil, err
	}
	return detail, nil
}

// loadAndAttachShipments loads shipments + their items + carriers, attaching them to detail.Shipments.
func (a *orderFulfillmentApp) loadAndAttachShipments(ctx context.Context, detail *OrderFulfillmentDetail) error {
	shipments, err := a.shipmentRepo.FindByOrderID(ctx, a.db,  detail.OrderID)
	if err != nil {
		return err
	}
	if len(shipments) == 0 {
		detail.Shipments = nil
		return nil
	}

	shipmentIDs := make([]int64, len(shipments))
	carrierCodes := make([]string, len(shipments))
	for i, s := range shipments {
		shipmentIDs[i] = s.Model.ID
		carrierCodes[i] = s.CarrierCode
	}
	shipmentItemsMap, _ := a.shipmentItemRepo.FindByShipmentIDs(ctx, a.db,  shipmentIDs)
	carrierMap, _ := a.carrierRepo.FindByCodes(ctx, a.db, carrierCodes)

	detail.Shipments = buildShipmentResponses(shipments, shipmentItemsMap, carrierMap, detail.OrderNo)
	return nil
}

// loadAndAttachRefund loads refunds for the order and attaches the latest one to detail.Refund.
func (a *orderFulfillmentApp) loadAndAttachRefund(ctx context.Context, detail *OrderFulfillmentDetail) error {
	refunds, err := a.refundRepo.FindByOrderID(ctx, a.db,  detail.OrderID)
	if err != nil {
		return err
	}
	detail.Refund = pickLatestRefund(refunds)
	return nil
}

func (a *orderFulfillmentApp) ShipOrder(ctx context.Context,  userID int64, orderID int64, req ShipOrderRequest) (*ShipmentResponse, error) {
	// Validate order exists and can be shipped
	if a.orderValidator != nil {
		orderInfo, err := a.orderValidator.GetOrderForShipment(ctx,  orderID)
		if err != nil {
			return nil, err
		}

		// If order info is nil, order service is not properly integrated
		if orderInfo == nil {
			return nil, code.ErrInternalServer
		}

		// Validate order is paid (BR-001: Cannot ship unpaid orders)
		if !orderInfo.IsPaid {
			return nil, code.ErrRefundOrderNotPaid
		}

		// Validate order is not already fully shipped
		if orderInfo.FulfillmentStatus == int8(fulfillment.OrderFulfillmentStatusShipped) || orderInfo.FulfillmentStatus == int8(fulfillment.OrderFulfillmentStatusDelivered) {
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

		shipmentNo := fulfillment.GenerateShipmentNo( int(shipmentID))

		// Create shipment entity
		shipment := &fulfillment.Shipment{
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
					Model:       application.Model{ID: itemID, CreatedAt: time.Now().UTC()},
					ShipmentID:  shipmentID,
					OrderItemID: itemReq.OrderItemID,
					ProductID:   itemReq.ProductID,
					SKUID:       itemReq.SKUID,
					ProductName: itemReq.ProductName,
					SKUName:     itemReq.SKUName,
					Image:       itemReq.Image,
					Quantity:    itemReq.Quantity,
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

	// Look up order_no so the response includes it (used by frontend detail page)
	var orderNo string
	if order, err := a.orderRepo.FindByID(ctx, a.db,  result.OrderID); err == nil && order != nil {
		orderNo = order.OrderNo
	}
	return toShipmentResponse(result, carrier, orderNo), nil
}

func (a *orderFulfillmentApp) GetFulfillmentSummary(ctx context.Context) (*FulfillmentSummary, error) {
	// Count shipments by status
	pendingCount, _ := a.shipmentRepo.CountByStatus(ctx, a.db,  fulfillment.ShipmentStatusPending)
	shippedCount, _ := a.shipmentRepo.CountByStatus(ctx, a.db,  fulfillment.ShipmentStatusShipped)
	inTransitCount, _ := a.shipmentRepo.CountByStatus(ctx, a.db,  fulfillment.ShipmentStatusInTransit)
	deliveredCount, _ := a.shipmentRepo.CountByStatus(ctx, a.db,  fulfillment.ShipmentStatusDelivered)

	// Count refunds by status
	pendingRefundCount, _ := a.refundRepo.CountByStatus(ctx, a.db,  fulfillment.RefundStatusPending)
	approvedRefundCount, _ := a.refundRepo.CountByStatus(ctx, a.db,  fulfillment.RefundStatusApproved)

	// Get today's stats
	todayOrders, _ := a.orderRepo.CountTodayOrders(ctx, a.db)
	todayGMV, _ := a.orderRepo.SumTodayGMV(ctx, a.db)

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
func (a *orderFulfillmentApp) UpdateOrderRemark(ctx context.Context,  orderID int64, remark string) error {
	return a.orderRepo.UpdateRemark(ctx, a.db,  orderID, remark)
}

// AdjustOrderPrice adjusts the price of an order
func (a *orderFulfillmentApp) AdjustOrderPrice(ctx context.Context,  userID int64, orderID int64, adjustAmount decimal.Decimal, reason string) (*fulfillment.AdjustPriceResponse, error) {
	// Get order
	order, err := a.orderRepo.FindByID(ctx, a.db,  orderID)
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
		AdjustedAt:     order.AdjustedAt.Format(time.RFC3339),
	}, nil
}

// ExportOrders exports orders to CSV
func (a *orderFulfillmentApp) ExportOrders(ctx context.Context,  req QueryOrderRequest) ([]*fulfillment.OrderExportRow, int64, error) {
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
	orders, err := a.orderRepo.FindForExport(ctx, a.db,  query)
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
			OrderNo:           o.OrderNo,
			Status:            o.Status.Text(),
			FulfillmentStatus: o.FulfillmentStatus.Text(),
			RefundStatus:      o.RefundStatus.Text(),
			TotalAmount:       o.TotalAmount,
			DiscountAmount:    o.DiscountAmount,
			ShippingFee:       o.ShippingFee,
			PayAmount:         o.PayAmount,
			ReceiverName:      o.ReceiverName,
			ReceiverPhone:     o.ReceiverPhone,
			ReceiverAddress:   o.ReceiverAddress,
			PaymentMethod:     o.PaymentMethod,
			CreatedAt:         o.Audit.CreatedAt,
			PaidAt:            o.PaidAt,
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
		ID:           r.Model.ID,
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
		CreatedAt:    r.Model.CreatedAt,
		UpdatedAt:    r.Model.UpdatedAt,
	}
}

// toOrderShippingAddress maps the order's receiver fields into a shipping address DTO.
func toOrderShippingAddress(o *fulfillment.Order) *OrderShippingAddress {
	return &OrderShippingAddress{
		ReceiverName:  o.ReceiverName,
		ReceiverPhone: o.ReceiverPhone,
		Address:       o.ReceiverAddress,
		FullAddress:   o.ReceiverAddress,
	}
}

// toOrderFulfillmentItem maps a single OrderItem into its DTO form.
func toOrderFulfillmentItem(item fulfillment.OrderItem) *OrderFulfillmentItem {
	return &OrderFulfillmentItem{
		OrderItemID: item.ID,
		ProductID:   item.ProductID,
		SKUID:       item.SKUID,
		ProductName: item.ProductName,
		SKUName:     item.SKUName,
		Image:       item.Image,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		Currency:    item.Currency,
	}
}

// toOrderFulfillmentItems maps a slice of OrderItem into DTOs.
// Returns nil when items is empty so JSON encoding produces "items": null
// only when expected; callers that need [] can substitute empty slice.
func toOrderFulfillmentItems(items []fulfillment.OrderItem) []*OrderFulfillmentItem {
	if len(items) == 0 {
		return nil
	}
	result := make([]*OrderFulfillmentItem, len(items))
	for i, item := range items {
		result[i] = toOrderFulfillmentItem(item)
	}
	return result
}

// toOrderFulfillmentDetail builds the base detail (no items, shipments, refund) from an Order.
// Items, shipments, and refund are attached separately so this mapper can be reused
// before those related records are loaded.
func toOrderFulfillmentDetail(o *fulfillment.Order) *OrderFulfillmentDetail {
	return &OrderFulfillmentDetail{
		OrderID:           o.ID,
		OrderNo:           o.OrderNo,
		Status:            string(o.Status),
		FulfillmentStatus: int8(o.FulfillmentStatus),
		FulfillmentText:   o.FulfillmentStatus.Text(),
		RefundStatus:      int8(o.RefundStatus),
		RefundText:        o.RefundStatus.Text(),
		TotalAmount:       o.TotalAmount,
		Currency:          o.Currency,
		UserID:            o.UserID,
		UserName:          o.ReceiverName,  // user_name mirrors receiver_name
		UserPhone:         o.ReceiverPhone, // user_phone mirrors receiver_phone
		ShippingAddress:   toOrderShippingAddress(o),
		Remark:            o.Remark, // 用户备注
		PaidAt:            o.PaidAt,
		ShippedAt:         o.ShippedAt,
		DeliveredAt:       o.DeliveredAt,
		CreatedAt:         o.Audit.CreatedAt,
		UpdatedAt:         o.Audit.UpdatedAt,
	}
}

// buildShipmentResponses converts already-loaded shipments to DTOs using
// the pre-loaded shipment-items map (keyed by shipment ID) and carriers map
// (keyed by carrier code). orderNo is provided by the caller because all
// shipments in a single call share the same parent order — avoids an extra
// query per shipment. Returns nil when shipments is empty.
func buildShipmentResponses(
	shipments []*fulfillment.Shipment,
	shipmentItemsMap map[int64][]fulfillment.ShipmentItem,
	carrierMap map[string]*fulfillment.Carrier,
	orderNo string,
) []*ShipmentResponse {
	if len(shipments) == 0 {
		return nil
	}
	result := make([]*ShipmentResponse, len(shipments))
	for i, s := range shipments {
		s.Items = shipmentItemsMap[s.Model.ID]
		result[i] = toShipmentResponse(s, carrierMap[s.CarrierCode], orderNo)
	}
	return result
}

// pickLatestRefund returns the most recent refund's DTO, or nil if none exist.
// "Latest" is defined by the largest CreatedAt timestamp.
func pickLatestRefund(refunds []*fulfillment.Refund) *RefundResponse {
	if len(refunds) == 0 {
		return nil
	}
	latest := refunds[0]
	for _, r := range refunds[1:] {
		if r.Model.CreatedAt.After(latest.Model.CreatedAt) {
			latest = r
		}
	}
	return toRefundResponse(latest)
}

// buildOrderQuery converts a QueryOrderRequest into the repository's OrderQuery.
func buildOrderQuery(req QueryOrderRequest) fulfillment.OrderQuery {
	return fulfillment.OrderQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		OrderNo:           req.OrderNo,
		UserID:            req.UserID,
		UserName:          req.UserName,
		Status:            fulfillment.OrderStatus(req.Status),
		FulfillmentStatus: fulfillment.OrderFulfillmentStatus(req.FulfillmentStatus),
		RefundStatus:      fulfillment.OrderRefundStatus(req.RefundStatus),
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
	}
}

// collectOrderIDs extracts the IDs from a slice of orders in order.
func collectOrderIDs(orders []*fulfillment.Order) []int64 {
	ids := make([]int64, len(orders))
	for i, o := range orders {
		ids[i] = o.ID
	}
	return ids
}
