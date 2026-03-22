package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromotionRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPromotionRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPromotionRulesLogic {
	return GetPromotionRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPromotionRulesLogic) GetPromotionRules(req *types.GetPromotionRulesReq) (resp *types.ListPromotionRulesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
