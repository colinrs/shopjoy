package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteMarketLogic {
	return DeleteMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMarketLogic) DeleteMarket(req *types.GetMarketReq) error {
	// todo: add your logic here and delete this line

	return nil
}
