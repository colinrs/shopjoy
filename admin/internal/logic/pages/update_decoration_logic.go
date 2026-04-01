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

type UpdateDecorationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDecorationLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateDecorationLogic {
	return UpdateDecorationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDecorationLogic) UpdateDecoration(req *types.UpdateDecorationRequest) error {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	var blockConfig map[string]any
	if err := json.Unmarshal([]byte(req.BlockConfig), &blockConfig); err != nil {
		blockConfig = make(map[string]any)
	}

	return l.svcCtx.DecorationService.UpdateDecoration(
		l.ctx,
		shared.TenantID(tenantID),
		req.ID,
		blockConfig,
	)
}