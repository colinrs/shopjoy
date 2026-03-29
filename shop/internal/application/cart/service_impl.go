package cart

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/shop/internal/domain/cart"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	db       *gorm.DB
	cartRepo cart.Repository
}

func NewService(db *gorm.DB, cartRepo cart.Repository) Service {
	return &ServiceImpl{
		db:       db,
		cartRepo: cartRepo,
	}
}

func (s *ServiceImpl) GetCart(ctx context.Context, tenantID shared.TenantID, userID int64) (*CartResponse, error) {
	c, err := s.cartRepo.FindByUserID(ctx, s.db, tenantID, userID)
	if err != nil {
		return nil, err
	}

	return toCartResponse(c), nil
}

func (s *ServiceImpl) AddItem(ctx context.Context, req AddToCartRequest) error {
	var c *cart.Cart
	var err error

	if req.UserID > 0 {
		c, err = s.cartRepo.FindByUserID(ctx, s.db, req.TenantID, req.UserID)
	} else {
		c, err = s.cartRepo.FindBySessionID(ctx, s.db, req.TenantID, req.SessionID)
	}

	if err != nil {
		return err
	}

	if c.ID == 0 {
		c.TenantID = req.TenantID
		c.UserID = req.UserID
		c.SessionID = req.SessionID
		if err := s.cartRepo.Save(ctx, s.db, c); err != nil {
			return err
		}
	}

	priceAmount, err := shared.ParseMoneyFromString(req.Price)
	if err != nil {
		return err
	}
	price := shared.NewMoney(priceAmount, req.Currency)
	return c.AddItem(req.ProductID, req.SKUId, req.Quantity, price, req.ProductName, req.SKUName, req.Image)
}

func (s *ServiceImpl) UpdateItem(ctx context.Context, req UpdateCartItemRequest) error {
	c, err := s.cartRepo.FindByUserID(ctx, s.db, req.TenantID, req.UserID)
	if err != nil {
		return err
	}

	if err := c.UpdateItem(req.ItemID, req.Quantity); err != nil {
		return err
	}

	for _, item := range c.Items {
		if item.ID == req.ItemID {
			return s.cartRepo.SaveItem(ctx, s.db, &item)
		}
	}

	return nil
}

func (s *ServiceImpl) RemoveItem(ctx context.Context, req RemoveFromCartRequest) error {
	return s.cartRepo.DeleteItem(ctx, s.db, req.TenantID, req.ItemID)
}

func (s *ServiceImpl) ClearCart(ctx context.Context, tenantID shared.TenantID, userID int64) error {
	return s.cartRepo.Clear(ctx, s.db, tenantID, userID)
}

func (s *ServiceImpl) MergeCart(ctx context.Context, tenantID shared.TenantID, sessionID string, userID int64) error {
	return s.cartRepo.Merge(ctx, s.db, tenantID, sessionID, userID)
}

func toCartResponse(c *cart.Cart) *CartResponse {
	resp := &CartResponse{
		Items:     make([]CartItemResponse, len(c.Items)),
		ItemCount: len(c.Items),
	}

	total := shared.NewMoney(decimal.Zero, "CNY")
	for i, item := range c.Items {
		resp.Items[i] = CartItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			SKUId:       item.SKUId,
			ProductName: item.ProductName,
			SKUName:     item.SKUName,
			Image:       item.Image,
			Price:       shared.FormatMoneyToStringOnly(item.Price.Amount),
			Quantity:    item.Quantity,
			TotalAmount: shared.FormatMoneyToStringOnly(item.TotalAmount.Amount),
			Selected:    item.Selected,
		}
		total, _ = total.Add(item.TotalAmount)
	}

	resp.TotalAmount = shared.FormatMoneyToStringOnly(total.Amount)
	return resp
}
