// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user_coupons

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
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
	issueReq := apppromotion.BatchIssueCouponToUserRequest{
		CouponID: req.CouponID,
		UserIDs:  req.UserIDs,
	}

	issueResp, err := l.svcCtx.CouponApp.BatchIssueCouponToUser(l.ctx, issueReq)
	if err != nil {
		return nil, err
	}

	return &types.BatchIssueUserCouponResp{
		Issued:        int64(len(issueResp.UserCouponIDs)),
		UserCouponIDs: issueResp.UserCouponIDs,
	}, nil
}
