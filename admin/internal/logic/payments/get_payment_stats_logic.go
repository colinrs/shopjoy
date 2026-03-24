package payments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPaymentStatsLogic {
	return GetPaymentStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentStatsLogic) GetPaymentStats(req *types.GetPaymentStatsReq) (resp *types.PaymentStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
