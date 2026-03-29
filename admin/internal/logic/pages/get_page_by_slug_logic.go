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

type GetPageBySlugLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPageBySlugLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPageBySlugLogic {
	return GetPageBySlugLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPageBySlugLogic) GetPageBySlug(req *types.GetPageBySlugRequest) (resp *types.PageDetailResponse, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	result, err := l.svcCtx.PageService.GetPageBySlug(l.ctx, shared.TenantID(tenantID), req.Slug)
	if err != nil {
		return nil, err
	}

	if result == nil || result.Page == nil {
		return &types.PageDetailResponse{}, nil
	}

	decorations := make([]*types.DecorationDTO, 0, len(result.Decorations))
	for _, d := range result.Decorations {
		blockConfigJSON, err := json.Marshal(d.BlockConfig)
		if err != nil {
			l.Logger.Errorf("failed to marshal block config for decoration %d: %v", d.ID, err)
			continue
		}
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
