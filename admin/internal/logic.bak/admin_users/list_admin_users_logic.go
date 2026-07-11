package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAdminUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAdminUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListAdminUsersLogic {
	return ListAdminUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAdminUsersLogic) ListAdminUsers(req *types.ListAdminUsersRequest) (resp *types.ListAdminUsersResponse, err error) {
	operatorID := contextx.GetCurrentUserID(l.ctx)

	listReq := adminuser.ListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Type:     req.Type,
		Status:   req.Status,
		TenantID: req.TenantID,
	}

	listResp, err := l.svcCtx.AdminUserService.List(l.ctx, operatorID, listReq)
	if err != nil {
		return nil, err
	}

	list := make([]*types.AdminUserInfo, len(listResp.List))
	for i, u := range listResp.List {
		list[i] = &types.AdminUserInfo{
			ID:        u.ID,
			TenantID:  u.TenantID,
			Username:  u.Username,
			Email:     u.Email,
			Mobile:    u.Mobile,
			RealName:  u.RealName,
			Avatar:    u.Avatar,
			Type:      u.Type,
			TypeText:  u.TypeText,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
		}
	}

	return &types.ListAdminUsersResponse{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
