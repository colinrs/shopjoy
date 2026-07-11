package users

import (
	"context"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	domainUser "github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListUsersLogic {
	return ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.ListUsersRequest) (resp *types.ListUsersResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		return nil, code.ErrTenantInvalidID
	}

	queryReq := appUser.QueryRequest{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name:   req.Name,
		Email:  req.Email,
		Status: domainUser.Status(req.Status),
	}

	listResp, err := l.svcCtx.UserService.List(l.ctx, tenantID, queryReq)
	if err != nil {
		return nil, err
	}

	list := make([]*types.GetUserResponse, len(listResp.List))
	for i, u := range listResp.List {
		list[i] = &types.GetUserResponse{
			ID:        u.ID,
			Email:     u.Email,
			Phone:     u.Phone,
			Name:      u.Name,
			Avatar:    u.Avatar,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
			LastLogin: u.LastLogin,
		}
	}

	return &types.ListUsersResponse{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
