package shop

import (
	"context"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShopSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateShopSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateShopSettingsLogic {
	return UpdateShopSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShopSettingsLogic) UpdateShopSettings(req *types.UpdateShopSettingsRequest) (resp *types.ShopSettings, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok || tenantID == 0 {
		return nil, code.ErrTenantInvalidID
	}

	// Find existing settings
	settings, err := l.svcCtx.ShopSettingsRepo.FindByTenantID(l.ctx, l.svcCtx.DB, tenantID)
	if err != nil {
		l.Logger.Errorf("find shop settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if settings == nil {
		// Create new settings if not exists
		settings = &shop.ShopSettings{
			TenantID: tenantID,
			Code:     generateShopCode(tenantID),
			Status:   1,
			Plan:     0,
		}
	}

	// Update fields if provided
	if req.Name != "" {
		settings.Name = req.Name
	}
	if req.Logo != "" {
		settings.Logo = req.Logo
	}
	if req.Description != "" {
		settings.Description = req.Description
	}
	if req.ContactName != "" {
		settings.ContactName = req.ContactName
	}
	if req.ContactPhone != "" {
		settings.ContactPhone = req.ContactPhone
	}
	if req.ContactEmail != "" {
		settings.ContactEmail = req.ContactEmail
	}
	if req.Address != "" {
		settings.Address = req.Address
	}
	if req.CustomDomain != "" {
		settings.CustomDomain = req.CustomDomain
	}
	if req.PrimaryColor != "" {
		settings.PrimaryColor = req.PrimaryColor
	}
	if req.SecondaryColor != "" {
		settings.SecondaryColor = req.SecondaryColor
	}
	if req.Favicon != "" {
		settings.Favicon = req.Favicon
	}
	if req.DefaultLanguage != "" {
		settings.DefaultLanguage = req.DefaultLanguage
	}
	if req.Timezone != "" {
		settings.Timezone = req.Timezone
	}

	settings.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.ShopSettingsRepo.Save(l.ctx, l.svcCtx.DB, settings); err != nil {
		l.Logger.Errorf("save shop settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	return toShopSettingsResponse(settings), nil
}

func generateShopCode(tenantID int64) string {
	return fmt.Sprintf("shop-%d", tenantID)
}
