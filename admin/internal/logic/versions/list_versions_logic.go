package versions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
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
	tenantID, _ := contextx.GetTenantID(l.ctx)

	versions, err := l.svcCtx.VersionService.ListVersions(l.ctx, shared.TenantID(tenantID), req.PageID, 20)
	if err != nil {
		return nil, err
	}

	items := make([]*types.VersionListItem, 0, len(versions))
	for _, v := range versions {
		items = append(items, &types.VersionListItem{
			ID:        v.ID,
			Version:   v.Version,
			CreatedBy: v.CreatedBy,
			CreatedAt: v.CreatedAt,
		})
	}

	return &types.ListVersionsResponse{
		Versions: items,
	}, nil
}