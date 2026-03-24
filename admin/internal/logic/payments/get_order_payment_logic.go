package payments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetOrderPaymentLogic {
	return GetOrderPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderPaymentLogic) GetOrderPayment(req *types.GetOrderPaymentReq) (resp *types.OrderPaymentResp, err error) {
	// todo: add your logic here and delete this line

	return
}
