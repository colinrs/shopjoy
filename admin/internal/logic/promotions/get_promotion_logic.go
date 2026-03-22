package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPromotionLogic {
	return GetPromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPromotionLogic) GetPromotion(req *types.GetPromotionReq) (resp *types.PromotionDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
