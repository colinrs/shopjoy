package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePromotionRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePromotionRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeletePromotionRuleLogic {
	return DeletePromotionRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePromotionRuleLogic) DeletePromotionRule(req *types.DeletePromotionRuleReq) (resp *types.CreatePromotionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
