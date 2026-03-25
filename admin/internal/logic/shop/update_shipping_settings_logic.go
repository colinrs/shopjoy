package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShippingSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateShippingSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateShippingSettingsLogic {
	return UpdateShippingSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShippingSettingsLogic) UpdateShippingSettings(req *types.UpdateShippingSettingsRequest) (resp *types.ShippingSettings, err error) {
	// todo: add your logic here and delete this line

	return
}
