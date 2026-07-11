package seo

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGlobalSEOLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGlobalSEOLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetGlobalSEOLogic {
	return GetGlobalSEOLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGlobalSEOLogic) GetGlobalSEO() (resp *types.GlobalSEOConfigResponse, err error) {

	result, err := l.svcCtx.SEOService.GetGlobalSEO(l.ctx)
	if err != nil {
		return nil, err
	}

	return &types.GlobalSEOConfigResponse{
		SEOConfigDTO: types.SEOConfigDTO{
			Title:       result.Title,
			Description: result.Description,
			Keywords:    result.Keywords,
		},
	}, nil
}
