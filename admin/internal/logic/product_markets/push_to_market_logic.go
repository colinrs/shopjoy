// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushToMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 推送商品到市场
func NewPushToMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushToMarketLogic {
	return &PushToMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushToMarketLogic) PushToMarket(req *types.PushToMarketReq) (resp *types.PushToMarketResp, err error) {
	// todo: add your logic here and delete this line

	return
}
