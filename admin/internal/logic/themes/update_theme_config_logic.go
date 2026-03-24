package themes

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	appStorefront "github.com/colinrs/shopjoy/admin/internal/application/storefront"
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
	tenantID, _ := contextx.GetTenantID(l.ctx)

	config := appStorefront.ThemeConfigDTO{
		PrimaryColor:   req.Config.PrimaryColor,
		SecondaryColor: req.Config.SecondaryColor,
		FontFamily:     req.Config.FontFamily,
		ButtonStyle:    req.Config.ButtonStyle,
	}

	return l.svcCtx.ThemeService.UpdateThemeConfig(l.ctx, shared.TenantID(tenantID), config)
}