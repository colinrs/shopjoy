package coupons

import (
	"context"
	"encoding/json"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCouponCodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCouponCodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GenerateCouponCodesLogic {
	return GenerateCouponCodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GenerateCouponCodes parses the wire CouponConfig JSON, hands the
// resulting map to PromotionApp.GenerateCodes, and reports the codes
// it produced. GenerateCodes internally builds a COUPON-kind
// Promotion per code via couponFromConfig (see promotion_app.go).
//
// Per handoff I2, COUPON rows may carry usage_limit=0 ("unlimited")
// after the data migration; the JSON parsing below preserves that —
// a zero usage_limit leaves TotalCount unset (nil) on the resulting
// Promotion so the consume-inventory SQL guard short-circuits.
func (l *GenerateCouponCodesLogic) GenerateCouponCodes(req *types.GenerateCouponCodesReq) (resp *types.GenerateCouponCodesResp, err error) {
	cfg := map[string]any{}
	if req.CouponConfig != "" {
		if err := json.Unmarshal([]byte(req.CouponConfig), &cfg); err != nil {
			return nil, err
		}
	}
	// Inject context-derived defaults so couponFromConfig in the
	// app layer can attach them to the COUPON-kind Promotion it
	// builds. The actor id is needed for audit, the tenant for
	// scoping.
	cfg["_actor_id"] = actorID(l.ctx)
	cfg["_tenant_id"] = tenantID(l.ctx)

	codes, err := l.svcCtx.PromotionApp.GenerateCodes(l.ctx, req.Prefix, req.Quantity, cfg)
	if err != nil {
		return nil, err
	}
	return &types.GenerateCouponCodesResp{
		Codes: codes,
		Count: len(codes),
	}, nil
}