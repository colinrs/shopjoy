package shipping_mappings

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTemplateMappingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTemplateMappingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListTemplateMappingsLogic {
	return ListTemplateMappingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTemplateMappingsLogic) ListTemplateMappings(req *types.ListTemplateMappingsReq) (resp *types.ListTemplateMappingsResp, err error) {
	// Get tenant ID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Verify template exists and belongs to tenant
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.TemplateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, code.ErrShippingTemplateNotFound
	}

	// Get mappings
	mappings, err := l.svcCtx.ShippingRepo.FindMappingsByTemplateID(l.ctx, l.svcCtx.DB, req.TemplateID)
	if err != nil {
		return nil, err
	}

	// Build response
	list := make([]*types.TemplateMappingDetail, 0, len(mappings))
	for _, m := range mappings {
		list = append(list, &types.TemplateMappingDetail{
			ID:         int64(m.ID),
			TemplateID: m.TemplateID,
			TargetType: string(m.TargetType),
			TargetID:   m.TargetID,
		})
	}

	return &types.ListTemplateMappingsResp{
		List: list,
	}, nil
}
