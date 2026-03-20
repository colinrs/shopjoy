package order

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/colinrs/shopjoy/shop/internal/domain/cart"
	"github.com/colinrs/shopjoy/shop/internal/domain/order"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	db        *gorm.DB
	orderRepo order.Repository
	cartRepo  cart.Repository
	idGen     snowflake.Snowflake
}

func NewService(db *gorm.DB, orderRepo order.Repository, cartRepo cart.Repository, idGen snowflake.Snowflake) Service {
	return &ServiceImpl{
		db:        db,
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
		idGen:     idGen,
	}
}

func (s *ServiceImpl) CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderResponse, error) {
	cart, err := s.cartRepo.FindByUserID(ctx, s.db, req.TenantID, req.UserID)
	if err != nil {
		return nil, err
	}

	if cart.IsEmpty() {
		return nil, code.ErrOrderCartEmpty
	}

	orderID := order.GenerateOrderNo(req.TenantID)

	o := &order.Order{
		ID:       orderID,
		TenantID: req.TenantID,
		UserID:   req.UserID,
		OrderNo:  orderID,
		Status:   order.StatusPendingPayment,
		Currency: "CNY",
		ExpireAt: time.Now().Add(30 * time.Minute),
		Address: order.ShippingAddress{
			Name:     req.Address.Name,
			Phone:    req.Address.Phone,
			Province: req.Address.Province,
			City:     req.Address.City,
			District: req.Address.District,
			Address:  req.Address.Address,
			ZipCode:  req.Address.ZipCode,
		},
		Remark: req.Remark,
		Audit:  shared.NewAuditInfo(req.UserID),
	}

	for _, item := range cart.Items {
		if !item.Selected {
			continue
		}
		o.Items = append(o.Items, order.OrderItem{
			ProductID:   item.ProductID,
			SKUId:       item.SKUId,
			ProductName: item.ProductName,
			SKUName:     item.SKUName,
			Image:       item.Image,
			Price:       item.Price,
			Quantity:    item.Quantity,
			TotalAmount: item.TotalAmount,
		})
	}

	o.CalculateTotals()

	if err := s.orderRepo.Create(ctx, s.db, o); err != nil {
		return nil, err
	}

	s.cartRepo.Clear(ctx, s.db, req.TenantID, req.UserID)

	return toOrderResponse(o), nil
}

func (s *ServiceImpl) GetOrder(ctx context.Context, tenantID shared.TenantID, orderID string) (*OrderResponse, error) {
	o, err := s.orderRepo.FindByID(ctx, s.db, tenantID, orderID)
	if err != nil {
		return nil, err
	}
	return toOrderResponse(o), nil
}

func (s *ServiceImpl) GetOrderByNo(ctx context.Context, tenantID shared.TenantID, orderNo string) (*OrderResponse, error) {
	o, err := s.orderRepo.FindByOrderNo(ctx, s.db, tenantID, orderNo)
	if err != nil {
		return nil, err
	}
	return toOrderResponse(o), nil
}

func (s *ServiceImpl) GetUserOrders(ctx context.Context, tenantID shared.TenantID, userID int64, query QueryRequest) (*OrderListResponse, error) {
	oQuery := order.Query{
		PageQuery: query.PageQuery,
		Status:    query.Status,
	}
	oQuery.PageQuery.Validate()

	orders, total, err := s.orderRepo.FindByUserID(ctx, s.db, tenantID, userID, oQuery)
	if err != nil {
		return nil, err
	}

	resp := &OrderListResponse{
		List:     make([]*OrderResponse, len(orders)),
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	for i, o := range orders {
		resp.List[i] = toOrderResponse(o)
	}

	return resp, nil
}

func (s *ServiceImpl) CancelOrder(ctx context.Context, tenantID shared.TenantID, userID int64, orderID string, reason string) error {
	o, err := s.orderRepo.FindByID(ctx, s.db, tenantID, orderID)
	if err != nil {
		return err
	}

	if o.UserID != userID {
		return code.ErrForbidden
	}

	return o.Cancel(reason)
}

func (s *ServiceImpl) PayOrder(ctx context.Context, tenantID shared.TenantID, orderID string, paymentID string) error {
	o, err := s.orderRepo.FindByID(ctx, s.db, tenantID, orderID)
	if err != nil {
		return err
	}

	if err := o.Pay(paymentID); err != nil {
		return err
	}

	return s.orderRepo.Update(ctx, s.db, o)
}

func toOrderResponse(o *order.Order) *OrderResponse {
	resp := &OrderResponse{
		ID:             o.ID,
		OrderNo:        o.OrderNo,
		UserID:         o.UserID,
		Status:         int(o.Status),
		TotalAmount:    o.TotalAmount.Amount,
		DiscountAmount: o.DiscountAmount.Amount,
		FreightAmount:  o.FreightAmount.Amount,
		PayAmount:      o.PayAmount.Amount,
		Currency:       o.Currency,
		Remark:         o.Remark,
		ExpireAt:       o.ExpireAt,
		CreatedAt:      o.Audit.CreatedAt,
		Address: AddressResponse{
			Name:     o.Address.Name,
			Phone:    o.Address.Phone,
			Province: o.Address.Province,
			City:     o.Address.City,
			District: o.Address.District,
			Address:  o.Address.Address,
			ZipCode:  o.Address.ZipCode,
		},
		Items: make([]OrderItemResponse, len(o.Items)),
	}

	if o.PaidAt != nil {
		resp.PaidAt = o.PaidAt
	}
	if o.ShippedAt != nil {
		resp.ShippedAt = o.ShippedAt
	}

	for i, item := range o.Items {
		resp.Items[i] = OrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			SKUId:       item.SKUId,
			SKUName:     item.SKUName,
			Image:       item.Image,
			Price:       item.Price.Amount,
			Quantity:    item.Quantity,
			TotalAmount: item.TotalAmount.Amount,
		}
	}

	return resp
}