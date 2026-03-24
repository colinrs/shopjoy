package users

import (
	"context"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	domain "github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersEnhancedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUsersEnhancedLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListUsersEnhancedLogic {
	return ListUsersEnhancedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersEnhancedLogic) ListUsersEnhanced(req *types.ListUsersEnhancedRequest) (resp *types.ListUsersEnhancedResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		tenantID = shared.TenantID(1) // 默认租户
	}

	queryReq := appUser.EnhancedQueryRequest{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Keyword:       req.Keyword,
		Status:        domain.Status(req.Status),
		RegisterStart: req.RegisterStart,
		RegisterEnd:   req.RegisterEnd,
	}

	listResp, err := l.svcCtx.UserService.ExtendedList(l.ctx, tenantID, queryReq)
	if err != nil {
		return nil, err
	}

	users := make([]*types.ExtendedUserResponse, 0, len(listResp.List))
	for _, u := range listResp.List {
		users = append(users, &types.ExtendedUserResponse{
			ID:            u.ID,
			TenantID:      u.TenantID,
			Email:         u.Email,
			Phone:         u.Phone,
			Name:          u.Name,
			Avatar:        u.Avatar,
			Status:        u.Status,
			StatusText:    u.StatusText,
			PointsBalance: u.PointsBalance,
			OrderCount:    u.OrderCount,
			TotalSpent:    u.TotalSpent,
			LastLogin:     u.LastLogin,
			CreatedAt:     u.CreatedAt,
		})
	}

	return &types.ListUsersEnhancedResponse{
		List:     users,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
