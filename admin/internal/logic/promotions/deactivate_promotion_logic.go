package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeactivatePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeactivatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeactivatePromotionLogic {
	return DeactivatePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeactivatePromotionLogic) DeactivatePromotion(req *types.DeactivatePromotionReq) (resp *types.PromotionDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
