package themes

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/storefront"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
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

	themes, err := l.svcCtx.ThemeService.ListThemes(l.ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*types.ThemeItem, 0, len(themes))
	for _, t := range themes {
		items = append(items, &types.ThemeItem{
			ID:            t.ID,
			Code:          t.Code,
			Name:          t.Name,
			Description:   t.Description,
			PreviewImage:  t.PreviewImage,
			Thumbnail:     t.Thumbnail,
			IsPreset:      t.IsPreset,
			IsCurrent:     t.IsCurrent,
			DefaultConfig: defaultConfigToWire(t.DefaultConfig),
		})
	}

	return &types.ListThemesResponse{
		Themes: items,
	}, nil
}

// defaultConfigToWire converts the application-layer *storefront.ThemeConfigDTO
// to the wire *types.ThemeConfigDTO. Returns nil when src is nil so the JSON
// `omitempty` tag still hides the field for themes that lack a default config.
func defaultConfigToWire(src *storefront.ThemeConfigDTO) *types.ThemeConfigDTO {
	if src == nil {
		return nil
	}
	return &types.ThemeConfigDTO{
		PrimaryColor:   src.PrimaryColor,
		SecondaryColor: src.SecondaryColor,
		FontFamily:     src.FontFamily,
		ButtonStyle:    src.ButtonStyle,
	}
}