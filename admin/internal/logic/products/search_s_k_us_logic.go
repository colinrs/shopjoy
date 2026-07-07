// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSKUsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索SKU（下拉选择用）
func NewSearchSKUsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSKUsLogic {
	return &SearchSKUsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchSKUsLogic) SearchSKUs(req *types.SearchSKUsReq) (resp *types.SearchSKUsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
