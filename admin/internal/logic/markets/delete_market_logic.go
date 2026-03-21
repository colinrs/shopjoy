package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMarketLogic {
	return &DeleteMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMarketLogic) DeleteMarket(req *types.GetMarketReq) error {
	repo := persistence.NewMarketRepository()
	m, err := repo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return err
	}

	// Cannot delete default market
	if m.IsDefault {
		return code.ErrMarketCannotDelete
	}

	if err := repo.Delete(l.ctx, l.svcCtx.DB, req.ID); err != nil {
		return err
	}

	return nil
}