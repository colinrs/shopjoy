package versions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListVersionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListVersionsLogic {
	return ListVersionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListVersionsLogic) ListVersions(req *types.ListVersionsRequest) (resp *types.ListVersionsResponse, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	result, err := l.svcCtx.VersionService.ListVersions(l.ctx, shared.TenantID(tenantID), req.PageID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*types.VersionListItem, 0, len(result.Items))
	for _, v := range result.Items {
		items = append(items, &types.VersionListItem{
			ID:        v.ID,
			Version:   v.Version,
			CreatedBy: v.CreatedBy,
			CreatedAt: v.CreatedAt,
		})
	}

	return &types.ListVersionsResponse{
		Versions: items,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}
