package brands

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

type CreateBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateBrandLogic {
	return CreateBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBrandLogic) CreateBrand(req *types.CreateBrandReq) (resp *types.CreateBrandResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Check for duplicate brand name
	existing, err := l.svcCtx.BrandRepo.FindByName(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, code.ErrBrandDuplicate
	}

	// Generate ID
	id, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		return nil, err
	}

	brand := &product.Brand{
		Model:            application.Model{ID: id, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		TenantID:         shared.TenantID(tenantID),
		Name:             req.Name,
		Logo:             req.Logo,
		Description:      req.Description,
		Website:          req.Website,
		TrademarkNumber:  req.TrademarkNumber,
		TrademarkCountry: req.TrademarkCountry,
		EnablePage:       req.EnablePage,
		Sort:             req.Sort,
		Status:           shared.StatusEnabled,
	}

	if err := l.svcCtx.BrandRepo.Create(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	return &types.CreateBrandResp{
		ID: brand.Model.ID,
	}, nil
}
