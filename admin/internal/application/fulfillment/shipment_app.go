package fulfillment

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// CreateShipmentRequest 创建发货单请求
type CreateShipmentRequest struct {
	OrderID      int64
	CarrierCode  string
	CarrierName  string
	TrackingNo   string
	ShippingCost decimal.Decimal
	Currency     string
	Weight       decimal.Decimal
	Remark       string
	Items        []CreateShipmentItemRequest
}

// CreateShipmentItemRequest 创建发货单明细请求
type CreateShipmentItemRequest struct {
	OrderItemID int64
	ProductID   int64
	SKUID       int64
	ProductName string
	SKUName     string
	Image       string
	Quantity    int
}

// UpdateShipmentRequest 更新发货单请求
type UpdateShipmentRequest struct {
	ID           int64
	CarrierCode  string
	CarrierName  string
	TrackingNo   string
	ShippingCost decimal.Decimal
	Currency     string
	Weight       decimal.Decimal
	Remark       string
}

// QueryShipmentRequest 查询发货单请求
type QueryShipmentRequest struct {
	Page       int
	PageSize   int
	ShipmentNo string
	OrderID    int64
	TrackingNo string
	// Status 为 nil 表示不过滤状态；指向 ShipmentStatusXxx 表示按该状态过滤。
	Status            *fulfillment.ShipmentStatus
	CarrierCode       string
	FulfillmentStatus int8
	StartTime         time.Time
	EndTime           time.Time
}

// ShipmentResponse 发货单响应
type ShipmentResponse struct {
	ID           int64                   `json:"id"`
	ShipmentNo   string                  `json:"shipment_no"`
	OrderID      int64                   `json:"order_id"`
	OrderNo      string                  `json:"order_no"`
	Status       int                     `json:"status"`
	StatusText   string                  `json:"status_text"`
	Carrier      string                  `json:"carrier"`
	CarrierCode  string                  `json:"carrier_code"`
	TrackingNo   string                  `json:"tracking_no"`
	TrackingURL  string                  `json:"tracking_url,omitempty"`
	ShippingCost string                  `json:"shipping_cost"`
	Currency     string                  `json:"currency"`
	Weight       string                  `json:"weight"`
	ShippedAt    *time.Time              `json:"shipped_at,omitempty"`
	DeliveredAt  *time.Time              `json:"delivered_at,omitempty"`
	Remark       string                  `json:"remark"`
	Items        []*ShipmentItemResponse `json:"items"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	CreatedBy    int64                   `json:"created_by"`
}

// ShipmentItemResponse 发货单明细响应
type ShipmentItemResponse struct {
	ID          int64  `json:"id"`
	ShipmentID  int64  `json:"shipment_id"`
	OrderItemID int64  `json:"order_item_id"`
	ProductID   int64  `json:"product_id"`
	SKUID       int64  `json:"sku_id"`
	ProductName string `json:"product_name"`
	SKUName     string `json:"sku_name"`
	Image       string `json:"image"`
	Quantity    int    `json:"quantity"`
}

// ShipmentListResponse 发货单列表响应
type ShipmentListResponse struct {
	List     []*ShipmentResponse `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// CarrierResponse 物流公司响应
type CarrierResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	TrackingURL string `json:"tracking_url"`
	IsActive    bool   `json:"is_active"`
	Sort        int    `json:"sort"`
}

// ShipmentApp 发货单应用服务接口
type ShipmentApp interface {
	CreateShipment(ctx context.Context,  userID int64, req CreateShipmentRequest) (*ShipmentResponse, error)
	BatchCreateShipments(ctx context.Context,  userID int64, carrierCode, carrierName string, shipments []BatchShipmentItem) (*BatchShipmentResult, error)
	GetShipment(ctx context.Context,  id int64) (*ShipmentResponse, error)
	ListShipments(ctx context.Context,  req QueryShipmentRequest) (*ShipmentListResponse, error)
	UpdateShipment(ctx context.Context,  userID int64, req UpdateShipmentRequest) (*ShipmentResponse, error)
	UpdateShipmentStatus(ctx context.Context,  userID int64, id int64, status fulfillment.ShipmentStatus) (*ShipmentResponse, error)
	GetOrderShipments(ctx context.Context,  orderID int64) ([]*ShipmentResponse, error)
	CancelShipment(ctx context.Context,  userID int64, id int64, reason string) (*ShipmentResponse, error)
}

// BatchShipmentItem 批量发货单项
type BatchShipmentItem struct {
	OrderID    int64
	TrackingNo string
}

// BatchShipmentResult 批量发货结果
type BatchShipmentResult struct {
	Total   int
	Success int
	Failed  int
	Results []BatchShipmentResultItem
}

// BatchShipmentResultItem 批量发货结果项
type BatchShipmentResultItem struct {
	OrderID    int64
	ShipmentID int64
	ShipmentNo string
	Success    bool
	Error      string
}

type shipmentApp struct {
	db               *gorm.DB
	shipmentRepo     fulfillment.ShipmentRepository
	shipmentItemRepo fulfillment.ShipmentItemRepository
	carrierRepo      fulfillment.CarrierRepository
	orderRepo        fulfillment.OrderRepository
	idGen            snowflake.Snowflake
}

// NewShipmentApp 创建发货单应用服务
func NewShipmentApp(
	db *gorm.DB,
	shipmentRepo fulfillment.ShipmentRepository,
	shipmentItemRepo fulfillment.ShipmentItemRepository,
	carrierRepo fulfillment.CarrierRepository,
	orderRepo fulfillment.OrderRepository,
	idGen snowflake.Snowflake,
) ShipmentApp {
	return &shipmentApp{
		db:               db,
		shipmentRepo:     shipmentRepo,
		shipmentItemRepo: shipmentItemRepo,
		carrierRepo:      carrierRepo,
		orderRepo:        orderRepo,
		idGen:            idGen,
	}
}

func (a *shipmentApp) CreateShipment(ctx context.Context, userID int64, req CreateShipmentRequest) (*ShipmentResponse, error) {
	// Validate carrier
	carrier, err := a.carrierRepo.FindByCode(ctx, a.db, req.CarrierCode)
	if err != nil {
		if err == code.ErrCarrierNotFound {
			// Allow custom carrier
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

		shipmentNo := fulfillment.GenerateShipmentNo( int(shipmentID%1000000))

		// Create shipment entity
		shipment := &fulfillment.Shipment{
			OrderID:          req.OrderID,
			ShipmentNo:       shipmentNo,
			Status:           fulfillment.ShipmentStatusShipped, // Automatically set to shipped
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
					Model: application.Model{
						ID:        itemID,
						CreatedAt: time.Now().UTC(),
					},
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

	orderNo := a.lookupOrderNo(ctx,  result.OrderID)
	return toShipmentResponse(result, carrier, orderNo), nil
}

func (a *shipmentApp) BatchCreateShipments(ctx context.Context, userID int64, carrierCode, carrierName string, shipments []BatchShipmentItem) (*BatchShipmentResult, error) {
	// Validate carrier
	_, err := a.carrierRepo.FindByCode(ctx, a.db, carrierCode)
	if err != nil {
		if err == code.ErrCarrierNotFound {
			if carrierName == "" {
				return nil, code.ErrShipmentCarrierRequired
			}
		} else {
			return nil, err
		}
	}

	result := &BatchShipmentResult{
		Total:   len(shipments),
		Results: make([]BatchShipmentResultItem, len(shipments)),
	}

	for i, item := range shipments {
		shipmentResp, err := a.CreateShipment(ctx,  userID, CreateShipmentRequest{
			OrderID:     item.OrderID,
			CarrierCode: carrierCode,
			CarrierName: carrierName,
			TrackingNo:  item.TrackingNo,
		})

		result.Results[i] = BatchShipmentResultItem{
			OrderID: item.OrderID,
		}

		if err != nil {
			result.Results[i].Success = false
			result.Results[i].Error = err.Error()
			result.Failed++
		} else {
			result.Results[i].ShipmentID = shipmentResp.ID
			result.Results[i].ShipmentNo = shipmentResp.ShipmentNo
			result.Results[i].Success = true
			result.Success++
		}
	}

	return result, nil
}

func (a *shipmentApp) GetShipment(ctx context.Context,  id int64) (*ShipmentResponse, error) {
	shipment, err := a.shipmentRepo.FindByID(ctx, a.db,  id)
	if err != nil {
		return nil, err
	}

	// Load items
	items, err := a.shipmentItemRepo.FindByShipmentID(ctx, a.db,  shipment.ID)
	if err != nil {
		return nil, err
	}
	shipment.Items = items

	// Get carrier for tracking URL
	carrier, err := a.carrierRepo.FindByCode(ctx, a.db, shipment.CarrierCode)
	if err != nil {
		log.Printf("GetShipment: find carrier by code %s error: %v", shipment.CarrierCode, err)
	}

	orderNo := a.lookupOrderNo(ctx,  shipment.OrderID)
	return toShipmentResponse(shipment, carrier, orderNo), nil
}

func (a *shipmentApp) ListShipments(ctx context.Context,  req QueryShipmentRequest) (*ShipmentListResponse, error) {
	query := fulfillment.ShipmentQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		OrderID:     req.OrderID,
		Status:      req.Status,
		CarrierCode: req.CarrierCode,
		TrackingNo:  req.TrackingNo,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	shipments, total, err := a.shipmentRepo.FindList(ctx, a.db,  query)
	if err != nil {
		return nil, err
	}

	if len(shipments) == 0 {
		return &ShipmentListResponse{
			List:     []*ShipmentResponse{},
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, nil
	}

	// Batch load items for all shipments (fix N+1 query)
	shipmentIDs := make([]int64, len(shipments))
	for i, s := range shipments {
		shipmentIDs[i] = s.ID
	}
	itemsMap, err := a.shipmentItemRepo.FindByShipmentIDs(ctx, a.db,  shipmentIDs)
	if err != nil {
		log.Printf("ListShipments: batch find items error: %v", err)
		itemsMap = make(map[int64][]fulfillment.ShipmentItem)
	}

	// Collect carrier codes for batch query
	carrierCodes := make([]string, 0, len(shipments))
	seenCodes := make(map[string]bool)
	for _, s := range shipments {
		if s.CarrierCode != "" && !seenCodes[s.CarrierCode] {
			carrierCodes = append(carrierCodes, s.CarrierCode)
			seenCodes[s.CarrierCode] = true
		}
	}

	// Batch load carriers (fix N+1 query)
	carriersMap, err := a.carrierRepo.FindByCodes(ctx, a.db, carrierCodes)
	if err != nil {
		log.Printf("ListShipments: batch find carriers error: %v", err)
		carriersMap = make(map[string]*fulfillment.Carrier)
	}

	// Batch load orders to enrich shipments with order_no (avoid N+1)
	orderIDs := make([]int64, 0, len(shipments))
	seenOrderIDs := make(map[int64]bool)
	for _, s := range shipments {
		if s.OrderID != 0 && !seenOrderIDs[s.OrderID] {
			orderIDs = append(orderIDs, s.OrderID)
			seenOrderIDs[s.OrderID] = true
		}
	}
	orderNoMap := a.lookupOrderNos(ctx,  orderIDs)

	list := make([]*ShipmentResponse, len(shipments))
	for i, s := range shipments {
		s.Items = itemsMap[s.ID]
		carrier := carriersMap[s.CarrierCode]
		list[i] = toShipmentResponse(s, carrier, orderNoMap[s.OrderID])
	}

	return &ShipmentListResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (a *shipmentApp) UpdateShipment(ctx context.Context,  userID int64, req UpdateShipmentRequest) (*ShipmentResponse, error) {
	shipment, err := a.shipmentRepo.FindByID(ctx, a.db,  req.ID)
	if err != nil {
		return nil, err
	}

	// Update carrier info
	if req.CarrierCode != "" {
		carrier, err := a.carrierRepo.FindByCode(ctx, a.db, req.CarrierCode)
		if err != nil {
			if err == code.ErrCarrierNotFound && req.CarrierName != "" {
				shipment.CarrierCode = req.CarrierCode
				shipment.Carrier = req.CarrierName
			} else if err != code.ErrCarrierNotFound {
				return nil, err
			}
		} else {
			shipment.CarrierCode = carrier.Code
			shipment.Carrier = carrier.Name
		}
	}

	if req.TrackingNo != "" {
		shipment.TrackingNo = req.TrackingNo
	}

	shipment.ShippingCost = req.ShippingCost
	if req.Currency != "" {
		shipment.ShippingCurrency = req.Currency
	}
	shipment.Weight = req.Weight
	shipment.Remark = req.Remark
	shipment.UpdatedAt = time.Now().UTC()

	if err := a.shipmentRepo.Update(ctx, a.db, shipment); err != nil {
		return nil, err
	}

	// Load items
	items, err := a.shipmentItemRepo.FindByShipmentID(ctx, a.db,  shipment.ID)
	if err != nil {
		log.Printf("UpdateShipment: find items by shipment ID %d error: %v", shipment.ID, err)
	}
	shipment.Items = items

	carrier, err := a.carrierRepo.FindByCode(ctx, a.db, shipment.CarrierCode)
	if err != nil {
		log.Printf("UpdateShipment: find carrier by code %s error: %v", shipment.CarrierCode, err)
	}
	orderNo := a.lookupOrderNo(ctx,  shipment.OrderID)
	return toShipmentResponse(shipment, carrier, orderNo), nil
}

func (a *shipmentApp) UpdateShipmentStatus(ctx context.Context,  userID int64, id int64, status fulfillment.ShipmentStatus) (*ShipmentResponse, error) {
	shipment, err := a.shipmentRepo.FindByID(ctx, a.db,  id)
	if err != nil {
		return nil, err
	}

	// Update status based on the new status
	switch status {
	case fulfillment.ShipmentStatusShipped:
		if err := shipment.Ship(shipment.Carrier, shipment.CarrierCode, shipment.TrackingNo, userID); err != nil {
			return nil, err
		}
	case fulfillment.ShipmentStatusInTransit:
		if err := shipment.MarkInTransit(userID); err != nil {
			return nil, err
		}
	case fulfillment.ShipmentStatusDelivered:
		if err := shipment.MarkDelivered(userID); err != nil {
			return nil, err
		}
	case fulfillment.ShipmentStatusFailed:
		if err := shipment.MarkFailed("", userID); err != nil {
			return nil, err
		}
	case fulfillment.ShipmentStatusCancelled:
		if err := shipment.CancelShipment("", userID); err != nil {
			return nil, err
		}
	default:
		return nil, code.ErrShipmentInvalidStatusTransition
	}

	if err := a.shipmentRepo.Update(ctx, a.db, shipment); err != nil {
		return nil, err
	}

	// Load items
	items, err := a.shipmentItemRepo.FindByShipmentID(ctx, a.db,  shipment.ID)
	if err != nil {
		log.Printf("UpdateShipmentStatus: find items by shipment ID %d error: %v", shipment.ID, err)
	}
	shipment.Items = items

	carrier, err := a.carrierRepo.FindByCode(ctx, a.db, shipment.CarrierCode)
	if err != nil {
		log.Printf("UpdateShipmentStatus: find carrier by code %s error: %v", shipment.CarrierCode, err)
	}
	orderNo := a.lookupOrderNo(ctx,  shipment.OrderID)
	return toShipmentResponse(shipment, carrier, orderNo), nil
}

func (a *shipmentApp) GetOrderShipments(ctx context.Context,  orderID int64) ([]*ShipmentResponse, error) {
	shipments, err := a.shipmentRepo.FindByOrderID(ctx, a.db,  orderID)
	if err != nil {
		return nil, err
	}

	if len(shipments) == 0 {
		return []*ShipmentResponse{}, nil
	}

	// Batch load items for all shipments (fix N+1 query)
	shipmentIDs := make([]int64, len(shipments))
	for i, s := range shipments {
		shipmentIDs[i] = s.ID
	}
	itemsMap, err := a.shipmentItemRepo.FindByShipmentIDs(ctx, a.db,  shipmentIDs)
	if err != nil {
		log.Printf("GetOrderShipments: batch find items error: %v", err)
		itemsMap = make(map[int64][]fulfillment.ShipmentItem)
	}

	// Collect carrier codes for batch query
	carrierCodes := make([]string, 0, len(shipments))
	seenCodes := make(map[string]bool)
	for _, s := range shipments {
		if s.CarrierCode != "" && !seenCodes[s.CarrierCode] {
			carrierCodes = append(carrierCodes, s.CarrierCode)
			seenCodes[s.CarrierCode] = true
		}
	}

	// Batch load carriers (fix N+1 query)
	carriersMap, err := a.carrierRepo.FindByCodes(ctx, a.db, carrierCodes)
	if err != nil {
		log.Printf("GetOrderShipments: batch find carriers error: %v", err)
		carriersMap = make(map[string]*fulfillment.Carrier)
	}

	// All shipments share the same orderID — one lookup is enough.
	orderNo := a.lookupOrderNo(ctx,  orderID)

	list := make([]*ShipmentResponse, len(shipments))
	for i, s := range shipments {
		s.Items = itemsMap[s.ID]
		carrier := carriersMap[s.CarrierCode]
		list[i] = toShipmentResponse(s, carrier, orderNo)
	}

	return list, nil
}

func (a *shipmentApp) CancelShipment(ctx context.Context,  userID int64, id int64, reason string) (*ShipmentResponse, error) {
	// Get the shipment
	shipment, err := a.shipmentRepo.FindByID(ctx, a.db,  id)
	if err != nil {
		return nil, err
	}

	// Call entity method to cancel
	if err := shipment.CancelShipment(reason, userID); err != nil {
		return nil, err
	}

	// Update the shipment
	if err := a.shipmentRepo.Update(ctx, a.db, shipment); err != nil {
		return nil, err
	}

	// Load items
	items, err := a.shipmentItemRepo.FindByShipmentID(ctx, a.db,  shipment.ID)
	if err != nil {
		log.Printf("CancelShipment: find items by shipment ID %d error: %v", shipment.ID, err)
	}
	shipment.Items = items

	// Get carrier for tracking URL
	carrier, err := a.carrierRepo.FindByCode(ctx, a.db, shipment.CarrierCode)
	if err != nil {
		log.Printf("CancelShipment: find carrier by code %s error: %v", shipment.CarrierCode, err)
	}

	orderNo := a.lookupOrderNo(ctx,  shipment.OrderID)
	return toShipmentResponse(shipment, carrier, orderNo), nil
}

// toShipmentResponse 转换为响应DTO.
// orderNo comes from the related Order entity (loaded by callers to avoid N+1).
// Pass an empty string when no order is associated (e.g., orphaned records).
func toShipmentResponse(s *fulfillment.Shipment, carrier *fulfillment.Carrier, orderNo string) *ShipmentResponse {
	items := make([]*ShipmentItemResponse, len(s.Items))
	for i, item := range s.Items {
		items[i] = &ShipmentItemResponse{
			ID:          item.Model.ID,
			ShipmentID:  item.ShipmentID,
			OrderItemID: item.OrderItemID,
			ProductID:   item.ProductID,
			SKUID:       item.SKUID,
			ProductName: item.ProductName,
			SKUName:     item.SKUName,
			Image:       item.Image,
			Quantity:    item.Quantity,
		}
	}

	var trackingURL string
	if carrier != nil && carrier.TrackingURL != "" {
		trackingURL = carrier.GetTrackingURL(s.TrackingNo)
	}

	return &ShipmentResponse{
		ID:           s.ID,
		ShipmentNo:   s.ShipmentNo,
		OrderID:      s.OrderID,
		OrderNo:      orderNo,
		Status:       int(s.Status),
		StatusText:   s.Status.String(),
		Carrier:      s.Carrier,
		CarrierCode:  s.CarrierCode,
		TrackingNo:   s.TrackingNo,
		TrackingURL:  trackingURL,
		ShippingCost: s.ShippingCost.StringFixed(2),
		Currency:     s.ShippingCurrency,
		Weight:       s.Weight.StringFixed(3),
		ShippedAt:    s.ShippedAt,
		DeliveredAt:  s.DeliveredAt,
		Remark:       s.Remark,
		Items:        items,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

// lookupOrderNo returns the order_no for a single order, or "" if the order
// is missing. Used at the read boundary to enrich shipment responses.
func (a *shipmentApp) lookupOrderNo(ctx context.Context, orderID int64) string {
	if orderID == 0 {
		return ""
	}
	order, err := a.orderRepo.FindByID(ctx, a.db,  orderID)
	if err != nil || order == nil {
		return ""
	}
	return order.OrderNo
}

// lookupOrderNos batch-loads orders to avoid N+1 queries on list endpoints.
// Returns a map keyed by order ID; missing orders are omitted.
func (a *shipmentApp) lookupOrderNos(ctx context.Context, orderIDs []int64) map[int64]string {
	result := make(map[int64]string, len(orderIDs))
	if len(orderIDs) == 0 {
		return result
	}
	for _, id := range orderIDs {
		if id == 0 {
			continue
		}
		order, err := a.orderRepo.FindByID(ctx, a.db,  id)
		if err != nil || order == nil {
			continue
		}
		result[id] = order.OrderNo
	}
	return result
}

// FormatMoneyToInt64 parses a money string to int64 (cents)
func FormatMoneyToInt64(s string) int64 {
	if s == "" {
		return 0
	}
	// Try parsing as integer (cents)
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// FormatInt64ToMoney converts int64 to money string
func FormatInt64ToMoney(v int64) string {
	return fmt.Sprintf("%d", v)
}
