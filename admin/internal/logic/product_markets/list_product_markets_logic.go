// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductMarketsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品市场配置列表
func NewListProductMarketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductMarketsLogic {
	return &ListProductMarketsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductMarketsLogic) ListProductMarkets() (resp *types.ListProductMarketsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
