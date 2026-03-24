package pages

import (
	"context"
	"encoding/json"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddDecorationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddDecorationLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddDecorationLogic {
	return AddDecorationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddDecorationLogic) AddDecoration(req *types.AddDecorationRequest) (resp *types.AddDecorationResponse, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	var blockConfig map[string]any
	if err := json.Unmarshal([]byte(req.BlockConfig), &blockConfig); err != nil {
		blockConfig = make(map[string]any)
	}

	result, err := l.svcCtx.DecorationService.AddDecoration(
		l.ctx,
		shared.TenantID(tenantID),
		req.PageID,
		req.BlockType,
		blockConfig,
		req.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	blockConfigJSON, _ := json.Marshal(result.BlockConfig)

	return &types.AddDecorationResponse{
		ID:          result.ID,
		BlockType:   result.BlockType,
		BlockConfig: string(blockConfigJSON),
		SortOrder:   result.SortOrder,
	}, nil
}