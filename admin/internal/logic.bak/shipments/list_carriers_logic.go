package shipments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCarriersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCarriersLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListCarriersLogic {
	return ListCarriersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCarriersLogic) ListCarriers(req *types.ListCarriersReq) (resp *types.ListCarriersResp, err error) {
	// Get carriers
	carriers, err := l.svcCtx.CarrierApp.ListCarriers(l.ctx)
	if err != nil {
		return nil, err
	}

	// Build response
	resp = &types.ListCarriersResp{
		List: make([]*types.CarrierResp, len(carriers)),
	}

	for i, c := range carriers {
		resp.List[i] = &types.CarrierResp{
			ID:          c.ID,
			Code:        c.Code,
			Name:        c.Name,
			TrackingURL: c.TrackingURL,
			IsActive:    c.IsActive,
			Sort:        c.Sort,
		}
	}

	return resp, nil
}
