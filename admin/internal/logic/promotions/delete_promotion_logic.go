package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeletePromotionLogic {
	return DeletePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePromotionLogic) DeletePromotion(req *types.DeletePromotionReq) (resp *types.CreatePromotionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
