package transactions

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTransactionLogic {
	return GetTransactionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransactionLogic) GetTransaction(req *types.GetPointsTransactionReq) (resp *types.PointsTransaction, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	transaction, err := l.svcCtx.PointsService.GetTransaction(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return &types.PointsTransaction{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
		AccountID:     transaction.AccountID,
		Points:        transaction.Points,
		BalanceAfter:  transaction.BalanceAfter,
		Type:          transaction.Type,
		ReferenceType: transaction.ReferenceType,
		ReferenceID:   transaction.ReferenceID,
		Description:   transaction.Description,
		ExpiresAt:     formatTimePtrFromTime(transaction.ExpiresAt),
		CreatedAt:     transaction.CreatedAt.Format(time.RFC3339),
	}, nil
}
