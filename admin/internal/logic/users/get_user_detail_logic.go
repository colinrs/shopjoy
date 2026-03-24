package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUserDetailLogic {
	return GetUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserDetailLogic) GetUserDetail(req *types.GetUserRequest) (resp *types.UserDetailResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		tenantID = shared.TenantID(1) // 默认租户
	}

	detail, err := l.svcCtx.UserService.GetDetail(l.ctx, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.UserDetailResponse{
		ID:            detail.ID,
		TenantID:      detail.TenantID,
		Email:         detail.Email,
		Phone:         detail.Phone,
		Name:          detail.Name,
		Avatar:        detail.Avatar,
		Gender:        detail.Gender,
		GenderText:    detail.GenderText,
		Birthday:      detail.Birthday,
		Status:        detail.Status,
		StatusText:    detail.StatusText,
		PointsBalance: detail.PointsBalance,
		PointsFrozen:  detail.PointsFrozen,
		TotalEarned:   detail.TotalEarned,
		TotalRedeemed: detail.TotalRedeemed,
		OrderCount:    detail.OrderCount,
		TotalSpent:    detail.TotalSpent,
		ReviewCount:   detail.ReviewCount,
		LastLogin:     detail.LastLogin,
		CreatedAt:     detail.CreatedAt,
		UpdatedAt:     detail.UpdatedAt,
	}, nil
}
