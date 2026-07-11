package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLocalizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLocalizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteProductLocalizationLogic {
	return DeleteProductLocalizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductLocalizationLogic) DeleteProductLocalization(req *types.GetProductLocalizationReq) (resp *types.CreateProductResp, err error) {
	// Get tenant ID from context

	// Delete localization
	if err := l.svcCtx.ProductLocalizationRepo.Delete(l.ctx, l.svcCtx.DB, req.ID); err != nil {
		return nil, err
	}

	return &types.CreateProductResp{ID: req.ID}, nil
}
