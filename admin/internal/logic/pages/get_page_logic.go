package pages

import (
	"context"
	"encoding/json"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPageLogic {
	return GetPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPageLogic) GetPage(req *types.GetPageRequest) (resp *types.PageDetailResponse, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	result, err := l.svcCtx.PageService.GetPage(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	if result == nil || result.Page == nil {
		return &types.PageDetailResponse{}, nil
	}

	decorations := make([]*types.DecorationDTO, 0, len(result.Decorations))
	for _, d := range result.Decorations {
		blockConfigJSON, _ := json.Marshal(d.BlockConfig)
		decorations = append(decorations, &types.DecorationDTO{
			ID:          d.ID,
			BlockType:   d.BlockType,
			BlockConfig: string(blockConfigJSON),
			SortOrder:   d.SortOrder,
		})
	}

	return &types.PageDetailResponse{
		Page: &types.PageListItem{
			ID:          result.Page.ID,
			PageType:    result.Page.PageType,
			Name:        result.Page.Name,
			Slug:        result.Page.Slug,
			IsPublished: result.Page.IsPublished,
			Version:     result.Page.Version,
		},
		Decorations: decorations,
	}, nil
}