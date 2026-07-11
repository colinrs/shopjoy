package themes

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
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

	userID, _ := contextx.GetUserID(l.ctx)
	userName := contextx.GetCurrentUserName(l.ctx)

	return l.svcCtx.ThemeService.SwitchTheme(l.ctx, req.ThemeID, userID, userName)
}
