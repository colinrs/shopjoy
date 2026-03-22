package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromotionRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePromotionRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePromotionRuleLogic {
	return UpdatePromotionRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePromotionRuleLogic) UpdatePromotionRule(req *types.UpdatePromotionRuleReq) (resp *types.PromotionRuleResp, err error) {
	// todo: add your logic here and delete this line

	return
}
