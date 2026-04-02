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

type ListAccountsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAccountsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListAccountsLogic {
	return ListAccountsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAccountsLogic) ListAccounts(req *types.ListAccountsReq) (resp *types.ListAccountsResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	query := points.PointsAccountQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		UserID: req.UserID,
		Email:  req.Email,
	}

	accounts, total, stats, err := l.svcCtx.PointsService.ListAccounts(l.ctx, shared.TenantID(tenantID), query)
	if err != nil {
		return nil, err
	}

	list := make([]*types.PointsAccount, len(accounts))
	for i, a := range accounts {
		list[i] = &types.PointsAccount{
			ID:            a.ID,
			UserID:        a.UserID,
			UserEmail:     "", // Would need to fetch from user service
			Balance:       a.Balance,
			FrozenBalance: a.FrozenBalance,
			TotalEarned:   a.TotalEarned,
			TotalRedeemed: a.TotalRedeemed,
			TotalExpired:  a.TotalExpired,
			CreatedAt:     a.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     a.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &types.ListAccountsResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Stats: types.AccountsStats{
			Total:        stats.Total,
			TotalBalance: stats.TotalBalance,
			Active:       stats.Active,
		},
	}, nil
}
