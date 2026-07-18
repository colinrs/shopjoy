package promotions

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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

	// Only attach status / type filters when the frontend actually
	// provided a value. Using the mapped int as a sentinel — as the
	// old code did — silently dropped filters whose mapped value
	// equals the iota-zero (StatusPending / TypeDiscount).
	queryReq := apppromotion.QueryPromotionRequest{
		Name:     req.Name,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if req.Status == "expired" {
		// "expired" is a derived wire status, not a stored value;
		// translate it into an end_at-based filter rather than a
		// status column match (StatusExpired doesn't exist in the enum).
		queryReq.ExpiredOnly = true
	} else if req.Status != "" {
		s := mapPromotionStatusToInt(req.Status)
		queryReq.Status = &s
	}
	if req.Type != "" {
		t := mapPromotionType(req.Type)
		queryReq.Type = &t
	}

	listResp, err := l.svcCtx.PromotionApp.ListPromotions(l.ctx, queryReq)
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
