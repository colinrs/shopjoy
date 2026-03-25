package shipping_mappings

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTemplateMappingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTemplateMappingLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteTemplateMappingLogic {
	return DeleteTemplateMappingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTemplateMappingLogic) DeleteTemplateMapping(req *types.DeleteTemplateMappingReq) error {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		return code.ErrUnauthorized
	}

	// Find existing mapping
	mapping, err := l.svcCtx.ShippingRepo.FindMappingByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return err
	}

	// Verify template belongs to tenant through the mapping
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, mapping.TemplateID)
	if err != nil {
		return code.ErrShippingTemplateNotFound
	}
	if template == nil {
		return code.ErrShippingTemplateNotFound
	}

	// Delete mapping
	return l.svcCtx.ShippingRepo.DeleteMapping(l.ctx, l.svcCtx.DB, req.ID)
}