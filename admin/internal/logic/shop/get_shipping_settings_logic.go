package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShippingSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShippingSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetShippingSettingsLogic {
	return GetShippingSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShippingSettingsLogic) GetShippingSettings() (resp *types.ShippingSettings, err error) {
	// todo: add your logic here and delete this line

	return
}
