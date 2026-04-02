package themes

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListThemesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListThemesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListThemesLogic {
	return ListThemesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListThemesLogic) ListThemes() (resp *types.ListThemesResponse, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	themes, err := l.svcCtx.ThemeService.ListThemes(l.ctx, shared.TenantID(tenantID))
	if err != nil {
		return nil, err
	}

	items := make([]*types.ThemeItem, 0, len(themes))
	for _, t := range themes {
		items = append(items, &types.ThemeItem{
			ID:           t.ID,
			Code:         t.Code,
			Name:         t.Name,
			Description:  t.Description,
			PreviewImage: t.PreviewImage,
			Thumbnail:    t.Thumbnail,
			IsPreset:     t.IsPreset,
			IsCurrent:    t.IsCurrent,
		})
	}

	return &types.ListThemesResponse{
		Themes: items,
	}, nil
}
