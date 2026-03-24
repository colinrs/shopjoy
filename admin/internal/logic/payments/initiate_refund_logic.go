package payments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitiateRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitiateRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) InitiateRefundLogic {
	return InitiateRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitiateRefundLogic) InitiateRefund(req *types.InitiateRefundReq) (resp *types.InitiateRefundResp, err error) {
	// todo: add your logic here and delete this line

	return
}
