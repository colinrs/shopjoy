package inventory

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	query := product.InventoryLogQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		SKUCode:    req.SKUCode,
		ChangeType: req.Type,
	}

	var logs []*product.InventoryLog
	var total int64

	if req.ProductID > 0 {
		logs, total, err = l.svcCtx.InventoryLogRepo.FindByProduct(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ProductID, query)
	} else if req.SKUCode != "" {
		logs, total, err = l.svcCtx.InventoryLogRepo.FindBySKU(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.SKUCode, query)
	} else {
		// Return empty result if no filter
		return &types.ListInventoryLogsResp{
			List:  []*types.InventoryLogResp{},
			Total: 0,
		}, nil
	}

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
