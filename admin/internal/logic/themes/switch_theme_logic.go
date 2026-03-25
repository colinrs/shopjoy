package themes

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchThemeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSwitchThemeLogic(ctx context.Context, svcCtx *svc.ServiceContext) SwitchThemeLogic {
	return SwitchThemeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwitchThemeLogic) SwitchTheme(req *types.SwitchThemeRequest) error {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)
	// userName is not available in context, using empty string
	// TODO: Add GetUserName to contextx or fetch from user service
	userName := ""

	return l.svcCtx.ThemeService.SwitchTheme(l.ctx, shared.TenantID(tenantID), req.ThemeID, userID, userName)
}