package shop

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaymentSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePaymentSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePaymentSettingsLogic {
	return UpdatePaymentSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePaymentSettingsLogic) UpdatePaymentSettings(req *types.UpdatePaymentSettingsRequest) (resp *types.PaymentSettings, err error) {

	// Find shop settings first to get shop ID
	settings, err := l.svcCtx.ShopSettingsRepo.FindByTenantID(l.ctx, l.svcCtx.DB)
	if err != nil {
		l.Logger.Errorf("find shop settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if settings == nil {
		return nil, code.ErrShopNotFound
	}

	// Find existing payment settings or create new
	paymentSettings, err := l.svcCtx.PaymentSettingsRepo.FindByShopID(l.ctx, l.svcCtx.DB, settings.ID)
	if err != nil {
		l.Logger.Errorf("find payment settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if paymentSettings == nil {
		paymentSettings = &shop.PaymentSettings{
			ShopID: settings.ID,
		}
	}

	// Update fields from request
	paymentSettings.StripeEnabled = req.StripeEnabled
	if req.StripeSecretKey != "" {
		paymentSettings.StripeSecretKey = req.StripeSecretKey
	}
	paymentSettings.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.PaymentSettingsRepo.Save(l.ctx, l.svcCtx.DB, paymentSettings); err != nil {
		l.Logger.Errorf("save payment settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	return toPaymentSettingsResponse(paymentSettings), nil
}
