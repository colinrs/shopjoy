package seo

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPageSEOLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPageSEOLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPageSEOLogic {
	return GetPageSEOLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPageSEOLogic) GetPageSEO(req *types.GetPageSEORequest) (resp *types.PageSEOConfigResponse, err error) {

	result, err := l.svcCtx.SEOService.GetPageSEO(l.ctx, req.PageType, nil)
	if err != nil {
		return nil, err
	}

	return &types.PageSEOConfigResponse{
		PageType: result.PageType,
		PageID:   result.PageID,
		SEOConfigDTO: types.SEOConfigDTO{
			Title:       result.Config.Title,
			Description: result.Config.Description,
			Keywords:    result.Config.Keywords,
		},
	}, nil
}
