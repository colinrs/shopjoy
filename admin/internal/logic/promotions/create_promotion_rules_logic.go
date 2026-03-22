package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePromotionRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePromotionRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreatePromotionRulesLogic {
	return CreatePromotionRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePromotionRulesLogic) CreatePromotionRules(req *types.CreatePromotionRulesReq) (resp *types.CreatePromotionRulesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
