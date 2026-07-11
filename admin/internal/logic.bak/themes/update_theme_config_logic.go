package themes

import (
	"context"

	appStorefront "github.com/colinrs/shopjoy/admin/internal/application/storefront"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateThemeConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateThemeConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateThemeConfigLogic {
	return UpdateThemeConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateThemeConfigLogic) UpdateThemeConfig(req *types.UpdateThemeConfigRequest) error {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	userID, _ := contextx.GetUserID(l.ctx)
	userName := contextx.GetCurrentUserName(l.ctx)

	config := appStorefront.ThemeConfigDTO{
		PrimaryColor:   req.Config.PrimaryColor,
		SecondaryColor: req.Config.SecondaryColor,
		FontFamily:     req.Config.FontFamily,
		ButtonStyle:    req.Config.ButtonStyle,
	}

	return l.svcCtx.ThemeService.UpdateThemeConfig(l.ctx, shared.TenantID(tenantID), config, userID, userName)
}
