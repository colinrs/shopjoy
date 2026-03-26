package accounts

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetAccountLogic {
	return GetAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccountLogic) GetAccount(req *types.GetAccountReq) (resp *types.PointsAccount, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	account, err := l.svcCtx.PointsService.GetAccount(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return &types.PointsAccount{
		ID:            account.ID,
		UserID:        account.UserID,
		UserEmail:     "",
		Balance:       account.Balance,
		FrozenBalance: account.FrozenBalance,
		TotalEarned:   account.TotalEarned,
		TotalRedeemed: account.TotalRedeemed,
		TotalExpired:  account.TotalExpired,
		CreatedAt:     account.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     account.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}