package accounts

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

type GetAccountTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccountTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetAccountTransactionsLogic {
	return GetAccountTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccountTransactionsLogic) GetAccountTransactions(req *types.ListAccountTransactionsReq) (resp *types.ListPointsTransactionsResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	query := points.PointsTransactionQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Type: points.TransactionType(req.Type),
	}

	transactions, total, err := l.svcCtx.PointsService.GetAccountTransactions(l.ctx, shared.TenantID(tenantID), req.ID, query)
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
			CreatedAt:     t.CreatedAt.Format(time.RFC3339),
		}
	}

	return &types.ListPointsTransactionsResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Stats: types.PointsTransactionsStats{
			TotalEarned:   0,
			TotalRedeemed: 0,
		},
	}, nil
}

func formatTimePtrFromTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
