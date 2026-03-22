package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActivatePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActivatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) ActivatePromotionLogic {
	return ActivatePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActivatePromotionLogic) ActivatePromotion(req *types.ActivatePromotionReq) (resp *types.PromotionDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
