package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
