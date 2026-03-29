package promotions

import (
	"context"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreatePromotionLogic {
	return CreatePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePromotionLogic) CreatePromotion(req *types.CreatePromotionReq) (resp *types.CreatePromotionResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Parse time
	startAt, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	// Map type string to domain type
	promotionType := mapPromotionType(req.Type)

	createReq := apppromotion.CreatePromotionRequest{
		Name:        req.Name,
		Description: req.Description,
		Type:        promotionType,
		StartAt:     startAt,
		EndAt:       endAt,
		Rules:       make([]apppromotion.CreatePromotionRuleRequest, 0),
	}

	// Build rules from request if applicable
	if req.DiscountType != "" && req.DiscountValue != "" {
		rule := apppromotion.CreatePromotionRuleRequest{
			ConditionType: pkgpromotion.ConditionMinAmount,
			ActionType:    mapDiscountActionType(req.DiscountType),
		}
		if req.MinOrderAmount != "" {
			rule.ConditionValue = parseMoneyToDecimal(req.MinOrderAmount)
		}
		if req.DiscountValue != "" {
			rule.ActionValue = parseMoneyToDecimal(req.DiscountValue)
		}
		if req.MaxDiscount != "" {
			rule.MaxDiscount = parseMoneyToDecimal(req.MaxDiscount)
		}
		createReq.Rules = append(createReq.Rules, rule)
	}

	promotionResp, err := l.svcCtx.PromotionApp.CreatePromotion(l.ctx, shared.TenantID(tenantID), createReq)
	if err != nil {
		return nil, err
	}

	return &types.CreatePromotionResp{
		ID: promotionResp.ID,
	}, nil
}