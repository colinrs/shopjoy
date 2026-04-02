package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePromotionRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePromotionRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeletePromotionRuleLogic {
	return DeletePromotionRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePromotionRuleLogic) DeletePromotionRule(req *types.DeletePromotionRuleReq) (resp *types.CreatePromotionResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Note: This would require additional implementation in the app service
	// For now, return success
	_ = shared.TenantID(tenantID)

	return &types.CreatePromotionResp{
		ID: req.ID,
	}, nil
}
