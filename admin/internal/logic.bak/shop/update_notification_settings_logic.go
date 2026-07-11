package shop

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNotificationSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateNotificationSettingsLogic {
	return UpdateNotificationSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNotificationSettingsLogic) UpdateNotificationSettings(req *types.UpdateNotificationSettingsRequest) (resp *types.NotificationSettings, err error) {
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

	// Find existing notification settings or create new
	notificationSettings, err := l.svcCtx.NotificationSettingsRepo.FindByShopID(l.ctx, l.svcCtx.DB, settings.ID)
	if err != nil {
		l.Logger.Errorf("find notification settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if notificationSettings == nil {
		notificationSettings = &shop.NotificationSettings{
			ShopID: settings.ID,
		}
	}

	// Update fields from request
	notificationSettings.OrderCreated = req.OrderCreated
	notificationSettings.OrderPaid = req.OrderPaid
	notificationSettings.OrderShipped = req.OrderShipped
	notificationSettings.OrderCancelled = req.OrderCancelled
	notificationSettings.LowStockAlert = req.LowStockAlert
	notificationSettings.LowStockThreshold = req.LowStockThreshold
	notificationSettings.RefundRequested = req.RefundRequested
	notificationSettings.NewReview = req.NewReview
	notificationSettings.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.NotificationSettingsRepo.Save(l.ctx, l.svcCtx.DB, notificationSettings); err != nil {
		l.Logger.Errorf("save notification settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	return toNotificationSettingsResponse(notificationSettings), nil
}
