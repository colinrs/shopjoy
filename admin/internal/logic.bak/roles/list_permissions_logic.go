package roles

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListPermissionsLogic {
	return ListPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPermissionsLogic) ListPermissions(req *types.ListPermissionsRequest) (resp *types.ListPermissionsResponse, err error) {
	// Get all permissions
	permissions, err := l.svcCtx.PermissionRepo.FindAll(l.ctx, l.svcCtx.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response
	list := make([]*types.PermissionInfo, 0, len(permissions))
	for _, p := range permissions {
		list = append(list, &types.PermissionInfo{
			ID:       p.ID,
			Name:     p.Name,
			Code:     p.Code,
			Type:     int8(p.Type),                        // #nosec G115 // type values are small (tinyint range)
			TypeText: getPermissionTypeText(int8(p.Type)), // #nosec G115 // type values are small (tinyint range)
			ParentID: p.ParentID,
			Path:     p.Path,
			Icon:     p.Icon,
			Sort:     p.Sort,
		})
	}

	return &types.ListPermissionsResponse{
		List: list,
	}, nil
}
