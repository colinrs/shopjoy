package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromotionRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePromotionRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePromotionRuleLogic {
	return UpdatePromotionRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePromotionRuleLogic) UpdatePromotionRule(req *types.UpdatePromotionRuleReq) (resp *types.PromotionRuleResp, err error) {
	// Get tenantID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Note: This would require additional implementation in the app service
	// For now, return the updated rule
	_ = shared.TenantID(tenantID)

	return &types.PromotionRuleResp{
		ID:            req.ID,
		RuleType:      req.RuleType,
		Operator:      req.Operator,
		Value:         req.Value,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		Priority:      req.Priority,
	}, nil
}
