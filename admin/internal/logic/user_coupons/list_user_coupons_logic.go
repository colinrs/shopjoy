package user_coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserCouponsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListUserCouponsLogic {
	return ListUserCouponsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserCouponsLogic) ListUserCoupons(req *types.ListUserCouponsReq) (resp *types.ListUserCouponsResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Map status string to domain type
	status := mapUserCouponStatus(req.Status)

	listResp, err := l.svcCtx.CouponApp.ListUserCoupons(
		l.ctx,
		shared.TenantID(tenantID),
		req.UserID,
		status,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return nil, err
	}

	list := make([]*types.UserCouponDetailResp, len(listResp.List))
	for i, uc := range listResp.List {
		list[i] = convertUserCouponToDetailResp(uc)
	}

	return &types.ListUserCouponsResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
