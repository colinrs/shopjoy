package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePromotionLogic {
	return UpdatePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePromotionLogic) UpdatePromotion(req *types.UpdatePromotionReq) (resp *types.PromotionDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
