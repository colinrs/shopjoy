package pages

import (
	"context"
	"encoding/json"

	appStorefront "github.com/colinrs/shopjoy/admin/internal/application/storefront"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SaveDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) SaveDraftLogic {
	return SaveDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveDraftLogic) SaveDraft(req *types.SaveDraftRequest) error {


	userID, _ := contextx.GetUserID(l.ctx)
	blocks := make([]*appStorefront.DecorationDTO, 0, len(req.Blocks))
	for _, b := range req.Blocks {
		var blockConfig map[string]any
		if err := json.Unmarshal([]byte(b.BlockConfig), &blockConfig); err != nil {
			return code.ErrInvalidBlockConfig
		}
		blocks = append(blocks, &appStorefront.DecorationDTO{
			BlockType:   b.BlockType,
			BlockConfig: blockConfig,
			SortOrder:   b.SortOrder,
		})
	}

	return l.svcCtx.PageService.SaveDraft(l.ctx, req.ID, blocks, userID)
}
