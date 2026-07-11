package user_coupons

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	issueReq := apppromotion.IssueCouponToUserRequest{
		CouponID: req.CouponID,
		UserID:   req.UserID,
	}

	issueResp, err := l.svcCtx.CouponApp.IssueCouponToUser(l.ctx, shared.TenantID(tenantID), issueReq)
	if err != nil {
		return nil, err
	}

	return &types.IssueUserCouponResp{
		ID: issueResp.UserCouponID,
	}, nil
}
