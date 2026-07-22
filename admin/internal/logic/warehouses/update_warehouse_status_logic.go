package warehouses

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWarehouseStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWarehouseStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateWarehouseStatusLogic {
	return UpdateWarehouseStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWarehouseStatusLogic) UpdateWarehouseStatus(req *types.UpdateWarehouseStatusReq) (resp *types.WarehouseDetailResp, err error) {

	// Resolve and authorize: rejects missing TenantID and cross-tenant access.
	warehouse, err := requireWarehouseOwnership(l.ctx, l.svcCtx, req.ID)
	if err != nil {
		return nil, err
	}

	// Update status
	if req.Status == 1 {
		warehouse.Enable()
	} else {
		warehouse.Disable()
	}
	warehouse.Model.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.WarehouseRepo.Update(l.ctx, l.svcCtx.DB, warehouse); err != nil {
		return nil, err
	}

	return toWarehouseDetailResp(warehouse), nil
}
