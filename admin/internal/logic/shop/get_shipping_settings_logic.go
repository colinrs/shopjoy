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

type GetShippingSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShippingSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetShippingSettingsLogic {
	return GetShippingSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShippingSettingsLogic) GetShippingSettings() (resp *types.ShippingSettings, err error) {
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

	shippingSettings, err := l.svcCtx.ShippingSettingsRepo.FindByShopID(l.ctx, l.svcCtx.DB, settings.ID)
	if err != nil {
		l.Logger.Errorf("get shipping settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if shippingSettings == nil {
		// Return default settings if not configured
		return &types.ShippingSettings{
			FreeShippingThreshold: "0.00",
			DefaultShippingFee:    "0.00",
			Currency:             "CNY",
		}, nil
	}

	return toShippingSettingsResponse(shippingSettings), nil
}

func toShippingSettingsResponse(s *shop.ShippingSettings) *types.ShippingSettings {
	return &types.ShippingSettings{
		FreeShippingThreshold: s.FreeShippingThreshold.StringFixed(2),
		DefaultShippingFee:    s.DefaultShippingFee.StringFixed(2),
		Currency:              s.Currency,
	}
}
