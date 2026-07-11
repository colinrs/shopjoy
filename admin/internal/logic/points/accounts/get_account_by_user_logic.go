package accounts

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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

	account, err := l.svcCtx.PointsService.GetAccountByUser(l.ctx, req.UserID)
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
