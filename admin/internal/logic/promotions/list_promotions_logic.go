package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPromotionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPromotionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListPromotionsLogic {
	return ListPromotionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPromotionsLogic) ListPromotions(req *types.ListPromotionsReq) (resp *types.ListPromotionsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
