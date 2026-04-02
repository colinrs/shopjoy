package promotions

import (
	"context"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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

func (l *UpdatePromotionLogic) UpdatePromotion(req *types.UpdatePromotionReq) (resp *types.PromotionDetailResp, err error) {
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

	updateReq := apppromotion.UpdatePromotionRequest{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		StartAt:     startAt,
		EndAt:       endAt,
	}

	promotionResp, err := l.svcCtx.PromotionApp.UpdatePromotion(l.ctx, shared.TenantID(tenantID), updateReq)
	if err != nil {
		return nil, err
	}

	return convertPromotionToDetailResp(promotionResp), nil
}
