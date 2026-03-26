package transactions

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListTransactionsLogic {
	return ListTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTransactionsLogic) ListTransactions(req *types.ListPointsTransactionsReq) (resp *types.ListPointsTransactionsResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	query := points.PointsTransactionQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		UserID:    req.UserID,
		AccountID: req.AccountID,
		Type:      points.TransactionType(req.Type),
	}

	if req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			query.StartTime = &t
		}
	}
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			query.EndTime = &t
		}
	}

	transactions, total, stats, err := l.svcCtx.PointsService.ListTransactions(l.ctx, shared.TenantID(tenantID), query)
	if err != nil {
		return nil, err
	}

	list := make([]*types.PointsTransaction, len(transactions))
	for i, t := range transactions {
		list[i] = &types.PointsTransaction{
			ID:            t.ID,
			UserID:        t.UserID,
			AccountID:     t.AccountID,
			Points:        t.Points,
			BalanceAfter:  t.BalanceAfter,
			Type:          t.Type,
			ReferenceType: t.ReferenceType,
			ReferenceID:   t.ReferenceID,
			Description:   t.Description,
			ExpiresAt:     formatTimePtrFromTime(t.ExpiresAt),
			CreatedAt:     t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &types.ListPointsTransactionsResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Stats: types.PointsTransactionsStats{
			TotalEarned:   stats.TotalEarned,
			TotalRedeemed: stats.TotalRedeemed,
		},
	}, nil
}

func formatTimePtrFromTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}