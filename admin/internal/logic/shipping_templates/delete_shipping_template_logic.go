package shipping_templates

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteShippingTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteShippingTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteShippingTemplateLogic {
	return DeleteShippingTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteShippingTemplateLogic) DeleteShippingTemplate(req *types.DeleteShippingTemplateReq) error {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		return code.ErrUnauthorized
	}

	// Find existing template
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return err
	}

	// Check if can delete
	if err := template.CanDelete(); err != nil {
		return err
	}

	// Delete template
	return l.svcCtx.ShippingRepo.Delete(l.ctx, l.svcCtx.DB, tenantID, req.ID)
}