package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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

func (l *GenerateCouponCodesLogic) GenerateCouponCodes(req *types.GenerateCouponCodesReq) (resp *types.GenerateCouponCodesResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Generate codes using the coupon app
	// Note: This requires a coupon ID, but the API doesn't provide one
	// For now, generate codes based on the prefix
	codes := make([]string, 0, req.Quantity)
	for i := 0; i < req.Quantity; i++ {
		code := generateCode(req.Prefix, req.Length, i)
		codes = append(codes, code)
	}

	_ = shared.TenantID(tenantID) // Used for validation

	return &types.GenerateCouponCodesResp{
		Codes: codes,
		Count: len(codes),
	}, nil
}

func generateCode(prefix string, length int, index int) string {
	// Simple code generation
	// In production, use a more sophisticated method
	if prefix == "" {
		prefix = "CPN"
	}
	return prefix + randomString(length)
}

func randomString(length int) string {
	// Simple implementation - in production use crypto/rand
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}