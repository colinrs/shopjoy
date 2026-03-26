package warehouses

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
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Find warehouse
	warehouse, err := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	// Update status
	if req.Status == 1 {
		warehouse.Enable()
	} else {
		warehouse.Disable()
	}
	warehouse.Audit.UpdatedAt = time.Now().Unix()

	if err := l.svcCtx.WarehouseRepo.Update(l.ctx, l.svcCtx.DB, warehouse); err != nil {
		return nil, err
	}

	return toWarehouseDetailResp(warehouse), nil
}
