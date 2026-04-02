package refunds

import (
	"context"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRefundsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRefundsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListRefundsLogic {
	return ListRefundsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRefundsLogic) ListRefunds(req *types.ListRefundsReq) (resp *types.ListRefundsResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Parse time filters
	var startTime, endTime time.Time
	if req.StartTime != "" {
		startTime, err = parseTime(req.StartTime)
		if err != nil {
			return nil, err
		}
	}
	if req.EndTime != "" {
		endTime, err = parseTime(req.EndTime)
		if err != nil {
			return nil, err
		}
	}

	queryReq := appfulfillment.QueryRefundRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		OrderID:    req.OrderID,
		UserID:     req.UserID,
		Status:     fulfillment.ParseRefundStatus(req.Status),
		ReasonType: req.ReasonType,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	listResp, err := l.svcCtx.RefundApp.ListRefunds(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return nil, err
	}

	list := make([]*types.RefundDetailResp, len(listResp.List))
	for i, r := range listResp.List {
		list[i] = convertRefundToDetailResp(r)
	}

	return &types.ListRefundsResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
