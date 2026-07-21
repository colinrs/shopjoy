package regions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/region"
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

// ListRegions 返回区域列表。
//
// 支持的过滤维度：
//   - country_code：ISO 3166-1 alpha-2（例如 CN、US、JP）
//   - parent_code：父级区域 code（空表示返回顶级 = 国家）
//   - level：1=国家 2=省/州 3=市 4=区（仅在 parent_code 非空时生效）
func (l *ListRegionsLogic) ListRegions(req *types.ListRegionsReq) (*types.ListRegionsResp, error) {
	regions, err := l.svcCtx.RegionRepo.FindTree(l.ctx, l.svcCtx.DB, region.RegionQuery{
		CountryCode: req.CountryCode,
		ParentCode:  req.ParentCode,
		Level:       req.Level,
	})
	if err != nil {
		return nil, err
	}

	items := make([]*types.RegionItem, 0, len(regions))
	for _, r := range regions {
		if !r.IsActive {
			continue
		}
		items = append(items, toRegionItem(r))
	}
	return &types.ListRegionsResp{List: items}, nil
}

// toRegionItem 将领域实体转换为 API 响应对象
func toRegionItem(r *region.Region) *types.RegionItem {
	if r == nil {
		return nil
	}
	return &types.RegionItem{
		Code:          r.Code,
		Name:          r.Name,
		Level:         r.Level,
		ParentCode:    r.ParentCode,
		CountryCode:   r.CountryCode,
		PostalPattern: r.PostalPattern,
	}
}