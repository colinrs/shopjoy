package refunds

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRefundReasonsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRefundReasonsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListRefundReasonsLogic {
	return ListRefundReasonsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRefundReasonsLogic) ListRefundReasons(req *types.ListRefundReasonsReq) (resp *types.ListRefundReasonsResp, err error) {
	result, err := l.svcCtx.RefundReasonApp.ListRefundReasons(l.ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*types.RefundReasonResp, len(result.List))
	for i, r := range result.List {
		list[i] = convertRefundReasonToResp(r)
	}

	return &types.ListRefundReasonsResp{
		List: list,
	}, nil
}
