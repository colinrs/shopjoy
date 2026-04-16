package promotions

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPromotionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPromotionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListPromotionsLogic {
	return ListPromotionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPromotionsLogic) ListPromotions(req *types.ListPromotionsReq) (resp *types.ListPromotionsResp, err error) {
	// Get tenantID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	queryReq := apppromotion.QueryPromotionRequest{
		Name:     req.Name,
		Type:     mapPromotionType(req.Type),
		Status:   mapPromotionStatusToInt(req.Status),
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	listResp, err := l.svcCtx.PromotionApp.ListPromotions(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return nil, err
	}

	list := make([]*types.PromotionDetailResp, len(listResp.List))
	for i, p := range listResp.List {
		list[i] = convertPromotionToDetailResp(p)
	}

	return &types.ListPromotionsResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
