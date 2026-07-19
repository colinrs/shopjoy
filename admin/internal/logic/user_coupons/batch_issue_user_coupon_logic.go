// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user_coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchIssueUserCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量发放优惠券
func NewBatchIssueUserCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchIssueUserCouponLogic {
	return &BatchIssueUserCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchIssueUserCouponLogic) BatchIssueUserCoupon(req *types.BatchIssueUserCouponReq) (resp *types.BatchIssueUserCouponResp, err error) {
	issued, ids, err := l.svcCtx.PromotionApp.BatchIssue(l.ctx, req.CouponID, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &types.BatchIssueUserCouponResp{
		Issued:        issued,
		UserCouponIDs: ids,
	}, nil
}