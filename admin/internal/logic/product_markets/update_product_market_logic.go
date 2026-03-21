// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新商品市场配置
func NewUpdateProductMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductMarketLogic {
	return &UpdateProductMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductMarketLogic) UpdateProductMarket(req *types.UpdateProductMarketReq) (resp *types.ProductMarketResp, err error) {
	// todo: add your logic here and delete this line

	return
}
