package promotions

import (
	"context"
	"strconv"
	"strings"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"

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

	// Parse time
	startAt, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	// Map type + scope. Type comes straight from the wire; scope
	// needs synthesizing because the form posts ID arrays but not
	// scope_type (see wire UpdatePromotionReq). If the form supplies
	// scope_type we trust it, otherwise we derive from whichever ID
	// array is non-empty (PRODUCTS wins when both are populated).
	scope := buildPromotionScope(req.ScopeType, req.ProductIDs, req.CategoryIDs)

	updateReq := apppromotion.UpdatePromotionRequest{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Type:        mapPromotionType(req.Type),
		Scope:       scope,
		UsageLimit:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		Tags:         req.Tags,
		StartAt:     startAt,
		EndAt:       endAt,
	}

	promotionResp, err := l.svcCtx.PromotionApp.UpdatePromotion(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return convertPromotionToDetailResp(promotionResp), nil
}

// buildPromotionScope maps the wire-level (scope_type, product_ids,
// category_ids) triple onto the domain's PromotionScope. Empty IDs
// for the chosen scope are normalized to nil so the storage layer
// doesn't carry stale empty arrays.
//
// Wire values are lowercase ("products", "categories", …) but the
// domain constants are uppercase ("PRODUCTS", …). We uppercase
// before comparing so the form's emitted value matches the enum.
func buildPromotionScope(scopeType string, productIDs, categoryIDs []string) pkgpromotion.PromotionScope {
	normalize := func(s string) pkgpromotion.ScopeType {
		switch strings.ToUpper(s) {
		case string(pkgpromotion.ScopeTypeProducts):
			return pkgpromotion.ScopeTypeProducts
		case string(pkgpromotion.ScopeTypeCategories):
			return pkgpromotion.ScopeTypeCategories
		case string(pkgpromotion.ScopeTypeBrands):
			return pkgpromotion.ScopeTypeBrands
		case string(pkgpromotion.ScopeTypeStorewide):
			return pkgpromotion.ScopeTypeStorewide
		}
		return ""
	}

	st := normalize(scopeType)
	if st == "" {
		switch {
		case len(productIDs) > 0:
			st = pkgpromotion.ScopeTypeProducts
		case len(categoryIDs) > 0:
			st = pkgpromotion.ScopeTypeCategories
		default:
			st = pkgpromotion.ScopeTypeStorewide
		}
	}

	var ids []int64
	switch st {
	case pkgpromotion.ScopeTypeProducts:
		ids = parseInt64Slice(productIDs)
	case pkgpromotion.ScopeTypeCategories:
		ids = parseInt64Slice(categoryIDs)
	default:
		ids = nil
	}
	return pkgpromotion.PromotionScope{Type: st, IDs: ids}
}

// parseInt64Slice converts []string of decimal numeric IDs into
// []int64, ignoring entries that fail to parse. The wire model uses
// string IDs to match the rest of the API; the domain uses int64.
func parseInt64Slice(ss []string) []int64 {
	if len(ss) == 0 {
		return nil
	}
	out := make([]int64, 0, len(ss))
	for _, s := range ss {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			out = append(out, n)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
