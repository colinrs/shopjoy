package warehouses

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateWarehouseLogic {
	return CreateWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWarehouseLogic) CreateWarehouse(req *types.CreateWarehouseReq) (resp *types.CreateWarehouseResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Check for duplicate code
	existing, err := l.svcCtx.WarehouseRepo.FindByCode(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, code.ErrInventoryDuplicateWarehouseCode
	}

	// Generate ID
	id, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		return nil, err
	}

	warehouse := &product.Warehouse{
		ID:        id,
		TenantID:  shared.TenantID(tenantID),
		Code:      req.Code,
		Name:      req.Name,
		Country:   req.Country,
		Address:   req.Address,
		IsDefault: req.IsDefault,
		Status:    shared.StatusEnabled,
		Audit: shared.AuditInfo{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	if err := l.svcCtx.WarehouseRepo.Create(l.ctx, l.svcCtx.DB, warehouse); err != nil {
		return nil, err
	}

	return &types.CreateWarehouseResp{
		ID: warehouse.ID,
	}, nil
}
