package warehouses

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/application"
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

	// Resolve TenantID from ctx (injected by AuthMiddleware).
	// TenantID == 0 means either AuthMiddleware was bypassed or the caller
	// is a platform admin without an explicit X-Tenant-ID header — in either
	// case we must not write a row with tenant_id=0 because that breaks
	// cross-tenant isolation.
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok || tenantID == 0 {
		return nil, code.ErrTenantNotFound
	}

	// Check for duplicate code
	existing, err := l.svcCtx.WarehouseRepo.FindByCode(l.ctx, l.svcCtx.DB, req.Code)
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
		Model:     application.Model{ID: id, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		TenantID:  shared.TenantID(tenantID),
		Code:      req.Code,
		Name:      req.Name,
		Country:   req.Country,
		Address:   req.Address,
		IsDefault: req.IsDefault,
		Status:    shared.StatusEnabled,
	}

	if err := l.svcCtx.WarehouseRepo.Create(l.ctx, l.svcCtx.DB, warehouse); err != nil {
		return nil, err
	}

	return &types.CreateWarehouseResp{
		ID: warehouse.Model.ID,
	}, nil
}