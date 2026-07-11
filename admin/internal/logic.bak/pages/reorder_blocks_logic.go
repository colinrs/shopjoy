package pages

import (
	"context"

	appStorefront "github.com/colinrs/shopjoy/admin/internal/application/storefront"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReorderBlocksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReorderBlocksLogic(ctx context.Context, svcCtx *svc.ServiceContext) ReorderBlocksLogic {
	return ReorderBlocksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReorderBlocksLogic) ReorderBlocks(req *types.ReorderBlocksRequest) error {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	orders := make([]appStorefront.BlockOrderDTO, 0, len(req.BlockOrders))
	for _, o := range req.BlockOrders {
		orders = append(orders, appStorefront.BlockOrderDTO{
			ID:        o.ID,
			SortOrder: o.SortOrder,
		})
	}

	return l.svcCtx.DecorationService.ReorderBlocks(l.ctx, shared.TenantID(tenantID), req.PageID, orders)
}
