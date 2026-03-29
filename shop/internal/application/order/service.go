package order

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/shop/internal/domain/order"
)

// DTOs
type CreateOrderRequest struct {
	TenantID shared.TenantID
	UserID   int64
	Items    []OrderItemRequest
	Address  AddressRequest
	Remark   string
	CouponID *int64
}

type OrderItemRequest struct {
	ProductID int64
	SKUId     int64
	Quantity  int
}

type AddressRequest struct {
	Name     string
	Phone    string
	Province string
	City     string
	District string
	Address  string
	ZipCode  string
}

type OrderResponse struct {
	ID             string
	OrderNo        string
	UserID         int64
	Status         int
	TotalAmount    string // API monetary values use string (yuan)
	DiscountAmount string
	FreightAmount  string
	PayAmount      string
	Currency       string
	Items          []OrderItemResponse
	Address        AddressResponse
	Remark         string
	ExpireAt       time.Time
	PaidAt         *time.Time
	ShippedAt      *time.Time
	CreatedAt      time.Time
}

type OrderItemResponse struct {
	ProductID   int64
	ProductName string
	SKUId       int64
	SKUName     string
	Image       string
	Price       string // API monetary values use string (yuan)
	Quantity    int
	TotalAmount string
}

type AddressResponse struct {
	Name     string
	Phone    string
	Province string
	City     string
	District string
	Address  string
	ZipCode  string
}

type OrderListResponse struct {
	List     []*OrderResponse
	Total    int64
	Page     int
	PageSize int
}

// Service 订单应用服务接口
type Service interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderResponse, error)
	GetOrder(ctx context.Context, tenantID shared.TenantID, orderID string) (*OrderResponse, error)
	GetOrderByNo(ctx context.Context, tenantID shared.TenantID, orderNo string) (*OrderResponse, error)
	GetUserOrders(ctx context.Context, tenantID shared.TenantID, userID int64, query QueryRequest) (*OrderListResponse, error)
	CancelOrder(ctx context.Context, tenantID shared.TenantID, userID int64, orderID string, reason string) error
	PayOrder(ctx context.Context, tenantID shared.TenantID, orderID string, paymentID string) error
}

type QueryRequest struct {
	shared.PageQuery
	Status    order.Status
	StartTime *time.Time
	EndTime   *time.Time
}
