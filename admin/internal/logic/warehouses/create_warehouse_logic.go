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
	//
	// FIX-2: platform admins (TenantID == 0 or missing from ctx) are now
	// allowed through this guard — the GORM tenant middleware is responsible
	// for platform-scope filtering, so this layer must not reject them. We
	// still resolve the value here because the row's tenant_id field has to
	// be populated: a normal user gets their TenantID, while a platform
	// admin who did not pass X-Tenant-ID falls back to 0 (representing
	// "platform-owned") so the row remains queryable from the admin console.
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Check for duplicate code within the tenant (codes are scoped per-tenant,
	// so the same code may exist across different tenants). When called by a
	// platform admin (tenantID == 0) this scopes the duplicate check to
	// platform-owned rows; ordinary users can never collide with them because
	// their TenantID is non-zero.
	existing, err := l.svcCtx.WarehouseRepo.FindByCode(l.ctx, l.svcCtx.DB, tenantID, req.Code)
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
