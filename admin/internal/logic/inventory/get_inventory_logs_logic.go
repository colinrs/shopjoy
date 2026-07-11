package inventory

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetInventoryLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInventoryLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetInventoryLogsLogic {
	return GetInventoryLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInventoryLogsLogic) GetInventoryLogs(req *types.GetInventoryLogsReq) (resp *types.ListInventoryLogsResp, err error) {

	query := product.InventoryLogQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		ProductID:  req.ProductID,
		SKUCode:    req.SKUCode,
		ChangeType: req.Type,
	}

	logs, total, err := l.svcCtx.InventoryLogRepo.Find(l.ctx, l.svcCtx.DB, query)
	if err != nil {
		return nil, err
	}

	list := make([]*types.InventoryLogResp, 0, len(logs))
	for _, log := range logs {
		list = append(list, &types.InventoryLogResp{
			ID:             log.ID,
			SKUCode:        log.SKUCode,
			ProductID:      log.ProductID,
			WarehouseID:    log.WarehouseID,
			ChangeType:     log.ChangeType,
			ChangeQuantity: log.ChangeQuantity,
			BeforeStock:    log.BeforeStock,
			AfterStock:     log.AfterStock,
			OrderNo:        log.OrderNo,
			Remark:         log.Remark,
			OperatorID:     log.OperatorID,
			CreatedAt:      log.CreatedAt.Format(time.RFC3339),
		})
	}

	return &types.ListInventoryLogsResp{
		List:  list,
		Total: total,
	}, nil
}
