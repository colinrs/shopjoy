package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
