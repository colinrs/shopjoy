package cart

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

type AddToCartRequest struct {
	TenantID    shared.TenantID
	UserID      int64
	SessionID   string
	ProductID   int64
	SKUId       int64
	Quantity    int
	Price       string
	Currency    string
	ProductName string
	SKUName     string
	Image       string
}

type UpdateCartItemRequest struct {
	TenantID shared.TenantID
	UserID   int64
	ItemID   int64
	Quantity int
}

type RemoveFromCartRequest struct {
	TenantID shared.TenantID
	UserID   int64
	ItemID   int64
}

type CartItemResponse struct {
	ID          int64
	ProductID   int64
	SKUId       int64
	ProductName string
	SKUName     string
	Image       string
	Price       string
	Quantity    int
	TotalAmount string
	Selected    bool
}

type CartResponse struct {
	Items       []CartItemResponse
	TotalAmount string
	ItemCount   int
}

type Service interface {
	GetCart(ctx context.Context, tenantID shared.TenantID, userID int64) (*CartResponse, error)
	AddItem(ctx context.Context, req AddToCartRequest) error
	UpdateItem(ctx context.Context, req UpdateCartItemRequest) error
	RemoveItem(ctx context.Context, req RemoveFromCartRequest) error
	ClearCart(ctx context.Context, tenantID shared.TenantID, userID int64) error
	MergeCart(ctx context.Context, tenantID shared.TenantID, sessionID string, userID int64) error
}
