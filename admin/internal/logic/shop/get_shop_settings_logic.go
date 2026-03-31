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

type GetShopSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShopSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetShopSettingsLogic {
	return GetShopSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShopSettingsLogic) GetShopSettings() (resp *types.ShopSettings, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok || tenantID == 0 {
		return nil, code.ErrTenantInvalidID
	}

	settings, err := l.svcCtx.ShopSettingsRepo.FindByTenantID(l.ctx, l.svcCtx.DB, tenantID)
	if err != nil {
		l.Logger.Errorf("get shop settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if settings == nil {
		return nil, code.ErrShopNotFound
	}

	return toShopSettingsResponse(settings), nil
}

func toShopSettingsResponse(s *shop.ShopSettings) *types.ShopSettings {
	return &types.ShopSettings{
		ID:               s.ID,
		Name:             s.Name,
		Code:             s.Code,
		Logo:             s.Logo,
		Description:      s.Description,
		ContactName:      s.ContactName,
		ContactPhone:     s.ContactPhone,
		ContactEmail:     s.ContactEmail,
		Address:          s.Address,
		Domain:           s.Domain,
		CustomDomain:     s.CustomDomain,
		PrimaryColor:     s.PrimaryColor,
		SecondaryColor:   s.SecondaryColor,
		Favicon:          s.Favicon,
		DefaultCurrency:  s.DefaultCurrency,
		DefaultLanguage:  s.DefaultLanguage,
		Timezone:         s.Timezone,
		Status:           s.Status,
		StatusText:       s.StatusText(),
		Plan:             s.Plan,
		PlanText:         s.PlanText(),
		ExpireAt:         s.ExpireAt,
		CreatedAt:        s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        s.UpdatedAt.Format(time.RFC3339),
	}
}
