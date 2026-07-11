package themes

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListThemeAuditLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListThemeAuditLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListThemeAuditLogsLogic {
	return ListThemeAuditLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListThemeAuditLogsLogic) ListThemeAuditLogs() (resp *types.ListThemeAuditLogsResponse, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Default pagination (API doesn't have request type with pagination yet)
	result, err := l.svcCtx.ThemeService.ListAuditLogs(l.ctx, shared.TenantID(tenantID), 1, 20)
	if err != nil {
		return nil, err
	}

	logs := make([]*types.ThemeAuditLog, 0, len(result.Items))
	for _, l := range result.Items {
		logs = append(logs, &types.ThemeAuditLog{
			ID:        l.ID,
			Action:    l.Action,
			ThemeID:   l.ThemeID,
			ThemeName: l.ThemeName,
			UserID:    l.UserID,
			UserName:  l.UserName,
			CreatedAt: l.CreatedAt,
		})
	}

	return &types.ListThemeAuditLogsResponse{
		Logs:  logs,
		Total: result.Total,
	}, nil
}
