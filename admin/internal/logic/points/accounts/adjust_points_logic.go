package accounts

import (
	"context"

	apppoints "github.com/colinrs/shopjoy/admin/internal/application/points"
	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdjustPointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdjustPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdjustPointsLogic {
	return AdjustPointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdjustPointsLogic) AdjustPoints(req *types.AdjustPointsReq) (resp *types.AdjustPointsResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)

	adjustReq := apppoints.AdjustPointsRequest{
		AccountID:      req.ID,
		AdjustmentType: points.AdjustmentType(req.AdjustmentType),
		Points:         req.Points,
		Reason:         req.Reason,
		OperatorID:     userID,
	}

	transaction, err := l.svcCtx.PointsService.AdjustPoints(l.ctx, shared.TenantID(tenantID), adjustReq)
	if err != nil {
		return nil, err
	}

	return &types.AdjustPointsResp{
		TransactionID: transaction.ID,
		Points:        transaction.Points,
		BalanceAfter:  transaction.BalanceAfter,
	}, nil
}
