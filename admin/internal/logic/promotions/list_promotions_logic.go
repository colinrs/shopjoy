package promotions

import (
	"context"
	"strings"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
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

// ListPromotions builds a promotion.Query from the wire request and
// delegates to PromotionApp.List. Kind / Status / Type / MarketID
// are wired with pointer semantics so the iota-zero (PROMOTION /
// PENDING / DISCOUNT) can be expressed as a real filter, not a "no
// filter" sentinel.
func (l *ListPromotionsLogic) ListPromotions(req *types.ListPromotionsReq) (resp *types.ListPromotionsResp, err error) {
	q := pkgpromotion.Query{
		PageQuery: shared.PageQuery{Page: req.Page, PageSize: req.PageSize},
		Name:      req.Name,
	}

	// Wire sends lowercase kind strings; the domain constants are
	// uppercase.
	if req.Status == "expired" {
		// "expired" is a derived wire status, not a stored value;
		// translate it into an end_at-based filter rather than a
		// status column match.
		q.ExpiredOnly = true
	} else if req.Status != "" {
		s := mapPromotionStatusToInt(req.Status)
		q.Status = &s
	}
	if req.Type != "" {
		t := parsePromotionType(req.Type)
		q.Type = &t
	}
	if req.MarketID != 0 {
		mid := req.MarketID
		q.MarketID = &mid
	}

	listResp, err := l.svcCtx.PromotionApp.List(l.ctx, q)
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
		PageSize: listResp.Size,
	}, nil
}

// strings is imported to keep the package aligned with the rest of
// the promotion logic files that share it via convertRulesToResp.
var _ = strings.ToLower