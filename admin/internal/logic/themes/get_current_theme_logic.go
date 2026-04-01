package themes

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentThemeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentThemeLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCurrentThemeLogic {
	return GetCurrentThemeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentThemeLogic) GetCurrentTheme() (resp *types.CurrentThemeResponse, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	result, err := l.svcCtx.ThemeService.GetCurrentTheme(l.ctx, shared.TenantID(tenantID))
	if err != nil {
		return nil, err
	}

	if result == nil || result.Theme == nil {
		return &types.CurrentThemeResponse{}, nil
	}

	return &types.CurrentThemeResponse{
		Theme: &types.ThemeItem{
			ID:           result.Theme.ID,
			Code:         result.Theme.Code,
			Name:         result.Theme.Name,
			Description:  result.Theme.Description,
			PreviewImage: result.Theme.PreviewImage,
			Thumbnail:    result.Theme.Thumbnail,
			IsPreset:     result.Theme.IsPreset,
			IsCurrent:    result.Theme.IsCurrent,
		},
		Config: types.ThemeConfigDTO{
			PrimaryColor:   result.Config.PrimaryColor,
			SecondaryColor: result.Config.SecondaryColor,
			FontFamily:     result.Config.FontFamily,
			ButtonStyle:    result.Config.ButtonStyle,
		},
	}, nil
}