package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNotificationSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNotificationSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateNotificationSettingsLogic {
	return UpdateNotificationSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNotificationSettingsLogic) UpdateNotificationSettings(req *types.UpdateNotificationSettingsRequest) (resp *types.NotificationSettings, err error) {
	// todo: add your logic here and delete this line

	return
}
