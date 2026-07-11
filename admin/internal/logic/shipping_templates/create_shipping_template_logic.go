package shipping_templates

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateShippingTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateShippingTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateShippingTemplateLogic {
	return CreateShippingTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShippingTemplateLogic) CreateShippingTemplate(req *types.CreateShippingTemplateReq) (resp *types.CreateShippingTemplateResp, err error) {

	// Create template entity
	template := &shipping.ShippingTemplate{
		Name:      req.Name,
		IsDefault: req.IsDefault,
		IsActive:  true,
	}

	// If setting as default, unset other defaults first
	if req.IsDefault {
		if err := l.svcCtx.ShippingRepo.UnsetAllDefault(l.ctx, l.svcCtx.DB); err != nil {
			return nil, err
		}
	}

	// Save template
	if err := l.svcCtx.ShippingRepo.Create(l.ctx, l.svcCtx.DB, template); err != nil {
		return nil, err
	}

	return &types.CreateShippingTemplateResp{
		ID:   int64(template.ID),
		Name: template.Name,
	}, nil
}
