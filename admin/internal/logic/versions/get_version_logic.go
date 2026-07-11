package versions

import (
	"context"
	"encoding/json"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetVersionLogic {
	return GetVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVersionLogic) GetVersion(req *types.GetVersionRequest) (resp *types.VersionDetailResponse, err error) {

	result, err := l.svcCtx.VersionService.GetVersion(l.ctx, req.PageID, req.Version)
	if err != nil {
		return nil, err
	}

	if result == nil || result.Version == nil {
		return &types.VersionDetailResponse{}, nil
	}

	blocks := make([]*types.DecorationDTO, 0, len(result.Blocks))
	for _, b := range result.Blocks {
		blockConfigJSON, _ := json.Marshal(b.BlockConfig)
		blocks = append(blocks, &types.DecorationDTO{
			BlockType:   b.BlockType,
			BlockConfig: string(blockConfigJSON),
			SortOrder:   b.SortOrder,
		})
	}

	return &types.VersionDetailResponse{
		Version: &types.VersionListItem{
			ID:        result.Version.ID,
			Version:   result.Version.Version,
			CreatedBy: result.Version.CreatedBy,
			CreatedAt: result.Version.CreatedAt,
		},
		Blocks: blocks,
	}, nil
}
