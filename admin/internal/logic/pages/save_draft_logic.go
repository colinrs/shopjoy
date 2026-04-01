package pages

import (
	"context"
	"encoding/json"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	appStorefront "github.com/colinrs/shopjoy/admin/internal/application/storefront"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

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

	return l.svcCtx.PageService.SaveDraft(l.ctx, shared.TenantID(tenantID), req.ID, blocks, userID)
}