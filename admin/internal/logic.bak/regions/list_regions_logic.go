package regions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListRegionsLogic {
	return ListRegionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRegionsLogic) ListRegions(req *types.ListRegionsReq) (resp *types.ListRegionsResp, err error) {
	// Get regions based on parent_code filter
	if req.ParentCode == "" {
		// Return provinces (level 1)
		return &types.ListRegionsResp{
			List: chinaProvinces,
		}, nil
	}

	// Find children by parent code
	children := findChildren(req.ParentCode)
	return &types.ListRegionsResp{
		List: children,
	}, nil
}

// findChildren finds regions by parent code
func findChildren(parentCode string) []*types.RegionItem {
	// Check cities first
	if cities, ok := chinaCities[parentCode]; ok {
		return cities
	}
	// Check districts
	if districts, ok := chinaDistricts[parentCode]; ok {
		return districts
	}
	return []*types.RegionItem{}
}
