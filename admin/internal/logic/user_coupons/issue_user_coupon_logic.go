package user_coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IssueUserCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIssueUserCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) IssueUserCouponLogic {
	return IssueUserCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IssueUserCouponLogic) IssueUserCoupon(req *types.IssueUserCouponReq) (resp *types.IssueUserCouponResp, err error) {
	uc, err := l.svcCtx.PromotionApp.IssueToUser(l.ctx, req.CouponID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &types.IssueUserCouponResp{
		ID: uc.ID,
	}, nil
}
