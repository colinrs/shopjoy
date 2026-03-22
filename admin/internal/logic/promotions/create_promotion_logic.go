package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreatePromotionLogic {
	return CreatePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePromotionLogic) CreatePromotion(req *types.CreatePromotionReq) (resp *types.CreatePromotionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
