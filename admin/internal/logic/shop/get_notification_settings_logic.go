package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetNotificationSettingsLogic {
	return GetNotificationSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotificationSettingsLogic) GetNotificationSettings() (resp *types.NotificationSettings, err error) {
	// todo: add your logic here and delete this line

	return
}
