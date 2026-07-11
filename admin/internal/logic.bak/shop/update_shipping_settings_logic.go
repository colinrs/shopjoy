package shop

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShippingSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateShippingSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateShippingSettingsLogic {
	return UpdateShippingSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShippingSettingsLogic) UpdateShippingSettings(req *types.UpdateShippingSettingsRequest) (resp *types.ShippingSettings, err error) {
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

	// Find existing shipping settings or create new
	shippingSettings, err := l.svcCtx.ShippingSettingsRepo.FindByShopID(l.ctx, l.svcCtx.DB, settings.ID)
	if err != nil {
		l.Logger.Errorf("find shipping settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if shippingSettings == nil {
		shippingSettings = &shop.ShippingSettings{
			ShopID:   settings.ID,
			Currency: "CNY",
		}
	}

	// Update fields from request
	if req.FreeShippingThreshold != "" {
		threshold, err := decimal.NewFromString(req.FreeShippingThreshold)
		if err != nil {
			return nil, code.ErrShopInvalid
		}
		shippingSettings.FreeShippingThreshold = threshold
	}
	if req.DefaultShippingFee != "" {
		fee, err := decimal.NewFromString(req.DefaultShippingFee)
		if err != nil {
			return nil, code.ErrShopInvalid
		}
		shippingSettings.DefaultShippingFee = fee
	}
	shippingSettings.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.ShippingSettingsRepo.Save(l.ctx, l.svcCtx.DB, shippingSettings); err != nil {
		l.Logger.Errorf("save shipping settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	return toShippingSettingsResponse(shippingSettings), nil
}
