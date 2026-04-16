package accounts

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountByUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccountByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetAccountByUserLogic {
	return GetAccountByUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccountByUserLogic) GetAccountByUser(req *types.GetAccountByUserReq) (resp *types.PointsAccount, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	account, err := l.svcCtx.PointsService.GetAccountByUser(l.ctx, shared.TenantID(tenantID), req.UserID)
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
		CreatedAt:     account.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     account.UpdatedAt.Format(time.RFC3339),
	}, nil
}
