package promotions

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePromotionLogic {
	return UpdatePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdatePromotion mirrors CreatePromotion: parse times, build a
// PromotionRule slice from the wire Rules, then call the unified
// PromotionApp.Update with Rules wrapped in a pointer (so the app
// layer knows to replace — not preserve — the rules).
func (l *UpdatePromotionLogic) UpdatePromotion(req *types.UpdatePromotionReq) (resp *types.PromotionDetailResp, err error) {
	startAt, err := parseTime(req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := parseTime(req.EndTime)
	if err != nil {
		return nil, err
	}

	scope := buildPromotionScopeFromWire(req.ScopeType, req.ScopeIDs, req.ExcludeIDs)
	rules := convertRuleReqsToDomainPtr(req.Rules)

	updateReq := &apppromotion.UpdatePromotionRequest{
		ID:           req.ID,
		Name:         req.Name,
		Description:  req.Description,
		Code:         optionalString(req.Code),
		Type:         parsePromotionType(req.Type),
		MarketID:     optionalInt64(req.MarketID),
		Currency:     currencyWithDefault(req.Currency),
		TotalCount:   optionalInt(req.TotalCount),
		UsageLimit:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		Tags:         req.Tags,
		Scope:        scope,
		StartAt:      startAt,
		EndAt:        endAt,
		Rules:        &rules,
		ActorID:      actorID(l.ctx),
	}

	promotionResp, err := l.svcCtx.PromotionApp.Update(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return convertPromotionToDetailResp(promotionResp), nil
}

// actorID extracts the audit user ID from context with a 0 fallback.
func actorID(ctx context.Context) int64 {
	if v := ctx.Value("user_id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}
