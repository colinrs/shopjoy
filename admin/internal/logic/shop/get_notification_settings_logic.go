package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetNotificationSettingsLogic {
	return GetNotificationSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotificationSettingsLogic) GetNotificationSettings() (resp *types.NotificationSettings, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok || tenantID == 0 {
		return nil, code.ErrTenantInvalidID
	}

	// Find shop settings first to get shop ID
	settings, err := l.svcCtx.ShopSettingsRepo.FindByTenantID(l.ctx, l.svcCtx.DB, tenantID)
	if err != nil {
		l.Logger.Errorf("find shop settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if settings == nil {
		return nil, code.ErrShopNotFound
	}

	notificationSettings, err := l.svcCtx.NotificationSettingsRepo.FindByShopID(l.ctx, l.svcCtx.DB, settings.ID)
	if err != nil {
		l.Logger.Errorf("get notification settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if notificationSettings == nil {
		// Return default settings if not configured
		return &types.NotificationSettings{
			OrderCreated:      true,
			OrderPaid:         true,
			OrderShipped:      true,
			OrderCancelled:    true,
			LowStockAlert:     true,
			LowStockThreshold: 10,
			RefundRequested:   true,
			NewReview:         true,
		}, nil
	}

	return toNotificationSettingsResponse(notificationSettings), nil
}

func toNotificationSettingsResponse(n *shop.NotificationSettings) *types.NotificationSettings {
	return &types.NotificationSettings{
		OrderCreated:      n.OrderCreated,
		OrderPaid:         n.OrderPaid,
		OrderShipped:      n.OrderShipped,
		OrderCancelled:    n.OrderCancelled,
		LowStockAlert:     n.LowStockAlert,
		LowStockThreshold: n.LowStockThreshold,
		RefundRequested:   n.RefundRequested,
		NewReview:         n.NewReview,
	}
}
