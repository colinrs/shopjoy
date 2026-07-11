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

type GetPaymentSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPaymentSettingsLogic {
	return GetPaymentSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentSettingsLogic) GetPaymentSettings() (resp *types.PaymentSettings, err error) {
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

	paymentSettings, err := l.svcCtx.PaymentSettingsRepo.FindByShopID(l.ctx, l.svcCtx.DB, settings.ID)
	if err != nil {
		l.Logger.Errorf("get payment settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if paymentSettings == nil {
		// Return default settings if not configured
		return &types.PaymentSettings{
			StripeEnabled: false,
		}, nil
	}

	return toPaymentSettingsResponse(paymentSettings), nil
}

func toPaymentSettingsResponse(p *shop.PaymentSettings) *types.PaymentSettings {
	return &types.PaymentSettings{
		StripeEnabled:   p.StripeEnabled,
		StripePublicKey: p.StripePublicKey,
	}
}
